package collector

import (
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/netapp-api-exporter/pkg/netapp"
	log "github.com/sirupsen/logrus"
)

type AggregateCollector struct {
	client               *netapp.Client
	filerName            string
	aggregatePattern     string
	aggregateMetrics     []AggregateMetric
	scrapeCounter        prometheus.Counter
	scrapeFailureCounter prometheus.Counter
	scrapeDurationGauge  prometheus.Gauge
}

type AggregateMetric struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	getterFn  func(aggr *netapp.Aggregate) float64
}

func NewAggregateCollector(client *netapp.Client, filerName, aggrPattern string) *AggregateCollector {
	aggrLabels := []string{"node", "aggregate"}
	aggrMetrics := []AggregateMetric{
		{
			desc: prometheus.NewDesc(
				"netapp_aggregate_total_bytes",
				"Netapp Aggregate Metrics: total size",
				aggrLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(m *netapp.Aggregate) float64 { return m.SizeTotal },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_available_bytes",
				"Netapp Aggregate Metrics: available size",
				aggrLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(m *netapp.Aggregate) float64 { return m.SizeAvailable },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_used_bytes",
				"Netapp Aggregate Metrics: used size",
				aggrLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(m *netapp.Aggregate) float64 { return m.SizeUsed },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_used_percentage",
				"Netapp Aggregate Metrics: used percentage",
				aggrLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(m *netapp.Aggregate) float64 { return m.PercentUsedCapacity },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_physical_used_bytes",
				"Netapp Aggregate Metrics: physical used size",
				aggrLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(m *netapp.Aggregate) float64 { return m.PhysicalUsed },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_physical_used_percentage",
				"Netapp Aggregate Metrics: physical used percentage",
				aggrLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn:  func(m *netapp.Aggregate) float64 { return m.PhysicalUsedPercent },
		}, {
			desc: prometheus.NewDesc(
				"netapp_aggregate_is_encrypted",
				"Netapp Aggregate Metrics: is encrypted",
				aggrLabels,
				nil),
			valueType: prometheus.GaugeValue,
			getterFn: func(m *netapp.Aggregate) float64 {
				if m.IsEncrypted {
					return 1.0
				}
				return 0.0
			},
		},
	}
	scrapeDurationGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "netapp_aggregate_scrape_duration_seconds",
			Help: "duration in seconds of fetching aggregates from filer",
		},
	)
	scrapeCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "netapp_aggregate_scrape_total",
			Help: "number of aggregate fetches from filer",
		},
	)
	scrapeFailureCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "netapp_aggregate_scrape_failure_total",
			Help: "number of failures for fetching aggregates from filer",
		},
	)
	return &AggregateCollector{
		client:               client,
		filerName:            filerName,
		aggregatePattern:     aggrPattern,
		aggregateMetrics:     aggrMetrics,
		scrapeDurationGauge:  scrapeDurationGauge,
		scrapeCounter:        scrapeCounter,
		scrapeFailureCounter: scrapeFailureCounter,
	}
}

func (c *AggregateCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.aggregateMetrics {
		ch <- m.desc
	}
	ch <- c.scrapeCounter.Desc()
	ch <- c.scrapeFailureCounter.Desc()
	ch <- c.scrapeDurationGauge.Desc()
}

func (c *AggregateCollector) Collect(ch chan<- prometheus.Metric) {
	// fetch aggregates
	aggregates := c.Fetch()

	// export metrics
	for _, aggr := range aggregates {
		// filter aggregate here
		matched, err := regexp.MatchString(c.aggregatePattern, aggr.Name)
		if err != nil {
			log.Error(err)
		}
		if !matched {
			log.Debugf("AggregateCollector[%v] Collect(): %s does not match "+
				"pattern %q", c.filerName, aggr.Name, c.aggregatePattern)
			continue
		}

		labels := []string{aggr.OwnerName, aggr.Name}
		for _, m := range c.aggregateMetrics {
			ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.getterFn(aggr), labels...)
		}
	}
	c.scrapeCounter.Collect(ch)
	c.scrapeFailureCounter.Collect(ch)
	c.scrapeDurationGauge.Collect(ch)
}

func (c *AggregateCollector) Fetch() []*netapp.Aggregate {
	start := time.Now()
	aggregates, err := c.client.ListAggregates()
	elapsed := time.Now().Sub(start)
	c.scrapeCounter.Inc()
	c.scrapeDurationGauge.Set(elapsed.Seconds())
	if err != nil {
		log.Error(err)
		c.scrapeFailureCounter.Inc()
		return nil
	}
	return aggregates
}
