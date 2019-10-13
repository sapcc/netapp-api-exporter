package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

type ManagerCollector interface {
	sync.Locker
	prometheus.Collector
	Fetch() (data []interface{}, err error)
	SaveDataWithTime(data []interface{}, t time.Time)
	LastFetchTime() time.Time
	MaxAge() time.Duration
}

type NetappCollector struct {
	aggrManager    *AggrManager
	volManager     *VolumeManager
	scrapesFailure prometheus.Counter
}

func NewNetappCollector(f *Filer) NetappCollector {
	return NetappCollector{
		aggrManager: &AggrManager{
			Mutex:         sync.Mutex{},
			filer:         f,
			Aggregates:    make([]*NetappAggregate, 0),
			lastFetchTime: time.Time{},
			maxAge:        5 * time.Minute,
		},
		volManager: &VolumeManager{
			Mutex:         sync.Mutex{},
			filer:         f,
			Volumes:       make([]*NetappVolume, 0),
			lastFetchTime: time.Time{},
			maxAge:        5 * time.Minute,
		},
		scrapesFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "netapp",
			Subsystem: "filer",
			Name:      "scrape_failure",
			Help:      "The number of scraping failures of filer.",
		}),
	}
}

func (fc NetappCollector) Describe(ch chan<- *prometheus.Desc) {
	logger.Debug("calling Describe()")
	ch <- fc.scrapesFailure.Desc()
	fc.volManager.Describe(ch)
	fc.aggrManager.Describe(ch)
}

func (fc NetappCollector) Collect(ch chan<- prometheus.Metric) {
	logger.Debug("calling Collect()")
	ch <- fc.scrapesFailure

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go fc.collectManager(fc.volManager, ch, wg)
	go fc.collectManager(fc.aggrManager, ch, wg)
	wg.Wait()
}

func (fc NetappCollector) collectManager(m ManagerCollector, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()

	success := make(chan bool)
	fail := make(chan bool)
	fc.fetch(m, success, fail)

	m.Lock()
	if time.Since(m.LastFetchTime()) < m.MaxAge() {
		m.Collect(ch)
		m.Unlock()
		return
	}

	m.Unlock()
	select {
	case <-success:
		m.Lock()
		m.Collect(ch)
		m.Unlock()
	case <-fail:
	}
	return
}

func (fc NetappCollector) fetch(m ManagerCollector, success, fail chan<- bool) {
	data, err := m.Fetch()
	if err != nil {
		logger.Error(err)
		fc.scrapesFailure.Inc()
		close(fail)
	} else {
		m.SaveDataWithTime(data, time.Now())
		close(success)
	}
}