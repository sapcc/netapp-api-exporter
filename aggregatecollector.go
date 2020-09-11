package main

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/netapp-api-exporter/netapp"
)

type AggregateCollector struct {
	filerName       string
	client          *netapp.Client
	metrics         []AggregateMetric
	aggregates      []*netapp.Aggregate
	mux             sync.Mutex
	retentionPeriod time.Duration
	errorCh         chan<- error
}

type AggregateMetric struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	getterFn  func(aggr *netapp.Aggregate) float64
}

func NewAggregateCollector(filerName string, client *netapp.Client, ch chan<- error, retentionPeriod time.Duration) *AggregateCollector {
	aggrLabels := []string{"node", "aggregate"}
	return &AggregateCollector{
		filerName:       filerName,
		client:          client,
		errorCh:         ch,
		retentionPeriod: retentionPeriod,
		metrics: []AggregateMetric{
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
			},
		},
	}
}

func (c *AggregateCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m.desc
	}
}

func (c *AggregateCollector) Collect(ch chan<- prometheus.Metric) {
	defer c.mux.Unlock()
	c.mux.Lock()
	if c.aggregates == nil {
		if err := c.Fetch(); err != nil {
			logger.Error(err)
			c.errorCh <- err
			return
		}
	}
	for _, aggr := range c.aggregates {
		labels := []string{aggr.OwnerName, aggr.Name}
		for _, m := range c.metrics {
			ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.getterFn(aggr), labels...)
		}
	}
}

func (c *AggregateCollector) Fetch() error {
	aggregates, err := c.client.ListAggregates()
	if err != nil {
		return err
	}
	c.aggregates = aggregates
	time.AfterFunc(c.retentionPeriod, func() {
		defer c.mux.Unlock()
		c.mux.Lock()
		c.aggregates = nil
	})
	return nil
}
