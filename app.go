package main

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type App struct {
	ctx       context.Context
	relayPool *RelayPool
	config    *Config
	logging   zerolog.Level
}

var (
	appName = "Greet"

	followedPks []string
	db          DB
)

const (
	QUERY_SIZE   = 25
	POLL_SECONDS = 60
	SECS_6H      = 21600
	SECS_12H     = 43200
	SECS_24H     = 86400
)

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {

	a.ctx = ctx
	setupLogging()

	log.Info().Msg("Starting up...")
	a.config = NewConfig()
	err := a.config.Load()
	if err != nil {
		log.Error().Msg("Error: Could not configuration file: " + err.Error())
	}
	db = NewDB()
	a.relayPool = NewRelayPool()
	for _, r := range a.config.Relays {
		if r.Enabled {
			err = a.relayPool.Add(r)
			if err != nil {
				log.Err(err)
			}
		}
	}

	// Maintenance loop
	go func() {
		for {
			a.CheckRelays()
			a.PingTimer()
			time.Sleep(time.Second * 10)
		}
	}()

	log.Info().Msg("...start up done")
}

func setupLogging() {
	zerolog.SetGlobalLevel(app.logging)
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
		log.Debug().Msg("...key blank. Launch login")
		go func() {
			time.Sleep(time.Second * 2)
			runtime.EventsEmit(a.ctx, "evLoginDialog")
		}()
	} else {
		if strings.HasPrefix(key, "ENC:") {
			log.Debug().Msg("...key ENC:encrypted. Launch PIN dialog")
			go func() {
				time.Sleep(time.Second * 2)
				runtime.EventsEmit(a.ctx, "evPinDialog")
			}()
		} else {
			log.Debug().Msg("...use configured key")
			a.config.privKeyHex = key
			a.config.pubkey, err = nostr.GetPublicKey(key)
			if err != nil {
				log.Panic()
			}
			go func() {
				time.Sleep(time.Second * 2)
				runtime.EventsEmit(a.ctx, "evPkChange", a.config.pubkey)
			}()
		}
	}
	a.CheckRelays()
}

func (a *App) BeginSubscriptions() {
	a.RefreshContactProfiles()
	a.SubscribeToFeedForPubkeys(followedPks, true)
}

func (A *App) DumpEvents() {
	db.DumpEvents()
}

func (a *App) RefreshContactProfiles() {
	log.Debug().Msg("Refreshing Contact Profiles")
	followedPks = a.GetContactList(a.config.pubkey)

	chks := chunkSlice(followedPks, QUERY_SIZE)
	for _, chk := range chks {
		a.GetMetadataEvents(chk)
	}
}

func (a *App) RefreshFeed(repost bool) {
	if len(followedPks) == 0 {
		return
	}

	chks := chunkSlice(followedPks, QUERY_SIZE)
	for _, chk := range chks {
		a.SubscribeToFeedForPubkeys(chk, repost)
	}
}

func (a *App) RefreshFeedReset() {
	log.Debug().Msg("Resetting feed")
	a.relayPool.UnsubscribeAll()
	a.RefreshFeed(true)
}

func (a *App) PkToNpub(pk string) (string, error) {
	npub, err := nip19.EncodePublicKey(pk)
	return npub, err
}

func (a *App) Nip19Decode(uri string) ([]string, error) {
	result := []string{}
	prefix, val, err := nip19.Decode(uri)
	if err != nil {
		log.Err(err)
		return result, err
	}
	log.Debug().Msgf("Nip19Decode: type %s, %s -> %s %v", reflect.TypeOf(val), uri, prefix, val)

	result = append(result, prefix)
	switch val.(type) {
	case string:
		result = append(result, val.(string))
	case nostr.EventPointer:
		ep := val.(nostr.EventPointer)
		result = append(result, ep.ID, ep.Author, string(ep.Kind))
		for _, r := range ep.Relays {
			result = append(result, r)
		}
	default:
		result = append(result, fmt.Sprint(val))
	}

	return result, nil
}

func (a *App) GetContactList(pk string) []string {
	log.Debug().Msgf("Getting contact list for %s", pk)

	pks := []string{}
	if pk == "" {
		return pks
	}

	ch := make(chan *nostr.Event)
	go func() {
		for ev := range ch {
			tags := ev.Tags.GetAll([]string{"p"})
			for a := 0; a < len(tags); a++ {
				if !contains(pks, tags[a].Value()) {
					pks = append(pks, tags[a].Value())
				}
			}
		}
	}()
	a.relayPool.QuerySync(&nostr.Filter{
		Authors: []string{pk},
		Kinds: []int{
			nostr.KindContactList,
		},
	}, ch)

	return pks
}

func (a *App) GetMetadataEvents(pks []string) {
	if len(pks) == 0 {
		log.Warn().Msg("Getting metadata events called with no contacts!")
		return
	}
	log.Debug().Msgf("Getting metadata events for %d keys: %s", len(pks), pks)

	ch := make(chan *nostr.Event)
	go func() {
		for ev := range ch {
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
	}()

	a.relayPool.QuerySync(&nostr.Filter{
		Authors: pks,
		Kinds: []int{
			nostr.KindSetMetadata,
		},
		Limit: len(pks),
	}, ch)
}

func (a *App) GetTaggedProfiles(parentEvent string) []*Profile {
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
	} else {
		log.Debug().Msgf("GetTaggedProfiles called for parent event %s (not cached)", parentEvent)
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

func (a *App) GetContactProfile(pk string) (*Profile, error) {
	if strings.HasPrefix(pk, "npub") {
		val, err := a.Nip19Decode(pk)
		if err != nil {
			return nil, err
		}
		pk = val[1]
	}
	if db.HasProfile(pk) {
		log.Trace().Msgf("GetContactProfile for PK %s (cache)", pk)
		return db.GetProfile(pk), nil
	}
	log.Trace().Msgf("GetContactProfile for PK %s (query)", pk)
	a.GetMetadataEvents([]string{pk})
	if db.HasProfile(pk) {
		return db.GetProfile(pk), nil
	}
	npub, _ := a.PkToNpub(pk)
	return &Profile{
		Pk:        pk,
		Following: false,
		Meta:      ProfileMetadata{},
		Npub:      npub,
		Relays:    nil,
	}, nil
}

func (a *App) GetReadableRelays() []*string {
	rs := []*string{}
	for _, r := range a.relayPool.pool {
		if r.Enabled && r.Read && (r.conn.ConnectionError == nil) {
			rs = append(rs, &r.Url)
		}
	}
	return rs
}

func (a *App) GetWritableRelays() []*string {
	rs := []*string{}
	for _, r := range a.relayPool.pool {
		if r.Enabled && r.Write && (r.conn.ConnectionError == nil) {
			rs = append(rs, &r.Url)
		}
	}
	return rs
}

func (a *App) GetRelays() []*RelayStruct {
	return a.config.Relays
}

func (a *App) SetRelays(r []*RelayStruct) {
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
	go a.RefreshFeed(false)
}

func (a *App) GetTextNotesForPubkeys(pks []string, postEvent string, repost bool) error {
	log.Debug().Msgf("Getting text events for pks...(%d)", len(pks))

	if len(pks) == 0 {
		return nil
	}

	ch := make(chan *nostr.Event)
	go func() {
		for ev := range ch {
			existingEvent := db.GetEvent(ev.ID)
			db.AddEvent(ev.ID, ev)
			if existingEvent == nil || repost {
				runtime.EventsEmit(app.ctx, postEvent, ev)
			}
		}
	}()

	a.relayPool.QuerySync(&nostr.Filter{
		Authors: pks,
		Kinds:   []int{nostr.KindTextNote, nostr.KindBoost},
		Limit:   100,
	}, ch)

	return nil
}

func (a *App) SubscribeToFeedForPubkeys(pks []string, repost bool) {
	if len(pks) == 0 {
		return
	}
	ch := make(chan *nostr.Event)
	ch1 := make(chan *nostr.Event)
	go func() {
		for ev := range ch {
			existingEvent := db.GetEvent(ev.ID)
			db.AddEvent(ev.ID, ev)
			if existingEvent == nil || repost {
				runtime.EventsEmit(app.ctx, "evFollowEventNote", ev)
			}
		}
	}()
	go func() {
		for ev := range ch1 {
			existingEvent := db.GetEvent(ev.ID)
			db.AddEvent(ev.ID, ev)
			if existingEvent == nil || repost {
				runtime.EventsEmit(app.ctx, "evRefreshNote", ev)
			}
		}
	}()
	since := nostr.Now() - SECS_6H
	filter := nostr.Filter{
		Authors: pks,
		Kinds: []int{
			nostr.KindTextNote,
			nostr.KindBoost,
		},
		Since: &since,
	}

	a.relayPool.Subscribe(&filter, ch, ch1)
}

func (a *App) GetTextNotesByEventIds(ids []string) []*nostr.Event {
	log.Debug().Msgf("GetTextNotesByEventIds: %s", ids)
	events := []*nostr.Event{}
	if len(ids) == 0 {
		return events
	}

	ch := make(chan *nostr.Event)
	go func() {
		for ev := range ch {
			db.AddEvent(ev.ID, ev)
			events = append(events, ev)
		}
	}()
	a.relayPool.QuerySync(&nostr.Filter{
		IDs: ids,
		Kinds: []int{
			nostr.KindTextNote,
			nostr.KindBoost,
		},
	}, ch)

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
			r.conn.Publish(context.Background(), ev)
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
			r.conn.Publish(context.Background(), ev)
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
			r.conn.Publish(context.Background(), ev)
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
		val, e := a.Nip19Decode(key)
		if e != nil {
			return err
		}
		key = val[1]
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
	runtime.EventsEmit(a.ctx, "evPkChange", a.config.pubkey)

	if len(a.config.Relays) == 0 {
		// Add some default relays
		relays := []*RelayStruct{}
		addrs := []string{
			"wss://nos.lol",
			"wss://relay.damus.io",
			"wss://relay.snort.social",
			"wss://nostr.mom",
		}

		for _, addr := range addrs {
			relays = append(relays, &RelayStruct{
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
	runtime.EventsEmit(a.ctx, "evPkChange", a.config.pubkey)

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
		Lud16:       "",
		Banner:      "",
		Website:     "",
	}

	return a.SaveProfile(meta)
}

func (a *App) SaveContacts() (*string, error) {
	profile := db.GetProfile(a.config.pubkey)
	filename := fmt.Sprintf("%s-%s.json", profile.Meta.Name, profile.Npub)
	path := filepath.Join(a.config.configDir, filename)

	configOutput, err := PrettyStruct(followedPks)
	if err != nil {
		return nil, err
	}

	f, err := openFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	f.Write([]byte(configOutput + "\n"))

	return &path, nil
}

func (a *App) RestoreContacts() (*string, error) {
	profile := db.GetProfile(a.config.pubkey)
	filename := fmt.Sprintf("%s-%s.json", profile.Meta.Name, profile.Npub)
	path := filepath.Join(a.config.configDir, filename)
	f, err := openFile(path, os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileinfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	_, err = f.Read(buffer)
	if err != nil {
		return nil, err
	}

	contacts := []string{}
	err = json.Unmarshal(buffer, &contacts)
	if err != nil {
		log.Err(err)
	}

	log.Debug().Msgf("Loaded %d contacts from %s", len(contacts), path)

	if len(contacts) > 0 {
		followedPks = []string{}
		a.FollowContact(contacts)
	} else {
		return nil, errors.New("No contacts in file. Changes not published")
	}

	return &path, nil
}

func (a *App) SaveProfile(metadata ProfileMetadata) error {
	log.Debug().Msg("Saving profile")
	content, err := json.Marshal(metadata)
	if err != nil {
		log.Err(err)
		return err
	}
	a.PostEvent(nostr.KindSetMetadata, nostr.Tags{}, string(content))
	return nil
}

func (a *App) CheckRelays() {
	readable := a.GetReadableRelays()
	writable := a.GetWritableRelays()
	numSubs := 0
	for _, url := range readable {
		relay := a.relayPool.GetRelayByUrl(*url)
		numSubs += len(relay.subs)
	}

	opts := make(map[string]int)
	opts["readable"] = len(readable)
	opts["writable"] = len(writable)
	opts["subs"] = numSubs
	runtime.EventsEmit(a.ctx, "evRelayStatus", opts)
}

func (a *App) PingTimer() {
	runtime.EventsEmit(a.ctx, "evTimer", time.Now().UnixMilli())
}

func TestEncodeDecodeNEventTestEncodeDecodeNEvent(t *zerolog.Event) string {
	nevent, err := nip19.EncodeEvent(
		"45326f5d6962ab1e3cd424e758c3002b8665f7b0d8dcee9fe9e288d7751ac194",
		[]string{"wss://banana.com"},
		"7fa56f5d6962ab1e3cd424e758c3002b8665f7b0d8dcee9fe9e288d7751abb88",
	)
	if err != nil {
		t.Msgf("shouldn't error: %s", err)
	}

	prefix, res, err := nip19.Decode(nevent)
	if err != nil {
		t.Msgf("shouldn't error: %s", err)
	}

	if prefix != "nevent" {
		t.Msgf("should have 'nevent' prefix, not '%s'", prefix)
	}

	ep, ok := res.(nostr.EventPointer)
	if !ok {
		t.Msgf("'%s' should be an nevent, not %v", nevent, res)
	}

	if ep.Author != "7fa56f5d6962ab1e3cd424e758c3002b8665f7b0d8dcee9fe9e288d7751abb88" {
		t.Msgf("wrong author")
	}

	if ep.ID != "45326f5d6962ab1e3cd424e758c3002b8665f7b0d8dcee9fe9e288d7751ac194" {
		t.Msgf("wrong id")
	}

	if len(ep.Relays) != 1 || ep.Relays[0] != "wss://banana.com" {
		t.Msgf("wrong relay")
	}
	return nevent
}
