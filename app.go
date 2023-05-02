package main

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"strings"
	"time"
)

type App struct {
	ctx         context.Context
	relayPool   *RelayPool
	lastRefresh *nostr.Timestamp
	config      *Config
}

var (
	appName = "Greet"

	followedPks []string
	db          DB
)

const (
	QUERY_SIZE   = 25
	POLL_SECONDS = 60
)

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {

	log.Info().Msg("Starting")
	a.ctx = ctx

	t := nostr.Now() - 43200 // 12h
	a.lastRefresh = &t
	db = NewDB()

	a.config = NewConfig()
	err := a.config.Load()
	if err != nil {
		log.Error().Msg("Error: Could not configuration file: " + err.Error())
	}

	setupLogging()

	a.relayPool = NewRelayPool()
	for _, r := range a.config.Relays {
		if r.Enabled {
			err := a.relayPool.Add(r)
			if err != nil {
				log.Err(err)
			}
		}
	}

}

func setupLogging() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

func (a *App) OnBeforeClose(ctx context.Context) bool {
	return false
}

func (a *App) OnShutdown(ctx context.Context) {
	log.Info().Msg("Shutting down")
	a.relayPool.DisconnectAll()
	a.relayPool.RemoveAll()
}

func (a *App) Quit() {
	log.Info().Msg("QUIT called")
	runtime.Quit(a.ctx)
}

func (a *App) OnDomReady(ctx context.Context) {
	var err error

	log.Debug().Msg("Checking private key...")
	key := string(a.config.Privkey)
	if key == "" {
		log.Debug().Msg("...blank. Launch login")
		go runtime.EventsEmit(a.ctx, "evLoginDialog")
	} else {
		if strings.HasPrefix(key, "ENC:") {
			log.Debug().Msg("...ENC:encrypted. Launch PIN dialog")
			go runtime.EventsEmit(a.ctx, "evPinDialog")
		} else {
			log.Debug().Msg("...use configured")
			a.config.privKeyHex = key
			a.config.pubkey, err = nostr.GetPublicKey(key)
			if err != nil {
				log.Err(err)
			}
			a.RefreshContactProfiles()
			a.RefreshFeed("evFollowEventNote", true)
		}
	}

	a.StartPolling()
}

func (A *App) DumpEvents() {
	db.DumpEvents()
}

func (a *App) StartPolling() {
	log.Debug().Msgf("Starting polling every %d seconds", POLL_SECONDS)

	// Check for new notes
	go func() {
		loops := 0
		for true {
			time.Sleep(time.Second * POLL_SECONDS)
			a.RefreshFeed("evRefreshNote", false)
			if loops >= 15 {
				a.relayPool.ReconnectAll()
				loops = 0
			}

			loops++
		}
	}()
}

func (a *App) RefreshContactProfiles() {
	log.Debug().Msg("Refreshing Contact Profiles")
	followedPks = a.GetContactList(a.config.pubkey)
	a.GetMetadataEvents(followedPks)
}

func (a *App) RefreshFeed(postEvent string, repost bool) {
	if len(followedPks) == 0 {
		return
	}
	chks := chunkSlice(followedPks, QUERY_SIZE)
	for _, chk := range chks {
		a.GetTextNotesByPubkeys(chk, postEvent, repost)
	}
}

func (a *App) RefreshFeedReset(postEvent string) {
	t := nostr.Now() - 21600 // 6h
	a.lastRefresh = &t
	a.RefreshFeed(postEvent, true)
}

func (a *App) PkToNpub(pk string) (string, error) {
	npub, err := nip19.EncodePublicKey(pk)
	return npub, err
}

func (a *App) Nip19Decode(uri string) (string, error) {
	_, val, err := nip19.Decode(uri)
	if err != nil {
		log.Err(err)
		return "", err
	}
	log.Trace().Msgf("Nip19Decode: %s -> %s", uri, val)
	return fmt.Sprintf("%v", val), nil
}

func (a *App) GetContactList(pk string) []string {
	log.Debug().Msg("Getting contact list")
	pks := []string{}
	if pk == "" {
		return pks
	}

	evs := make(chan *nostr.Event)
	defer close(evs)
	pending := a.relayPool.Query(&nostr.Filter{
		Authors: []string{pk},
		Kinds: []int{
			nostr.KindContactList,
		},
	}, evs)

	for ev := range evs {
		if ev == nil {
			pending--
			if pending == 0 {
				break
			}
			continue
		}

		fmt.Println("GetContactList: looking at event", ev.ID)
		tags := ev.Tags.GetAll([]string{"p"})

		fmt.Println("GetContactList: pTag count", len(tags))
		// TODO: Change this to track the relays the contact is known at
		for a := 0; a < len(tags); a++ {
			if !contains(pks, tags[a].Value()) {
				pks = append(pks, tags[a].Value())
			}
		}
	}

	return pks
}

func (a *App) GetMetadataEvents(pks []string) {
	if len(pks) == 0 {
		log.Warn().Msg("Getting metadata events called with no contacts!")
		return
	}
	log.Debug().Msgf("Getting metadata events for %d keys: %s", len(pks), pks)

	for i := range pks {
		if pks[i] == "" {
			log.Error().Msg("Caught empty PK") // FIXME: why does this happen?
			return
		}
	}

	evs := make(chan *nostr.Event)
	pending := a.relayPool.Query(&nostr.Filter{
		Authors: pks,
		Kinds: []int{
			nostr.KindSetMetadata,
		},
		Limit: len(pks),
	}, evs)

	for ev := range evs {
		if ev == nil {
			pending--
			if pending == 0 {
				break
			}
			continue
		}
		db.AddEvent(ev.ID, ev)
		cm, err := getContentMeta(ev)
		if err != nil {
			fmt.Errorf("Error parsing metadata for event", ev.ID, ":", err.Error())
			continue
		}
		npub, err := a.PkToNpub(ev.PubKey)
		if err != nil {
			fmt.Errorf("Error converting PK to NPUB for event", ev.ID, ":", err.Error())
			continue
		}

		profile := Profile{
			Pk:        ev.PubKey,
			Following: contains(followedPks, ev.PubKey),
			Meta:      *cm,
			Npub:      npub,
		}

		db.AddProfile(ev.PubKey, &profile) // Overwrite if existing

		if profile.Following {
			go runtime.EventsEmit(app.ctx, "evMetadata", profile)
		}
	}
}

func (a *App) GetTaggedProfiles(parentEvent string) []*Profile {
	log.Debug().Msgf("GetTaggedProfiles called for parent event %s", parentEvent)
	cachedProfiles := []*Profile{}
	missingProfiles := []string{}
	dups := []string{}

	ev := db.GetEvent(parentEvent)
	if ev != nil {
		var pTags []nostr.Tag = nostr.Tags.GetAll(ev.Tags, []string{"p"})
		for a := 0; a < len(pTags); a++ {
			pk := pTags[a].Value()
			if !contains(dups, pk) {
				dups = append(dups, pk)
				profile := db.GetProfile(pk)
				if profile != nil {
					cachedProfiles = append(cachedProfiles, profile)
				} else {
					missingProfiles = append(missingProfiles, pk)
				}
			}
		}

		if len(missingProfiles) > 0 {
			a.GetMetadataEvents(missingProfiles)
			for a := 0; a < len(missingProfiles); a++ {
				pk := missingProfiles[a]
				profile := db.GetProfile(pk)
				if profile != nil {
					cachedProfiles = append(cachedProfiles, profile)
				}
			}
		}
	}

	return cachedProfiles
}

func (a *App) GetTaggedEvents(parentEvent string) []*nostr.Event {
	cachedEvents := []*nostr.Event{}
	missingEvents := []string{}

	ev := db.GetEvent(parentEvent)
	if ev != nil {
		var eTags []nostr.Tag = nostr.Tags.GetAll(ev.Tags, []string{"e"})
		for a := 0; a < len(eTags); a++ {
			evId := eTags[a].Value()
			event := db.GetEvent(evId)
			if event != nil {
				cachedEvents = append(cachedEvents, event)
			} else {
				missingEvents = append(missingEvents, evId)
			}
		}

		if len(missingEvents) > 0 {
			a.GetTextNotesByEventIds(missingEvents)
			for a := 0; a < len(missingEvents); a++ {
				evId := missingEvents[a]
				event := db.GetEvent(evId)
				if event != nil {
					cachedEvents = append(cachedEvents, event)
				}
			}
		}
	}

	return cachedEvents
}

func (a *App) GetContactProfile(pk string) *Profile {
	if strings.HasPrefix(pk, "npub") {
		pk, _ = a.Nip19Decode(pk)
	}
	if db.HasProfile(pk) {
		log.Debug().Msgf("GetContactProfile for PK %s (cache)", pk)
		return db.GetProfile(pk)
	}
	log.Debug().Msgf("GetContactProfile for PK %s (query)", pk)
	a.GetMetadataEvents([]string{pk})
	if db.HasProfile(pk) {
		return db.GetProfile(pk)
	}
	npub, _ := a.PkToNpub(pk)
	return &Profile{
		Pk:        pk,
		Following: false,
		Meta:      ProfileMetadata{},
		Npub:      npub,
		Relays:    nil,
	}
}

func (a *App) GetWritableRelays() []*string {
	rs := []*string{}
	for _, r := range a.relayPool.pool {
		if r.Enabled && r.Write {
			rs = append(rs, &r.Url)
		}
	}
	return rs
}

func (a *App) GetRelays() []*Relay {
	return a.config.Relays
}

func (a *App) SetRelays(r []*Relay) {
	a.relayPool.DisconnectAll()
	a.relayPool.RemoveAll()

	a.relayPool.AddAll(r)
	a.config.Relays = r

	err := a.config.Save()
	if err != nil {
		log.Error().Msgf("Error saving config file: %s", err.Error())
		return
	}

	a.RefreshContactProfiles()
	go a.RefreshFeed("evFollowEventNote", true)
}

func (a *App) GetTextNotesByPubkeys(pks []string, postEvent string, repost bool) {
	log.Debug().Msgf("Getting text events from contact pks...%d", len(pks))
	a.GetTextNotesByPubkeysOptions(pks, []int{nostr.KindTextNote, nostr.KindBoost}, a.lastRefresh, 0, postEvent, repost)
	t := nostr.Now()
	a.lastRefresh = &t
}

func (a *App) GetTextNotesByPubkeysOptions(pks []string, kinds []int, since *nostr.Timestamp, limit int, postEvent string, repost bool) {
	if len(pks) == 0 {
		return
	}

	evs := make(chan *nostr.Event)
	filter := nostr.Filter{
		Authors: pks,
		Kinds:   kinds,
		Since:   since,
	}
	if limit > 0 {
		filter.Limit = limit
	}

	pending := a.relayPool.Query(&filter, evs)

	for ev := range evs {
		if ev == nil {
			pending--
			if pending == 0 {
				break
			}
			continue
		}
		ev.SetExtra("when", humanize.Time(ev.CreatedAt.Time()))
		existingEvent := db.GetEvent(ev.ID)
		db.AddEvent(ev.ID, ev)
		if existingEvent == nil || repost {
			if postEvent != "" {
				runtime.EventsEmit(app.ctx, postEvent, ev)
			}
		}
	}
}

func (a *App) GetTextNotesByEventIds(ids []string) []*nostr.Event {
	log.Debug().Msgf("GetTextNotesByEventIds: %s", ids)
	events := []*nostr.Event{}
	if len(ids) == 0 {
		return events
	}

	evs := make(chan *nostr.Event)
	pending := a.relayPool.Query(&nostr.Filter{
		IDs: ids,
		Kinds: []int{
			nostr.KindTextNote,
		},
	}, evs)

	for ev := range evs {
		if ev == nil {
			pending--
			if pending == 0 {
				break
			}
			continue
		}
		ev.SetExtra("when", humanize.Time(ev.CreatedAt.Time()))
		db.AddEvent(ev.ID, ev)
		events = append(events, ev)
	}

	log.Debug().Msgf("GetTextNotesByEventIds returning %d events", len(events))
	return events
}

func (a *App) PostEvent(kind int, tags nostr.Tags, content string) {
	ev := nostr.Event{
		PubKey:    a.config.pubkey,
		CreatedAt: nostr.Now(),
		Kind:      kind,
		Tags:      tags,
		Content:   content,
	}
	ev.Sign(a.config.privKeyHex)

	for _, r := range a.relayPool.pool {
		if r.Enabled && r.Write {
			r.conn.Publish(r.rootCtx, ev)
			log.Info().Msgf("Published %s to relay %s", ev.ID, r.Url)
		}
	}
	runtime.EventsEmit(app.ctx, "evRefreshNote", ev)
}

func (a *App) PublishContentToSelectedRelays(kind int, content string, ts [][]string, relays []string) {
	tags := nostr.Tags{}

	for _, tag := range ts {
		tags = append(tags, nostr.Tag{
			tag[0],
			tag[1],
		})
	}

	ev := nostr.Event{
		PubKey:    a.config.pubkey,
		CreatedAt: nostr.Now(),
		Kind:      kind,
		Tags:      tags,
		Content:   content,
	}
	ev.Sign(a.config.privKeyHex)

	for _, url := range relays {
		r := a.relayPool.GetRelayByUrl(url)
		if r.Enabled && r.Write {
			r.conn.Publish(r.rootCtx, ev)
			log.Info().Msgf("Published %s to %s", ev.ID, r.Url)
		}
	}
	runtime.EventsEmit(app.ctx, "evRefreshNote", ev)
}

func (a *App) FollowContact(pk []string) error {
	// Append to existing follows
	followedPks = append(followedPks, pk...)

	// Post a new list to all relays
	tags := nostr.Tags{}
	for _, tag := range followedPks {
		tags = append(tags, nostr.Tag{
			"p",
			tag,
		})
	}

	a.PostEvent(3, tags, "")

	runtime.EventsEmit(a.ctx, "evRefreshContacts")

	return nil
}

func (a *App) UnfollowContact(pk string) error {
	// Remove PK from existing follows
	newFollows := []string{}

	for _, follow := range followedPks {
		if pk != follow {
			newFollows = append(newFollows, follow)
		}
	}

	// Post a new list to all relays
	tags := nostr.Tags{}
	for _, tag := range newFollows {
		tags = append(tags, nostr.Tag{
			"p",
			tag,
		})
	}

	followedPks = newFollows
	a.PostEvent(3, tags, "")

	// Tell frontend to refresh list
	runtime.EventsEmit(a.ctx, "evRefreshContacts")

	return nil
}

func (a *App) DeleteEvent(evId string) {
	ev := nostr.Event{
		PubKey: a.config.pubkey,
		Kind:   5,
		Tags: nostr.Tags{
			nostr.Tag{"e", evId},
		},
		Content: "Deletion request",
	}
	ev.Sign(a.config.privKeyHex)

	for _, r := range a.relayPool.pool {
		if r.Enabled && r.Write {
			r.conn.Publish(r.rootCtx, ev)
		}
		log.Info().Msgf("Delete %s requested to %s", ev.ID, r.Url)
	}
}

func (a *App) GetMyPubkey() string {
	return a.config.pubkey
}

func (a *App) SaveConfigDark(dark bool) {
	a.config.Dark = dark
	err := a.config.Save()
	if err != nil {
		log.Err(err)
		return
	}
}

func (a *App) SetLoginWithPrivKey(keypin []string) error {
	var err error
	var cipher []byte

	if len(keypin) != 2 {
		return errors.New("Input error: expected key and PIN")
	}
	key := keypin[0]
	pin := keypin[1]

	if strings.HasPrefix(key, "nsec") {
		key, err = a.Nip19Decode(key)
	}
	if err != nil {
		return err
	}

	if pin == "" {
		a.config.Privkey = key
	} else {
		cipher, err = encrypt([]byte(key), pin)
		if err != nil {
			return err
		}
		a.config.Privkey = "ENC:" + b64.StdEncoding.EncodeToString([]byte(cipher))
	}

	a.config.pubkey, err = nostr.GetPublicKey(key)
	a.config.privKeyHex = key
	if err != nil {
		return err
	}
	a.config.Save()

	if len(a.config.Relays) == 0 {
		// Add some default relays
		relays := []*Relay{}
		addrs := []string{
			"wss://nos.lol",
			"wss://relay.damus.io",
			"wss://relay.snort.social",
			"wss://nostr.mom",
		}

		for _, addr := range addrs {
			relays = append(relays, &Relay{
				Url:     addr,
				Read:    true,
				Write:   true,
				Enabled: true,
			})
		}
		a.SetRelays(relays)
	}

	return nil
}

func (a *App) LoginWithPin(pin string) error {
	log.Debug().Msg("PIN login called")
	parts := strings.SplitAfter(a.config.Privkey, "ENC:")
	if len(parts) != 2 {
		return errors.New("Private key does not appear to be encrypted")
	}
	dec64, err := b64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Err(err)
		return errors.New("Error decoding the private key")
	}

	key, err := decrypt(dec64, pin)
	if err != nil {
		return errors.New("Wrong PIN")
	}
	a.config.privKeyHex = string(key)
	a.config.pubkey, err = nostr.GetPublicKey(a.config.privKeyHex)

	go func() {
		a.RefreshContactProfiles()
		a.RefreshFeed("evFollowEventNote", true)
	}()

	log.Info().Msgf("PIN login success for %s", a.config.pubkey)
	return nil
}

func (a *App) GenerateKeys() (*map[string]string, error) {
	log.Debug().Msg("Generating new keys")
	keyDetail := make(map[string]string)

	key := nostr.GeneratePrivateKey()
	pk, _ := nostr.GetPublicKey(key)

	keyDetail["key"] = key
	keyDetail["pk"] = pk
	return &keyDetail, nil
}

func (a *App) SaveNewKeys(creds map[string]string) error {
	log.Info().Msg("Saving new account details...")
	key := creds["key"]
	name := creds["name"]
	displayName := creds["displayName"]
	pin := creds["pin"]

	a.SetLoginWithPrivKey([]string{key, pin})

	// Set profile with this name/display name
	meta := ProfileMetadata{
		Name:        name,
		About:       "",
		Picture:     "",
		NIP05:       "",
		DisplayName: displayName,
		Lud06:       "",
		Banner:      "",
		Website:     "",
	}

	content, err := json.Marshal(&meta)
	if err != nil {
		log.Err(err)
		return err
	}

	a.PostEvent(nostr.KindSetMetadata, nostr.Tags{}, string(content))

	return nil
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(slice) == 0 {
			break
		}

		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
