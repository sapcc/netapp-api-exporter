package collector

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/netapp-api-exporter/pkg/netapp"
	log "github.com/sirupsen/logrus"
)

type VolumeCollector struct {
	filerName            string
	client               *netapp.Client
	volumes              []*netapp.Volume
	volumeMetrics        []VolumeMetric
	volumeTotalGauge     prometheus.Gauge
	scrapeCounter        prometheus.Counter
	scrapeFailureCounter prometheus.Counter
	scrapeDurationGauge  prometheus.Gauge
	mux                  sync.Mutex
	fetchPeriod          time.Duration
}

type VolumeMetric struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	getterFn  func(volume *netapp.Volume) float64
}

func NewVolumeCollector(filerName string, client *netapp.Client, fetchPeriod time.Duration) *VolumeCollector {
	volumeLabels := []string{"vserver", "volume", "volume_type", "project_id", "share_id", "share_name", "share_type"}
	volumeMetrics := []VolumeMetric{
		{
			desc: prometheus.NewDesc(
				"netapp_volume_state",
				"Netapp Volume Metrics: state (1: online; 2: restricted; 3: offline; 4: quiesced)",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return float64(v.State) },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_total_bytes",
				"Netapp Volume Metrics: total size",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return v.SizeTotal },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_used_bytes",
				"Netapp Volume Metrics: used size",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return v.SizeUsed },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_available_bytes",
				"Netapp Volume Metrics: available size",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return v.SizeAvailable },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_snapshot_used_bytes",
				"Netapp Volume Metrics: size used by snapshots",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return v.SizeUsedBySnapshots },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_snapshot_available_bytes",
				"Netapp Volume Metrics: size available for snapshots",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return v.SizeAvailableForSnapshots },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_snapshot_reserved_bytes",
				"Netapp Volume Metrics: size reserved for snapshots",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return v.SnapshotReserveSize },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_used_percentage",
				"Netapp Volume Metrics: used percentage ",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return v.PercentageSizeUsed },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_saved_total_percentage",
				"Netapp Volume Metrics: percentage of space compression and deduplication saved",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return v.PercentageTotalSpaceSaved },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_saved_compression_percentage",
				"Netapp Volume Metrics: percentage of space compression saved",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return v.PercentageCompressionSpaceSaved },
		}, {
			desc: prometheus.NewDesc(
				"netapp_volume_saved_deduplication_percentage",
				"Netapp Volume Metrics: percentage of space deduplication saved",
				volumeLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(v *netapp.Volume) float64 { return v.PercentageDeduplicationSpaceSaved },
		},
	}
	volumeTotalGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "netapp_volume_total",
			Help: "number of volumes scraped from Netapp filer",
		},
	)
	scrapeDurationGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "netapp_volume_scrape_duration_seconds",
			Help: "duration in seconds used to fetch volumes from filer",
		},
	)
	scrapeCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "netapp_volume_scrape_total",
			Help: "number of volume fetches from filer",
		},
	)
	scrapeFailureCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "netapp_volume_scrape_failure_total",
			Help: "number of failures for fetching volumes from filer",
		},
	)
	c := &VolumeCollector{
		filerName:            filerName,
		client:               client,
		fetchPeriod:          fetchPeriod,
		volumeMetrics:        volumeMetrics,
		volumeTotalGauge:     volumeTotalGauge,
		scrapeCounter:        scrapeCounter,
		scrapeFailureCounter: scrapeFailureCounter,
		scrapeDurationGauge:  scrapeDurationGauge,
	}
	go c.PeriodicFetch(nil)
	return c
}

func (c *VolumeCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.volumeMetrics {
		ch <- m.desc
	}
	ch <- c.volumeTotalGauge.Desc()
	ch <- c.scrapeCounter.Desc()
	ch <- c.scrapeFailureCounter.Desc()
	ch <- c.scrapeDurationGauge.Desc()
}

func (c *VolumeCollector) Collect(ch chan<- prometheus.Metric) {
	defer c.mux.Unlock()
	c.mux.Lock()

	// export metrics
	log.Debugf("VolumeCollector[%v] Collect() exporting %d volumes", c.filerName, len(c.volumes))
	for _, volume := range c.volumes {
		volumeLabels := []string{volume.Vserver, volume.Volume, volume.VolumeType, volume.ProjectID, volume.ShareID, volume.ShareName, volume.ShareType}
		for _, m := range c.volumeMetrics {
			ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.getterFn(volume), volumeLabels...)
		}
	}
	c.volumeTotalGauge.Collect(ch)
	c.scrapeCounter.Collect(ch)
	c.scrapeFailureCounter.Collect(ch)
	c.scrapeDurationGauge.Collect(ch)
	return
}

func (c *VolumeCollector) PeriodicFetch(cancelCh <-chan int) {
	var clearTimer *time.Timer
	startTimer := time.NewTimer(time.Millisecond)
	fetchTicker := time.NewTicker(c.fetchPeriod)

	for {
		select {
		case <-cancelCh:
			fetchTicker.Stop()
			clearTimer.Stop()
			break
		case <-fetchTicker.C:
		case <-startTimer.C:
			// Fetch immediately without waiting for the first tick
		}

		volumes := c.Fetch()
		if len(volumes) > 0 {
			c.mux.Lock()
			c.volumes = volumes
			c.mux.Unlock()
			c.volumeTotalGauge.Set(float64(len(volumes)))

			// Clear cached data if next fetching fails. This prevents exporting outdated data.
			if clearTimer != nil {
				clearTimer.Stop()
			}
			clearTimer = time.AfterFunc(2*c.fetchPeriod, func() {
				c.mux.Lock()
				c.volumes = nil
				c.mux.Unlock()
				c.volumeTotalGauge.Set(0)
				log.Debugf("VolumeCollector[%v] cleared cached volumes", c.filerName)
			})
		}
	}
}

func (c *VolumeCollector) Fetch() []*netapp.Volume {
	log.Debugf("VolumeCollector[%v] fetch() starts fetching volumes", c.filerName)
	start := time.Now()
	volumes, err := c.client.ListVolumes()
	elapsed := time.Now().Sub(start)
	c.scrapeCounter.Inc()
	c.scrapeDurationGauge.Set(elapsed.Seconds())
	if err != nil {
		log.Error(err)
		c.scrapeFailureCounter.Inc()
		return nil
	}
	log.Debugf("VolumeCollector[%v] fetch() fetched %d volumes", c.filerName, len(volumes))
	return volumes
}
