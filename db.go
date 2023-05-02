package main

import (
	"fmt"
	"github.com/arriqaaq/hash"
	"github.com/nbd-wtf/go-nostr"
	"sync"
)

type DB struct {
	cache *hash.Hash
	mu    sync.Mutex
}

const (
	EVENT = "event"
	META  = "meta"
)

func NewDB() DB {
	return DB{
		cache: hash.New(),
	}
}

func (p *DB) GetLock() {
	p.mu.Lock()
}
func (p *DB) ReleaseLock() {
	p.mu.Unlock()
}

func (p *DB) HasEvent(evId string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cache.HExists(EVENT, evId)
}

func (p *DB) HasProfile(pk string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cache.HExists(META, pk)
}

func (p *DB) GetEvent(evId string) *nostr.Event {
	p.mu.Lock()
	defer p.mu.Unlock()
	r := p.cache.HGet(EVENT, evId)
	if r == nil {
		return nil
	}
	return r.(*nostr.Event)
}

func (p *DB) GetProfile(pk string) *Profile {
	p.mu.Lock()
	defer p.mu.Unlock()
	r := p.cache.HGet(META, pk)
	if r == nil {
		return nil
	}
	return r.(*Profile)
}

func (p *DB) AddProfile(pk string, profile *Profile) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cache.HSet(META, pk, profile)
}

func (p *DB) AddEvent(evId string, event *nostr.Event) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cache.HSet(EVENT, evId, event)
}

func (p *DB) DumpEvents() {
	p.mu.Lock()
	defer p.mu.Unlock()
	//entries := p.cache.HGetAll(EVENT)
	keys := p.cache.HKeys(EVENT)

	fmt.Println("Dumping", len(keys), "keys")
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		fmt.Println(key)
	}
}
