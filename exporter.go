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
			log.Println(err)
		}

		for _, v := range netappVolumes {

			_SizeAvailable, _ := strconv.ParseFloat(v.SizeAvailable, 64)
			_SizeTotal, _ := strconv.ParseFloat(v.SizeTotal, 64)
			_SizeUsed, _ := strconv.ParseFloat(v.SizeUsed, 64)
			_PercentageSizeUsed, _ := strconv.ParseFloat(v.PercentageSizeUsed, 64)
			_PercentageCompressionSpaceSaved, _ := strconv.ParseFloat(v.PercentageCompressionSpaceSaved, 64)
			_PercentageDeduplicationSpaceSaved, _ := strconv.ParseFloat(v.PercentageDeduplicationSpaceSaved, 64)
			_PercentageTotalSpaceSaved, _ := strconv.ParseFloat(v.PercentageTotalSpaceSaved, 64)

			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "total").Set(_SizeTotal)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "available").Set(_SizeAvailable)
			netappCapacity.WithLabelValues(v.ProjectID, v.ShareID, f.Name, v.Vserver, v.Volume, "used").Set(_SizeUsed)
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
			log.Println(err)
		}
		p.share = s
		time.Sleep(t * time.Second)
	}
}
