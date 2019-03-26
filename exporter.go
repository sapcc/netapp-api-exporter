package main

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	volumeCapacity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netapp",
			Subsystem: "capacity",
			Name:      "svm",
			Help:      "netapp SVM capacity",
		},
		[]string{
			"project_id",
			"share_id",
			"filer",
			"vserver",
			"volume",
			"metric",
		},
	)

	aggregateCapacity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netapp",
			Subsystem: "capacity",
			Name:      "aggregate",
			Help:      "Netapp aggregate capacity",
		},
		[]string{
			"availability_zone",
			"filer",
			"node",
			"aggregate",
			"metric",
		},
	)
)

type CapacityExporter struct {
	volumeCollector    prometheus.Collector
	aggregateCollector prometheus.Collector
	share              map[string]ManilaShare
}

func NewCapacityExporter() *CapacityExporter {
	return &CapacityExporter{
		volumeCollector:    volumeCapacity,
		aggregateCollector: aggregateCapacity,
	}
}

func (p *CapacityExporter) runGetNetappShare(f *Filer, t time.Duration) {
	for {
		netappVolumes, err := f.GetNetappVolume()
		if err != nil {
			logger.Println(err)
		}

		for _, v := range netappVolumes {
			_SizeTotal, _ := strconv.ParseFloat(v.SizeTotal, 64)
			_SizeAvailable, _ := strconv.ParseFloat(v.SizeAvailable, 64)
			_SizeUsed, _ := strconv.ParseFloat(v.SizeUsed, 64)
			_SizeUsedBySnapshots, _ := strconv.ParseFloat(v.SizeUsedBySnapshots, 64)
			_SizeAvailableForSnapshots, _ := strconv.ParseFloat(v.SizeAvailableForSnapshots, 64)
			_SnapshotReserveSize, _ := strconv.ParseFloat(v.SnapshotReserveSize, 64)
			_PercentageSizeUsed, _ := strconv.ParseFloat(v.PercentageSizeUsed, 64)
			_PercentageCompressionSpaceSaved, _ := strconv.ParseFloat(v.PercentageCompressionSpaceSaved, 64)
			_PercentageDeduplicationSpaceSaved, _ := strconv.ParseFloat(v.PercentageDeduplicationSpaceSaved, 64)
			_PercentageTotalSpaceSaved, _ := strconv.ParseFloat(v.PercentageTotalSpaceSaved, 64)

			// netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size").Set(float64(v.Size))
			volumeCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_total").Set(_SizeTotal)
			volumeCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_used").Set(_SizeUsed)
			volumeCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_available").Set(_SizeAvailable)
			volumeCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_used_by_snapshots").Set(_SizeUsedBySnapshots)
			volumeCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_available_for_snapshots").Set(_SizeAvailableForSnapshots)
			volumeCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_reserved_by_snapshots").Set(_SnapshotReserveSize)
			volumeCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "percentage_used").Set(_PercentageSizeUsed)
			volumeCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "percentage_compression_saved").Set(_PercentageCompressionSpaceSaved)
			volumeCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "percentage_deduplication_saved").Set(_PercentageDeduplicationSpaceSaved)
			volumeCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "percentage_total_saved").Set(_PercentageTotalSpaceSaved)
		}

		time.Sleep(t * time.Second)
	}
}

func (p *CapacityExporter) runGetNetappAggregate(f *Filer, t time.Duration) {
	for {
		aggrList := f.GetAggrData()

		for _, v := range aggrList {
			_percentrageUsed, _ := strconv.ParseFloat(v.PercentUsedCapacity, 64)
			aggregateCapacity.WithLabelValues(f.AvailabilityZone, f.Name, v.OwnerName, v.Name, "size_used").Set(float64(v.SizeUsed))
			aggregateCapacity.WithLabelValues(f.AvailabilityZone, f.Name, v.OwnerName, v.Name, "size_available").Set(float64(v.SizeAvailable))
			aggregateCapacity.WithLabelValues(f.AvailabilityZone, f.Name, v.OwnerName, v.Name, "size_total").Set(float64(v.SizeTotal))
			aggregateCapacity.WithLabelValues(f.AvailabilityZone, f.Name, v.OwnerName, v.Name, "size_total").Set(float64(v.SizeTotal))
			aggregateCapacity.WithLabelValues(f.AvailabilityZone, f.Name, v.OwnerName, v.Name, "percentage_used").Set(_percentrageUsed)
			aggregateCapacity.WithLabelValues(f.AvailabilityZone, f.Name, v.OwnerName, v.Name, "physical_used").Set(float64(v.PhysicalUsed))
			aggregateCapacity.WithLabelValues(f.AvailabilityZone, f.Name, v.OwnerName, v.Name, "physical_used_percent").Set(float64(v.PhysicalUsedPercent))
		}

		time.Sleep(t * time.Second)
	}
}

func (p *CapacityExporter) runGetOSShare(f *Filer, t time.Duration) {
	for {
		s, err := f.GetManilaShare()
		if err != nil {
			logger.Println(err)
		}
		p.share = s
		time.Sleep(t * time.Second)
	}
}
