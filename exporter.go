package main

import (
	"fmt"
	"time"

	"github.com/pepabo/go-netapp/netapp"

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

	vserverOptions = netapp.VServerOptions{
		Query: &netapp.VServerQuery{
			VServerInfo: &netapp.VServerInfo{
				VserverType: "cluster | data",
			},
		},
		DesiredAttributes: &netapp.VServerQuery{
			VServerInfo: &netapp.VServerInfo{
				VserverName: "x",
				UUID:        "x",
			},
		},
		MaxRecords: 100,
	}

	volumeOptions = netapp.VolumeOptions{
		MaxRecords: 200,
		Query: &netapp.VolumeQuery{
			VolumeInfo: &netapp.VolumeInfo{
				VolumeIDAttributes: &netapp.VolumeIDAttributes{
					OwningVserverUUID: "x",
				},
			},
		},
		DesiredAttributes: &netapp.VolumeQuery{
			VolumeInfo: &netapp.VolumeInfo{
				VolumeIDAttributes: &netapp.VolumeIDAttributes{
					Name:              "x",
					OwningVserverName: "x",
					OwningVserverUUID: "x",
				},
				VolumeSpaceAttributes: &netapp.VolumeSpaceAttributes{
					//
					Size:                1,
					SizeTotal:           "x",
					SizeAvailable:       "x",
					SizeUsed:            "x",
					SizeUsedBySnapshots: "x",
					PercentageSizeUsed:  "x",
				},
			},
		},
	}

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

func (p *CapacityExporter) run(f *filer, t time.Duration) {

	for {
		// capa := []CapacityData{}
		// fmt.Println(capa)

		vserverList, _, _ := f.Client.VServer.List(&vserverOptions)
		fmt.Println("vserverList ", vserverList)

		for i, vserver := range vserverList.Results.AttributesList.VserverInfo {
			if i > 1 {
				break
			}
			fmt.Println(volumeOptions.Query)
			volumeOptions.Query.VolumeInfo.VolumeIDAttributes.OwningVserverUUID = vserver.UUID
			vol, _, _ := f.Client.Volume.List(&volumeOptions)
			fmt.Println(vol)
		}

		// for _, d := range capa {
		// 	netappCapacity.WithLabelValues(f.Name, d.Vserver, d.Volume, "total").Set(d.Space.TotalSize)
		// 	netappCapacity.WithLabelValues(f.Name, d.Vserver, d.Volume, "available").Set(d.Space.AvailabeSize)
		// 	netappCapacity.WithLabelValues(f.Name, d.Vserver, d.Volume, "used").Set(d.Space.UsedSize)
		// 	netappCapacity.WithLabelValues(f.Name, d.Vserver, d.Volume, "percentage_used").Set(d.Space.UsedPercentage)
		// }

		time.Sleep(t * time.Second)
	}
}
