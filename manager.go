package main

import (
	"sync"
	"time"
)

type Manager struct {
	sync.Mutex
	filer         *Filer
	lastFetchTime time.Time
	maxAge        time.Duration
}

func (m Manager) MaxAge() time.Duration {
	return m.maxAge
}

func (m Manager) LastFetchTime() time.Time {
	return m.lastFetchTime
}
