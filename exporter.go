package main

import (
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
			"filer",
			"svm",
			"volume",
			"metric",
		},
	)

// 	Query: &netapp.VolumeQuery{
// 			VolumeInfo: &netapp.VolumeInfo{
// 				VolumeIDAttributes: &netapp.VolumeIDAttributes{
// 					OwningVserverUUID: "x",
// 				},
// 			},
// 		},
// 	}
// }

)

type CapacityExporter struct {
	collector prometheus.Collector
}

type CapacityData struct {
	Vserver string
	Volume  string
	Project string
	Space   struct {
		AvailabeSize   float64
		TotalSize      float64
		UsedSize       float64
		UsedPercentage float64
	}
}

func NewCapacityExporter() *CapacityExporter {
	return &CapacityExporter{netappCapacity}
}

func (p *CapacityExporter) runGetNetappShare(f *Filer, t time.Duration) {

	for {
		f.GetNetappShare()

		// for _, d := range capa {
		// 	netappCapacity.WithLabelValues(f.Name, d.Vserver, d.Volume, "total").Set(d.Space.TotalSize)
		// 	netappCapacity.WithLabelValues(f.Name, d.Vserver, d.Volume, "available").Set(d.Space.AvailabeSize)
		// 	netappCapacity.WithLabelValues(f.Name, d.Vserver, d.Volume, "used").Set(d.Space.UsedSize)
		// 	netappCapacity.WithLabelValues(f.Name, d.Vserver, d.Volume, "percentage_used").Set(d.Space.UsedPercentage)
		// }

		time.Sleep(t * time.Second)
	}
}

func (p *CapacityExporter) runGetOSShare(f *Filer, t time.Duration) {
	for {
		f.GetOSShare()
		time.Sleep(t * time.Second)
	}
}
