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
)

type CapacityExporter struct {
	collector prometheus.Collector
}

func NewCapacityExporter() *CapacityExporter {
	return &CapacityExporter{netappCapacity}
}

func (p *CapacityExporter) run(f *filer, t time.Duration) {
	for {
		// getData(&f)
		time.Sleep(t * time.Second)
	}
}

// vserverInfo := &netapp.VServerInfo{
// 	VserverName:   "1",
// 	UUID:          "1",
// 	State:         "1",
// 	AggregateList: &[]string{"x"},
// }

// volumeQuery:= &netapp.VolumeInfo{

// }

// volumeInfo := &netapp.VolumeInfo{
// 	VolumeIDAttributes: &netapp.VolumeIDAttributes{
// 		Name:              "x",
// 		OwningVserverName: "x",
// 	},
// }
