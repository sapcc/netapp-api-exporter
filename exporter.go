package main

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	netappCapacity = prometheus.NewGaugeVec(
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
)

type CapacityExporter struct {
	collector prometheus.Collector
	share     map[string]ManilaShare
}

func NewCapacityExporter() *CapacityExporter {
	return &CapacityExporter{collector: netappCapacity}
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
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_total").Set(_SizeTotal)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_used").Set(_SizeUsed)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_available").Set(_SizeAvailable)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_used_by_snapshots").Set(_SizeUsedBySnapshots)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_available_for_snapshots").Set(_SizeAvailableForSnapshots)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "size_reserved_by_snapshots").Set(_SnapshotReserveSize)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "percentage_used").Set(_PercentageSizeUsed)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "percentage_compression_saved").Set(_PercentageCompressionSpaceSaved)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "percentage_deduplication_saved").Set(_PercentageDeduplicationSpaceSaved)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "percentage_total_saved").Set(_PercentageTotalSpaceSaved)
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
