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
	var projectID string

	for {
		netappVolumes, err := f.GetNetappVolume()
		if err != nil {
			log.Println(err)
		}

		for _, d := range netappVolumes {
			if strings.HasPrefix(d.Vserver, "ma_") && strings.HasPrefix(d.Volume, "share_") {
				siid := strings.TrimPrefix(d.Volume, "share_")
				if share, ok := p.share[siid]; ok {
					projectID = share.ProjectId
					// log.Printf("%+v", share)
					// log.Printf("%+v\n", d)
				}
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
		s, err := f.GetManilaShare()
		if err != nil {
			log.Println(err)
		}
		p.share = s
		time.Sleep(t * time.Second)
	}
}
