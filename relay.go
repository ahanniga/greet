package main

import (
	"context"
	"github.com/nbd-wtf/go-nostr"
	"github.com/rs/zerolog/log"
)

type RelayStruct struct {
	Url       string `json:"url"`
	Read      bool   `json:"read"`
	Write     bool   `json:"write"`
	Enabled   bool   `json:"enabled"`
	conn      *nostr.Relay
	subs      []*nostr.Subscription
	relayMeta *RelayMetadata
}

func NewRelay() *RelayStruct {
	return &RelayStruct{
		Url:     "",
		Read:    false,
		Write:   false,
		Enabled: false,
		subs:    []*nostr.Subscription{},
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

func (r *RelayStruct) Connect(ctx context.Context) error {
	conn, err := nostr.RelayConnect(ctx, r.Url)
	if err != nil {
		return err
	}
	log.Debug().Msgf("Successful connection to %s", r.Url)
	r.conn = conn
	return nil
}

func (r *RelayStruct) RemoveSub(id string) {
	for i, s := range r.subs {
		if s.GetID() == id {
			r.subs = append(r.subs[:i], r.subs[i+1:]...)
		}
	}
}
