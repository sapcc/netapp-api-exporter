package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

type Collector struct {
	Filer
	aggrManager    *AggrManager
	volManager     *VolumeManager
	scrapesFailure prometheus.Counter
}

func NewCollector(filer Filer) *Collector {
	return &Collector{
		Filer:       filer,
		aggrManager: &AggrManager{maxAge: 5 * time.Minute},
		volManager:  &VolumeManager{maxAge: 5 * time.Minute},
		scrapesFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "netapp",
			Subsystem: "filer",
			Name:      "scrape_failure",
			Help:      "The number of scraping failures of filer.",
		}),
	}
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	logger.Debug("calling Describe()")
	ch <- c.scrapesFailure.Desc()
	c.volManager.Describe(ch)
	c.aggrManager.Describe(ch)
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	logger.Debug("calling Collect()")
	ch <- c.scrapesFailure

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go c.CollectVolume(ch, wg)
	go c.CollectAggr(ch, wg)
	wg.Wait()
}

func (c Collector) CollectAggr(ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	var (
		fail = make(chan bool)
		done = make(chan bool)
	)

	defer wg.Done()

	// Fetch data concurrently.
	go func() {
		aggrs, err := c.aggrManager.Fetch(c.Filer)
		if err != nil {
			logger.Error(err)
			c.scrapesFailure.Inc()
			close(fail)
		} else {
			c.aggrManager.Lock()
			c.aggrManager.lastFetchTime = time.Now()
			c.aggrManager.Aggregates = aggrs
			c.aggrManager.Unlock()
			close(done)
		}
	}()

	// Cached data are recent enough. Collect and return.
	c.aggrManager.Lock()
	if time.Since(c.aggrManager.lastFetchTime) < c.aggrManager.maxAge {
		c.aggrManager.Collect(ch)
		c.aggrManager.Unlock()
		return
	}

	// Cached data are not recent. Wait for fetch.
	c.aggrManager.Unlock()
	select {
	case <-done:
		c.aggrManager.Lock()
		c.aggrManager.Collect(ch)
		c.aggrManager.Unlock()
	case <-fail:
	}
	return
}

func (c Collector) CollectVolume(ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	var (
		fail = make(chan bool)
		done = make(chan bool)
	)

	defer wg.Done()

	// Fetch data concurrently.
	go func() {
		vols, err := c.volManager.Fetch(c.Filer)
		if err != nil {
			logger.Error(err)
			c.scrapesFailure.Inc()
			close(fail)
		} else {
			c.volManager.Lock()
			c.volManager.lastFetchTime = time.Now()
			c.volManager.Volumes = vols
			c.volManager.Unlock()
			close(done)
		}
	}()

	// Cached data are recent enough. Collect and return.
	c.volManager.Lock()
	if time.Since(c.volManager.lastFetchTime) < c.volManager.maxAge {
		c.volManager.Collect(ch)
		c.volManager.Unlock()
		return
	}

	// Cached data are not recent. Wait for fetch.
	c.volManager.Unlock()
	select {
	case <-done:
		c.volManager.Lock()
		c.volManager.Collect(ch)
		c.volManager.Unlock()
	case <-fail:
	}
	return
}
