package main

import (
	"log"
	"strings"
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
			"filer",
			"svm",
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
			log.Panic(err)
		}

		// projectId := p.share["maurice_test"].ProjectId
		var projectID string

		for _, d := range netappVolumes {
			if strings.HasPrefix(d.Vserver, "ma_") {
				projectID = "05f9781218b7401d9955f9b8a05a5aea"
			} else {
				projectID = ""
			}

			netappCapacity.WithLabelValues(projectID, f.Name, d.Vserver, d.Volume, "total").Set(d.SizeTotal)
			netappCapacity.WithLabelValues(projectID, f.Name, d.Vserver, d.Volume, "available").Set(d.SizeAvailable)
			netappCapacity.WithLabelValues(projectID, f.Name, d.Vserver, d.Volume, "used").Set(d.SizeUsed)
			netappCapacity.WithLabelValues(projectID, f.Name, d.Vserver, d.Volume, "percentage_used").Set(d.PercentageSizeUsed)
		}

		time.Sleep(t * time.Second)
	}
}

func (p *CapacityExporter) runGetOSShare(f *Filer, t time.Duration) {
	for {
		p.share = f.GetManilaShare()
		time.Sleep(t * time.Second)
	}
}
