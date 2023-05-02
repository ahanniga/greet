package main

import (
	"github.com/nbd-wtf/go-nostr"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type RelayPool struct {
	pool []*Relay
}

func NewRelayPool() *RelayPool {
	return &RelayPool{
		pool: []*Relay{},
	}
}

func (p *RelayPool) Add(relay *Relay) error {
	log.Debug().Msgf("Adding relay %s to pool", relay.Url)
	if relay.Enabled {
		err := relay.Connect()
		if err != nil {
			return err
		}
	}
	p.pool = append(p.pool, relay)
	return nil
}

func (p *RelayPool) AddAll(relays []*Relay) {
	for _, r := range relays {
		err := p.Add(r)
		if err != nil {
			log.Error().Msgf("Could not add %s to pool: %s", r.Url, err.Error())
		}
	}
}

func (p *RelayPool) Remove(relay *Relay) {
	relay.Disconnect()
	for i, r := range p.pool {
		if r.Url == relay.Url {
			log.Info().Msgf("Removing relay %s from pool", r.Url)
			p.pool = slices.Delete(p.pool, i, i+1)
		}
	}
}

func (p *RelayPool) RemoveAll() {
	p.pool = nil
	p.pool = []*Relay{}
}

func (p *RelayPool) Query(filter *nostr.Filter, evts chan *nostr.Event) int {
	log.Debug().Msgf("Query filter %s", filter)
	pending := 0
	for _, r := range p.pool {
		if r.Enabled && r.Read {
			log.Debug().Msgf("Enabled && read: %s", r.Url)
			pending++
			go func(relay *Relay, ch chan *nostr.Event) {
				evs, err := relay.SubscribeWithTimeout(filter)
				if err != nil {
					log.Error().Msgf("Channel subscribe error: %s", err.Error())
				} else {
					for _, ev := range evs {
						ev.SetExtra("relay", relay.Url)
						log.Trace().Msgf("Event kind %d from relay %s:", ev.Kind, relay.Url, ev)
						ch <- ev
					}
				}
				log.Debug().Msgf("Received %d events from relay %s", len(evs), relay.Url)
				ch <- nil
			}(r, evts)
		}
	}
	log.Debug().Msgf("Returning %d pending relay queries", pending)
	return pending
}

func (p *RelayPool) ReconnectAll() {
	for _, r := range p.pool {
		err := r.Disconnect()
		if err != nil {
			log.Err(err)
		}
	}
	for _, r := range p.pool {
		err := r.Connect()
		if err != nil {
			log.Err(err)
		}
	}
}

func (p *RelayPool) DisconnectAll() {
	for _, r := range p.pool {
		err := r.Disconnect()
		if err != nil {
			log.Err(err)
		}
	}
}

func (p *RelayPool) GetRelayByUrl(url string) *Relay {
	for _, r := range p.pool {
		if r.Url == url {
			return r
		}
	}
	return nil
}
