package main

import (
	"github.com/pepabo/go-netapp/netapp"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"sync"
	"time"
)

type NetappAggregate struct {
	AvailabilityZone    string
	FilerName           string
	Name                string
	OwnerName           string
	SizeUsed            float64
	SizeTotal           float64
	SizeAvailable       float64
	TotalReservedSpace  float64
	PercentUsedCapacity float64
	PhysicalUsed        float64
	PhysicalUsedPercent float64
}

type AggrManager struct {
	sync.Mutex
	Aggregates    []*NetappAggregate
	lastFetchTime time.Time
	maxAge        time.Duration
}

type aggregateMetrics []struct {
	desc    *prometheus.Desc
	valType prometheus.ValueType
	evalFn  func(agg *NetappAggregate) float64
}

var (
	aggregateLabels = []string{
		"node",
		"aggregate",
	}

	aggMetrics = aggregateMetrics{
		{
			desc: prometheus.NewDesc(
				"netapp_aggregate_total_bytes",
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

func (a AggrManager) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range aggMetrics {
		ch <- v.desc
	}
}

func (a AggrManager) Collect(ch chan<- prometheus.Metric) {
	for _, v := range a.Aggregates {
		labels := []string{v.OwnerName, v.Name}
		for _, m := range aggMetrics {
			ch <- prometheus.MustNewConstMetric(m.desc, m.valType, m.evalFn(v), labels...)
		}
	}
}

func (a AggrManager) Fetch(f Filer) (aggregates []*NetappAggregate, err error) {
	ff := new(bool)
	*ff = false
	opts := &netapp.AggrOptions{
		Query: &netapp.AggrInfo{
			AggrRaidAttributes: &netapp.AggrRaidAttributes{
				IsRootAggregate: ff,
			},
		},
		DesiredAttributes: &netapp.AggrInfo{
			AggrOwnershipAttributes: &netapp.AggrOwnershipAttributes{},
			AggrSpaceAttributes:     &netapp.AggrSpaceAttributes{},
		},
	}

	aggrs, err := a.fetch(f, opts)

	if err == nil {
		logger.Printf("%s: %d aggregates fetched", f.Host, len(aggrs))
		for _, n := range aggrs {
			percentUsedCapacity, _ := strconv.ParseFloat(n.AggrSpaceAttributes.PercentUsedCapacity, 64)
			aggregates = append(aggregates, &NetappAggregate{
				AvailabilityZone:    f.AvailabilityZone,
				FilerName:           f.Name,
				Name:                n.AggregateName,
				OwnerName:           n.AggrOwnershipAttributes.OwnerName,
				SizeUsed:            float64(n.AggrSpaceAttributes.SizeUsed),
				SizeTotal:           float64(n.AggrSpaceAttributes.SizeTotal),
				SizeAvailable:       float64(n.AggrSpaceAttributes.SizeAvailable),
				TotalReservedSpace:  float64(n.AggrSpaceAttributes.TotalReservedSpace),
				PercentUsedCapacity: percentUsedCapacity,
				PhysicalUsed:        float64(n.AggrSpaceAttributes.PhysicalUsed),
				PhysicalUsedPercent: float64(n.AggrSpaceAttributes.PhysicalUsedPercent),
			})
		}
	}
	return
}

func (a AggrManager) fetch(f Filer, opts *netapp.AggrOptions) (res []netapp.AggrInfo, err error) {
	pageHandler := func(r netapp.AggrListPagesResponse) bool {
		if r.Error != nil {
			err = r.Error
			return false
		}
		res = append(res, r.Response.Results.AggrAttributes...)
		return true
	}
	f.NetappClient.Aggregate.ListPages(opts, pageHandler)
	return
}
