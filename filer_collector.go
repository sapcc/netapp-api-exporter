package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type volumeMetrics []struct {
	desc    *prometheus.Desc
	valType prometheus.ValueType
	evalFn  func(volume *NetappVolume) float64
}

type aggregateMetrics []struct {
	desc    *prometheus.Desc
	valType prometheus.ValueType
	evalFn  func(agg *NetappAggregate) float64
}

type FilerCollector struct {
	*FilerManager
	up                prometheus.Gauge
	lastScrapeSuccess prometheus.Gauge
	scrapesTotal      prometheus.Counter
	scrapesFailure    prometheus.Counter
}

var (
	volumeLabels = []string{
		"vserver",
		"volume",
		"project_id",
		"share_id",
	}
	aggregateLabels = []string{
		"node",
		"aggregate",
	}

	volMetrics = volumeMetrics{
		{
			desc: prometheus.NewDesc(
				"netapp_volume_total_size_bytes",
				"Netapp Volume Metrics: total size",
				volumeLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(v *NetappVolume) float64 { return v.SizeTotal },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_used_bytes",
				"Netapp Volume Metrics: used size",
				volumeLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(v *NetappVolume) float64 { return v.SizeUsed },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_available_bytes",
				"Netapp Volume Metrics: available size",
				volumeLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(v *NetappVolume) float64 { return v.SizeAvailable },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_snapshot_used_bytes",
				"Netapp Volume Metrics: size used by snapshots",
				volumeLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(v *NetappVolume) float64 { return v.SizeUsedBySnapshots },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_snapshot_available_bytes",
				"Netapp Volume Metrics: size available for snapshots",
				volumeLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(v *NetappVolume) float64 { return v.SizeAvailableForSnapshots },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_snapshot_reserved_bytes",
				"Netapp Volume Metrics: size reserved for snapshots",
				volumeLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(v *NetappVolume) float64 { return v.SnapshotReserveSize },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_used_percentage",
				"Netapp Volume Metrics: used percentage ",
				volumeLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(v *NetappVolume) float64 { return v.PercentageSizeUsed },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_saved_percentage",
				"Netapp Volume Metrics: percentage of space compression and deduplication saved",
				volumeLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(v *NetappVolume) float64 { return v.PercentageTotalSpaceSaved },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_compression_saved_percentage",
				"Netapp Volume Metrics: percentage of space compression saved",
				volumeLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(v *NetappVolume) float64 { return v.PercentageCompressionSpaceSaved },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_deduplication_saved_percentage",
				"Netapp Volume Metrics: percentage of space deduplication saved",
				volumeLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(v *NetappVolume) float64 { return v.PercentageDeduplicationSpaceSaved },
		},
	}

	aggMetrics = aggregateMetrics{
		{
			desc: prometheus.NewDesc(
				"netapp_aggregate_total_size_bytes",
				"Netapp Aggregate Metrics: total size",
				aggregateLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(m *NetappAggregate) float64 { return m.SizeTotal },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_available_bytes",
				"Netapp Aggregate Metrics: available size",
				aggregateLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(m *NetappAggregate) float64 { return m.SizeAvailable },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_used_bytes",
				"Netapp Aggregate Metrics: used size",
				aggregateLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(m *NetappAggregate) float64 { return m.SizeUsed },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_used_percentage",
				"Netapp Aggregate Metrics: used percentage",
				aggregateLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(m *NetappAggregate) float64 { return m.PercentUsedCapacity },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_physical_used_bytes",
				"Netapp Aggregate Metrics: physical used size",
				aggregateLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(m *NetappAggregate) float64 { return m.PhysicalUsed },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_physical_percentage",
				"Netapp Aggregate Metrics: physical used percentage",
				aggregateLabels,
				nil),
			valType: prometheus.GaugeValue,
			evalFn:  func(m *NetappAggregate) float64 { return m.PhysicalUsedPercent },
		},
	}
)

func RegisterFilerCollector(filer *FilerManager, reg prometheus.Registerer) {
	cc := &FilerCollector{
		FilerManager: filer,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "netapp",
			Subsystem: "exporter",
			Name:      "up",
			Help:      "'1' if the last scrape of filer was successful, '0' otherwise.",
		}),
		lastScrapeSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "netapp",
			Subsystem: "exporter",
			Name:      "scrape_last_timestamp",
			Help:      "Timestamp of the last successful scrape.",
		}),
		scrapesTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "netapp",
			Subsystem: "exporter",
			Name:      "scrape_total",
			Help:      "The total number of filer scrapes.",
		}),
		scrapesFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "netapp",
			Subsystem: "exporter",
			Name:      "scrape_failure",
			Help:      "The number of filer scrape failures.",
		}),
	}
	labels := prometheus.Labels{
		"filer":             filer.Name,
		"availability_zone": filer.AvailabilityZone,
	}
	prometheus.WrapRegistererWith(labels, reg).MustRegister(cc)
}

func (c FilerCollector) Describe(ch chan<- *prometheus.Desc) {
	logger.Debug("calling Describe()")
	ch <- c.up.Desc()
	ch <- c.scrapesTotal.Desc()
	ch <- c.scrapesFailure.Desc()
	ch <- c.lastScrapeSuccess.Desc()
	for _, v := range volMetrics {
		ch <- v.desc
	}
	for _, v := range aggMetrics {
		ch <- v.desc
	}
}

func (c FilerCollector) Collect(ch chan<- prometheus.Metric) {
	logger.Debug("calling Collect()")
	defer func() {
		ch <- c.scrapesTotal
		ch <- c.scrapesFailure
		ch <- c.lastScrapeSuccess
		ch <- c.up
	}()

	c.volScrapeAndCollect(ch)
	c.aggScrapeAndCollect(ch)
}

func (c FilerCollector) volScrapeAndCollect(ch chan<- prometheus.Metric) {
	var (
		err  error
		vols []*NetappVolume
		done = make(chan struct{})
	)

	// Scrape Filer metrics.
	go func() {
		vols, err = c.GetNetappVolume()
		if err != nil {
			c.up.Set(0)
			c.scrapesTotal.Inc()
			c.scrapesFailure.Inc()
		} else {
			c.up.Set(1)
			c.scrapesTotal.Inc()
			c.lastScrapeSuccess.SetToCurrentTime()
			c.mtxVol.Lock()
			c.Volumes = vols
			c.lastVolScrape = time.Now()
			c.mtxVol.Unlock()
		}
		done <- struct{}{}
	}()

	// Last metrics are recent enough. Export them immediately.
	c.mtxVol.Lock()
	if time.Since(c.lastVolScrape) < c.volMaxAge {
		volCollect(c.Volumes, ch)
		c.mtxVol.Unlock()
		return
	}
	c.mtxVol.Unlock()

	// Last metrics are not recent. Wait for scraping finish.
	<-done

	// No error in scraping. Export Filer metrics.
	if err == nil {
		c.mtxVol.Lock()
		volCollect(c.Volumes, ch)
		c.mtxVol.Unlock()
	} else {
		logger.Error(err)
	}
}

func volCollect(volumes []*NetappVolume, ch chan<- prometheus.Metric) {
	for _, v := range volumes {
		labels := []string{v.Vserver, v.Volume, v.ProjectID, v.ShareID}
		for _, m := range volMetrics {
			ch <- prometheus.MustNewConstMetric(m.desc, m.valType, m.evalFn(v), labels...)
		}
	}
}

func (c FilerCollector) aggScrapeAndCollect(ch chan<- prometheus.Metric) {
	var (
		err  error
		aggs []*NetappAggregate
		done = make(chan struct{})
	)

	// Scrape Filer metrics.
	go func() {
		aggs, err = c.GetNetappAggregate()
		if err != nil {
			c.up.Set(0)
			c.scrapesTotal.Inc()
			c.scrapesFailure.Inc()
		} else {
			c.up.Set(1)
			c.scrapesTotal.Inc()
			c.lastScrapeSuccess.SetToCurrentTime()
			c.mtxAgg.Lock()
			c.Aggregates = aggs
			c.lastAggScrape = time.Now()
			c.mtxAgg.Unlock()
		}
		done <- struct{}{}
	}()

	// Last metrics are recent enough. Export them immediately.
	c.mtxAgg.Lock()
	if time.Since(c.lastAggScrape) < c.aggMaxAge {
		aggCollect(c.Aggregates, ch)
		c.mtxAgg.Unlock()
		return
	}
	c.mtxAgg.Unlock()

	// Last metrics are not recent. Wait for scraping finish.
	<-done

	// No error in scraping. Export Filer metrics.
	if err == nil {
		c.mtxAgg.Lock()
		aggCollect(c.Aggregates, ch)
		c.mtxAgg.Unlock()
	} else {
		logger.Error(err)
	}
}

func aggCollect(aggregates []*NetappAggregate, ch chan<- prometheus.Metric) {
	for _, v := range aggregates {
		labels := []string{v.OwnerName, v.Name}
		for _, m := range aggMetrics {
			ch <- prometheus.MustNewConstMetric(m.desc, m.valType, m.evalFn(v), labels...)
		}
	}
}
