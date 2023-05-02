package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nbd-wtf/go-nostr"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Relay struct {
	Url       string `json:"url"`
	Read      bool   `json:"read"`
	Write     bool   `json:"write"`
	Enabled   bool   `json:"enabled"`
	conn      *nostr.Relay
	sub       *nostr.Subscription
	relayMeta *RelayMetadata
	mu        sync.Mutex
	rootCtx   context.Context
}

func NewRelay() *Relay {
	return &Relay{
		Url:     "",
		Read:    false,
		Write:   false,
		Enabled: false,
		relayMeta: &RelayMetadata{
			Name:                   "",
			Description:            "",
			Pubkey:                 "",
			Contact:                "",
			SupportedNips:          nil,
			SupportedNipExtensions: nil,
			Software:               "",
			Version:                "",
			Limitation: struct {
				MaxMessageLength int  `json:"max_message_length"`
				MaxSubscriptions int  `json:"max_subscriptions"`
				MaxFilters       int  `json:"max_filters"`
				MaxLimit         int  `json:"max_limit"`
				MaxSubidLength   int  `json:"max_subid_length"`
				MinPrefix        int  `json:"min_prefix"`
				MaxEventTags     int  `json:"max_event_tags"`
				MaxContentLength int  `json:"max_content_length"`
				MinPowDifficulty int  `json:"min_pow_difficulty"`
				AuthRequired     bool `json:"auth_required"`
				PaymentRequired  bool `json:"payment_required"`
			}{},
			PaymentsUrl: "",
			Fees: struct {
				Admission []struct {
					Amount int    `json:"amount"`
					Unit   string `json:"unit"`
				} `json:"admission"`
				Publication []interface{} `json:"publication"`
			}{},
		},
	}
}

type RelayMetadata struct {
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	Pubkey                 string   `json:"pubkey"`
	Contact                string   `json:"contact"`
	SupportedNips          []int    `json:"supported_nips"`
	SupportedNipExtensions []string `json:"supported_nip_extensions"`
	Software               string   `json:"software"`
	Version                string   `json:"version"`
	Limitation             struct {
		MaxMessageLength int  `json:"max_message_length"`
		MaxSubscriptions int  `json:"max_subscriptions"`
		MaxFilters       int  `json:"max_filters"`
		MaxLimit         int  `json:"max_limit"`
		MaxSubidLength   int  `json:"max_subid_length"`
		MinPrefix        int  `json:"min_prefix"`
		MaxEventTags     int  `json:"max_event_tags"`
		MaxContentLength int  `json:"max_content_length"`
		MinPowDifficulty int  `json:"min_pow_difficulty"`
		AuthRequired     bool `json:"auth_required"`
		PaymentRequired  bool `json:"payment_required"`
	} `json:"limitation"`
	PaymentsUrl string `json:"payments_url"`
	Fees        struct {
		Admission []struct {
			Amount int    `json:"amount"`
			Unit   string `json:"unit"`
		} `json:"admission"`
		Publication []interface{} `json:"publication"`
	} `json:"fees"`
}

var wsMatch = regexp.MustCompile("ws[s]?://")

func (r *Relay) Connect() error {
	endpoint, err := r.getHttpEndpoint(r.Url)
	if err != nil {
		return err
	}

	// Get capabilities
	info, err := getRelayMetadata(endpoint)
	if err != nil {
		return err
	}
	r.relayMeta = info
	r.rootCtx = context.Background()
	r.conn, err = nostr.RelayConnect(r.rootCtx, r.Url)
	if err != nil {
		return err
	}
	return nil
}

func (r *Relay) Disconnect() error {
	log.Info().Msgf("Closing %s (%s)", r.Url, r.Enabled)
	if r.Enabled {
		r.conn.Close()
	}
	return nil
}

func (r *Relay) Reconnect() error {
	err := r.Disconnect()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	err = r.Connect()
	if err != nil {
		return err
	}
	return nil
}

func (r *Relay) Subscribe(filter *nostr.Filter, ctx context.Context) ([]*nostr.Event, error) {
	log.Debug().Msgf("Waiting subscription to %s", r.Url)
	return r.conn.QuerySync(ctx, *filter)
}

func (r *Relay) SubscribeWithTimeout(filter *nostr.Filter) ([]*nostr.Event, error) {
	ctx, cancel := context.WithTimeout(r.rootCtx, time.Duration(time.Second*5))
	defer cancel()

	res := make(chan []*nostr.Event)
	go func() {
		ar, err := r.Subscribe(filter, ctx)
		if err == nil {
			res <- ar
		}
		close(res)
	}()

	for {
		select {
		case aa := <-res:
			return aa, nil
		case <-ctx.Done():
			//log.Debug().Msgf("SubscribeWithTimeout: %s - disabling read for %s", ctx.Err(), r.Url)
			//r.Read = false
			return nil, ctx.Err()
		}
	}
	return []*nostr.Event{}, nil
}

func (r *Relay) getHttpEndpoint(url string) (string, error) {
	prefix := ""
	switch wsMatch.FindString(url) {
	case "ws://":
		prefix = "http://"
		break
	case "wss://":
		prefix = "https://"
		break
	default:
		return "", errors.New("Invalid websocket scheme:" + url)
	}
	return prefix + strings.Split(url, "//")[1], nil
}

func getRelayMetadata(url string) (*RelayMetadata, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/nostr+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	metadata := RelayMetadata{}
	err = json.Unmarshal([]byte(body), &metadata)
	if err != nil {
		return nil, err
	}

	// Check the default max subscriptions, default to 1
	if metadata.Limitation.MaxSubscriptions == 0 {
		metadata.Limitation.MaxSubscriptions = 1
	}
	return &metadata, nil
}
