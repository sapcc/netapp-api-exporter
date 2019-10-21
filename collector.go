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
	AggrManager   *AggrManager
	VolumeManager *VolumeManager
	scrapeFailure prometheus.Counter
	scrapeCounter prometheus.Counter
}

func NewNetappCollector(filer Filer) NetappCollector {
	return NetappCollector{
		AggrManager: &AggrManager{
			Manager: Manager{
				filer:  filer,
				maxAge: 5 * time.Minute,
			},
		},
		VolumeManager: &VolumeManager{
			Manager: Manager{
				filer:  filer,
				maxAge: 5 * time.Minute,
			},
		},
		scrapeFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "netapp",
			Subsystem: "filer",
			Name:      "scrape_failure",
			Help:      "The number of scraping failures of filer.",
		}),
		scrapeCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "netapp",
			Subsystem: "filer",
			Name:      "scrape_counter",
			Help:      "The number of scrapes.",
		}),
	}
}

func (n NetappCollector) Describe(ch chan<- *prometheus.Desc) {
	logger.Debug("calling Describe()")
	ch <- n.scrapeFailure.Desc()
	ch <- n.scrapeCounter.Desc()
	n.VolumeManager.Describe(ch)
	n.AggrManager.Describe(ch)
}

func (n NetappCollector) Collect(ch chan<- prometheus.Metric) {
	logger.Debug("calling Collect()")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go n.collectManager(n.VolumeManager, ch, wg)
	go n.collectManager(n.AggrManager, ch, wg)
	wg.Wait()

	ch <- n.scrapeFailure
	ch <- n.scrapeCounter
}

func (n NetappCollector) collectManager(m ManagerCollector, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()

	// fetch() makes expensive http request to netapp's ONTAP system. The success channel is closed when
	// request is returned successfully, otherwise fail channel is closed.
	success := make(chan bool)
	fail := make(chan bool)
	go n.fetch(m, success, fail)

	// Since fetch() is called in go routine, metrics can be exported right away, when data is recent enough.
	m.Lock()
	if time.Since(m.LastFetchTime()) < m.MaxAge() {
		m.Collect(ch)
		m.Unlock()
		return
	}

	// Data is not recent, and we have to wait for fetch() to finish.
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

func (n NetappCollector) fetch(m ManagerCollector, success, fail chan<- bool) {
	data, err := m.Fetch()
	if err != nil {
		logger.Error(err)
		n.scrapeFailure.Inc()
		close(fail)
	} else {
		m.SaveDataWithTime(data, time.Now())
		close(success)
	}
	n.scrapeCounter.Inc()
}
