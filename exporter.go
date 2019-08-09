package main

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type VolumeGaugeVec struct{ *prometheus.GaugeVec }
type AggrGaugeVec struct{ *prometheus.GaugeVec }

func NewVolumeGaugeVec() VolumeGaugeVec {
	return VolumeGaugeVec{
		prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "netapp",
				Subsystem: "capacity",
				Name:      "svm",
				Help:      "Netapp Volume Capacity",
			},
			[]string{
				"project_id",
				"share_id",
				"filer",
				"vserver",
				"volume",
				"metric",
			},
		),
	}
}

func NewAggrGaugeVec() AggrGaugeVec {
	return AggrGaugeVec{
		prometheus.NewGaugeVec(
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
		),
	}
}

type CapacityExporter struct {
	volumeCollector    prometheus.Collector
	aggregateCollector prometheus.Collector
}

func (vg *VolumeGaugeVec) SetMetric(v *NetappVolume) {
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
	vg.WithLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_total").Set(_SizeTotal)
	vg.WithLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_used").Set(_SizeUsed)
	vg.WithLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_available").Set(_SizeAvailable)
	vg.WithLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_used_by_snapshots").Set(_SizeUsedBySnapshots)
	vg.WithLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_available_for_snapshots").Set(_SizeAvailableForSnapshots)
	vg.WithLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_reserved_by_snapshots").Set(_SnapshotReserveSize)
	vg.WithLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "percentage_used").Set(_PercentageSizeUsed)
	vg.WithLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "percentage_compression_saved").Set(_PercentageCompressionSpaceSaved)
	vg.WithLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "percentage_deduplication_saved").Set(_PercentageDeduplicationSpaceSaved)
	vg.WithLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "percentage_total_saved").Set(_PercentageTotalSpaceSaved)
}

func (vg *VolumeGaugeVec) DeleteMetric(v *NetappVolume) {
	vg.DeleteLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_total")
	vg.DeleteLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_used")
	vg.DeleteLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_available")
	vg.DeleteLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_used_by_snapshots")
	vg.DeleteLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_available_for_snapshots")
	vg.DeleteLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "size_reserved_by_snapshots")
	vg.DeleteLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "percentage_used")
	vg.DeleteLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "percentage_compression_saved")
	vg.DeleteLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "percentage_deduplication_saved")
	vg.DeleteLabelValues(v.ProjectID, v.ShareID, v.FilerName, v.Vserver, v.Volume, "percentage_total_saved")
}

func (ag *AggrGaugeVec) SetMetric(v *Aggregate) {
	// aggrList := f.GetAggrData()
	_percentrageUsed, _ := strconv.ParseFloat(v.PercentUsedCapacity, 64)
	ag.WithLabelValues(v.AvailabilityZone, v.FilerName, v.OwnerName, v.Name, "size_used").Set(float64(v.SizeUsed))
	ag.WithLabelValues(v.AvailabilityZone, v.FilerName, v.OwnerName, v.Name, "size_available").Set(float64(v.SizeAvailable))
	ag.WithLabelValues(v.AvailabilityZone, v.FilerName, v.OwnerName, v.Name, "size_total").Set(float64(v.SizeTotal))
	ag.WithLabelValues(v.AvailabilityZone, v.FilerName, v.OwnerName, v.Name, "size_total").Set(float64(v.SizeTotal))
	ag.WithLabelValues(v.AvailabilityZone, v.FilerName, v.OwnerName, v.Name, "percentage_used").Set(_percentrageUsed)
	ag.WithLabelValues(v.AvailabilityZone, v.FilerName, v.OwnerName, v.Name, "physical_used").Set(float64(v.PhysicalUsed))
	ag.WithLabelValues(v.AvailabilityZone, v.FilerName, v.OwnerName, v.Name, "physical_used_percent").Set(float64(v.PhysicalUsedPercent))
}
