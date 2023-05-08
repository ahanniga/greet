package main

import (
	"context"
	"github.com/nbd-wtf/go-nostr"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type RelayPool struct {
	pool    []*RelayStruct
	rootCtx context.Context
}

func NewRelayPool() *RelayPool {
	return &RelayPool{
		pool:    []*RelayStruct{},
		rootCtx: context.Background(),
	}
}

func (p *RelayPool) UnsubscribeAll() {
	for _, relay := range p.pool {
		for _, sub := range relay.subs {
			log.Debug().Msgf("Unsubscribing from relay %s, ID %s", relay.Url, sub.GetID())
			sub.Unsub()
		}
		relay.subs = []*nostr.Subscription{}
	}
}

func (p *RelayPool) Add(relay *RelayStruct) error {
	if relay.Enabled {
		log.Debug().Msgf("Adding relay %s to pool", relay.Url)
		err := relay.Connect(p.rootCtx)
		if err != nil {
			return err
		}
		p.pool = append(p.pool, relay)
	}

	return nil
}

func (p *RelayPool) QuerySync(f *nostr.Filter, c chan *nostr.Event) {
	wg := sync.WaitGroup{}
	for _, relay := range p.pool {
		if relay.Enabled && relay.Read {
			wg.Add(1)
			go func(r *RelayStruct) {
				result, err := r.conn.QuerySync(p.rootCtx, *f)
				if err != nil {
					log.Error().Msgf("QuerySync error from %s: %s", r.Url, err.Error())

					// Try a reconnect, but disable relay if it fails
					time.Sleep(time.Second * 2)
					relay.Connect(p.rootCtx)
					result, err = r.conn.QuerySync(p.rootCtx, *f)
					if err != nil {
						log.Error().Msgf("QuerySync (retry) error from %s: %s: disabled", r.Url, err.Error())
						r.Enabled = false
						return
					}
				}
				for i := 0; i < len(result); i++ {
					ev := result[i]
					ev.SetExtra("relay", r.Url)
					c <- ev
				}
				wg.Done()
			}(relay)
		}
	}
	wg.Wait()
	close(c)
}

func (p *RelayPool) Subscribe(f *nostr.Filter, c chan *nostr.Event, ac chan *nostr.Event) {
	for _, relay := range p.pool {
		if relay.Enabled && relay.Read {
			go func(r *RelayStruct) {

				gotEose := false
				sub, err := r.conn.Subscribe(p.rootCtx, []nostr.Filter{*f})

				if err != nil {
					log.Error().Msgf(err.Error())
					return
				}
				r.subs = append(r.subs, sub)
				log.Debug().Msgf("Subscribed to relay %s", r.Url)
				defer r.RemoveSub(sub.GetID())

				for {
					select {
					case ev := <-sub.Events:
						if ev == nil {
							return
						}
						log.Trace().Msgf("Got event from %s %s", r.Url, ev.ID)
						ev.SetExtra("relay", r.Url)
						if gotEose {
							ac <- ev
						} else {
							c <- ev
						}
					case <-sub.EndOfStoredEvents:
						log.Debug().Msgf("Got EOSE from %s", r.Url)
						gotEose = true
					case <-sub.Context.Done():
						log.Debug().Msgf("Subscription to relay %s completed: ", r.Url)
						return
					}
				}
			}(relay)
		}
	}
}

func (p *RelayPool) AddAll(relays []*RelayStruct) {
	for _, r := range relays {
		err := p.Add(r)
		if err != nil {
			log.Error().Msgf("Could not add %s to pool: %s", r.Url, err.Error())
		}
	}
}

func (p *RelayPool) RemoveAll() {
	p.DisconnectAll()
	p.pool = []*RelayStruct{}
}

func (p *RelayPool) DisconnectAll() {
	p.UnsubscribeAll()
	for _, r := range p.pool {
		if r.Enabled {
			log.Debug().Msgf("Closing connection to relay %s", r.Url)
			err := r.conn.Close()
			if err != nil {
				log.Err(err)
			}
		}
	}
}

func (p *RelayPool) GetRelayByUrl(url string) *RelayStruct {
	for _, r := range p.pool {
		if r.Url == url {
			return r
		}
	}
	return nil
}
