package main

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type NetappCollector struct {
	AggrCollector   *AggrCollector
	VolumeCollector *VolumeCollector
	scrapeFailure   prometheus.Counter
	scrapeCounter   prometheus.Counter
}

func NewNetappCollector(filer NetappFilerClient) NetappCollector {
	return NetappCollector{
		AggrCollector: &AggrCollector{
			Filer: filer,
			ApiCollectorBase: ApiCollectorBase{
				maxAge: 5 * time.Minute,
			},
		},
		VolumeCollector: &VolumeCollector{
			Filer: filer,
			ApiCollectorBase: ApiCollectorBase{
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
	n.VolumeCollector.Describe(ch)
	n.AggrCollector.Describe(ch)
}

func (n NetappCollector) Collect(ch chan<- prometheus.Metric) {
	logger.Debug("calling Collect()")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go n.fetchAndCollect(n.VolumeCollector, ch, wg)
	go n.fetchAndCollect(n.AggrCollector, ch, wg)
	wg.Wait()

	ch <- n.scrapeFailure
	ch <- n.scrapeCounter
}

func (n NetappCollector) fetchAndCollect(m ApiCollector, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()

	// fetch() makes expensive http request to netapp's ONTAP system. The success channel is closed when
	// request is returned successfully, otherwise fail channel is closed.
	success := make(chan bool)
	fail := make(chan bool)
	go n.fetch(m, success, fail)

	// Since fetch() is called in go routine, metrics can be exported right away, when data is recent enough.
	m.Lock()
	if m.IsDataFresh() {
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

func (n NetappCollector) fetch(m ApiCollector, success, fail chan<- bool) {
	data, err := m.Fetch()
	if err != nil {
		logger.Error(err)
		n.scrapeFailure.Inc()
		close(fail)
	} else {
		m.Lock()
		err = m.SaveData(data)
		if err != nil {
			m.Unlock()
			logger.Error(err)
			close(fail)
		}
		m.SetFetchTime(time.Now())
		m.Unlock()
		close(success)
	}
	n.scrapeCounter.Inc()
}
