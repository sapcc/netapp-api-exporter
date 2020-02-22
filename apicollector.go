package main

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ApiCollector interface {
	sync.Locker
	prometheus.Collector
	// fetch data from api endpoint
	Fetch() (data []interface{}, err error)
	// save data in collector
	SaveData(data []interface{}) error
	SetFetchTime(time.Time)
	IsDataFresh() bool
}

type ApiCollectorBase struct {
	sync.Mutex
	lastFetchTime time.Time
	maxAge        time.Duration
}

func (m *ApiCollectorBase) IsDataFresh() bool {
	return time.Since(m.lastFetchTime) < m.maxAge
}

func (m *ApiCollectorBase) SetFetchTime(t time.Time) {
	m.lastFetchTime = t
}
