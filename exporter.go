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
			log.Println(err)
		}

		for _, v := range netappVolumes {
			projectID := ""

			if strings.HasPrefix(v.Vserver, "ma_") && strings.HasPrefix(v.Volume, "share_") {
				siid := strings.TrimPrefix(v.Volume, "share_")
				if share, ok := p.share[siid]; ok {
					projectID = share.ProjectId
					// log.Printf("%+v", share)
					// log.Printf("%+v\n", d)
				}
			}

			netappCapacity.WithLabelValues(projectID, f.Name, v.Vserver, v.Volume, "total").Set(v.SizeTotal)
			netappCapacity.WithLabelValues(projectID, f.Name, v.Vserver, v.Volume, "available").Set(v.SizeAvailable)
			netappCapacity.WithLabelValues(projectID, f.Name, v.Vserver, v.Volume, "used").Set(v.SizeUsed)
			netappCapacity.WithLabelValues(projectID, f.Name, v.Vserver, v.Volume, "percentage_used").Set(v.PercentageSizeUsed)
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
