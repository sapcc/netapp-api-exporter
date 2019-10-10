package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type FilerCollector struct {
	Filer
	aggrManager    *AggrManager
	volManager     *VolumeManager
	scrapesFailure prometheus.Counter
}

func NewFilerCollector(filer Filer) *FilerCollector {
	return &FilerCollector{
		Filer:       filer,
		aggrManager: &AggrManager{maxAge: 5 * time.Minute},
		volManager:  &VolumeManager{maxAge: 5 * time.Minute},
		scrapesFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "netapp",
			Subsystem: "exporter",
			Name:      "scrape_failure",
			Help:      "The number of filer scrape failures.",
		}),
	}
}

func (c FilerCollector) Describe(ch chan<- *prometheus.Desc) {
	logger.Debug("calling Describe()")
	ch <- c.scrapesFailure.Desc()
	c.volManager.Describe(ch)
	c.aggrManager.Describe(ch)
}

func (c FilerCollector) Collect(ch chan<- prometheus.Metric) {
	logger.Debug("calling Collect()")
	ch <- c.scrapesFailure
	c.CollectAggr(ch)
	c.CollectVolume(ch)
}

func (c FilerCollector) CollectAggr(ch chan<- prometheus.Metric) {
	var (
		aggrs       []*NetappAggregate
		err         error
		doneFetch   = make(chan bool)
		doneCollect = make(chan bool)
	)

	go func() {
		defer close(doneFetch)
		aggrs, err = c.aggrManager.Fetch(c.Filer)
		if err != nil {
			logger.Error(err)
			c.scrapesFailure.Inc()
		} else {
			c.aggrManager.mtx.Lock()
			c.aggrManager.lastFetchTime = time.Now()
			c.aggrManager.Aggregates = aggrs
			c.aggrManager.mtx.Unlock()
		}
	}()

	go func() {
		c.aggrManager.mtx.Lock()
		defer c.aggrManager.mtx.Unlock()
		if time.Since(c.aggrManager.lastFetchTime) < c.aggrManager.maxAge {
			c.aggrManager.Collect(ch)
			close(doneCollect)
		}
	}()

	select {
	case <-doneFetch:
		if err == nil {
			c.aggrManager.mtx.Lock()
			c.aggrManager.Collect(ch)
			c.aggrManager.mtx.Unlock()
		}
	case <-doneCollect:
	}
}

func (c FilerCollector) CollectVolume(ch chan<- prometheus.Metric) {
	var (
		vols        []*NetappVolume
		err         error
		doneFetch   = make(chan bool)
		doneCollect = make(chan bool)
	)

	go func() {
		defer close(doneFetch)
		vols, err = c.volManager.Fetch(c.Filer)
		if err != nil {
			logger.Error(err)
			c.scrapesFailure.Inc()
		} else {
			c.volManager.mtx.Lock()
			c.volManager.lastFetchTime = time.Now()
			c.volManager.Volumes = vols
			c.volManager.mtx.Unlock()
		}
	}()

	go func() {
		c.volManager.mtx.Lock()
		defer c.volManager.mtx.Unlock()
		if time.Since(c.volManager.lastFetchTime) < c.volManager.maxAge {
			c.volManager.Collect(ch)
			close(doneCollect)
		}
	}()

	select {
	case <-doneFetch:
		if err == nil {
			c.volManager.mtx.Lock()
			c.volManager.Collect(ch)
			c.volManager.mtx.Unlock()
		}
	case <-doneCollect:
	}
}

//
//func (c FilerCollector) aggScrapeAndCollect(ch chan<- prometheus.Metric) {
//	var (
//		err  error
//		aggs []*NetappAggregate
//		done = make(chan struct{})
//	)
//
//	// Scrape FilerBase metrics.
//	go func() {
//		aggs, err = c.GetNetappAggregate()
//		if err != nil {
//			c.up.Set(0)
//			c.scrapesTotal.Inc()
//			c.scrapesFailure.Inc()
//		} else {
//			c.up.Set(1)
//			c.scrapesTotal.Inc()
//			c.lastScrapeSuccess.SetToCurrentTime()
//			c.mtxAgg.Lock()
//			c.Aggregates = aggs
//			c.lastAggScrape = time.Now()
//			c.mtxAgg.Unlock()
//		}
//		done <- struct{}{}
//	}()
//
//	// Last metrics are recent enough. Export them immediately.
//	c.mtxAgg.Lock()
//	if time.Since(c.lastAggScrape) < c.aggMaxAge {
//		aggCollect(c.Aggregates, ch)
//		c.mtxAgg.Unlock()
//		return
//	}
//	c.mtxAgg.Unlock()
//
//	// Last metrics are not recent. Wait for scraping finish.
//	<-done
//
//	// No error in scraping. Export FilerBase metrics.
//	if err == nil {
//		c.mtxAgg.Lock()
//		aggCollect(c.Aggregates, ch)
//		c.mtxAgg.Unlock()
//	} else {
//		logger.Error(err)
//	}
//}

//func (c FilerCollector) volScrapeAndCollect(ch chan<- prometheus.Metric) {
//	var (
//		err  error
//		vols []*NetappVolume
//		done = make(chan struct{})
//	)
//
//	// Scrape FilerBase metrics.
//	go func() {
//		vols, err = c.GetNetappVolume()
//		if err != nil {
//			c.up.Set(0)
//			c.scrapesTotal.Inc()
//			c.scrapesFailure.Inc()
//		} else {
//			c.up.Set(1)
//			c.scrapesTotal.Inc()
//			c.lastScrapeSuccess.SetToCurrentTime()
//			c.mtxVol.Lock()
//			c.Volumes = vols
//			c.lastVolScrape = time.Now()
//			c.mtxVol.Unlock()
//		}
//		done <- struct{}{}
//	}()
//
//	// Last metrics are recent enough. Export them immediately.
//	c.mtxVol.Lock()
//	if time.Since(c.lastVolScrape) < c.volMaxAge {
//		volCollect(c.Volumes, ch)
//		c.mtxVol.Unlock()
//		return
//	}
//	c.mtxVol.Unlock()
//
//	// Last metrics are not recent. Wait for scraping finish.
//	<-done
//
//	// No error in scraping. Export FilerBase metrics.
//	if err == nil {
//		c.mtxVol.Lock()
//		volCollect(c.Volumes, ch)
//		c.mtxVol.Unlock()
//	} else {
//		logger.Error(err)
//	}
//}
