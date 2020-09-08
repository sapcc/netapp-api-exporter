package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/netapp-api-exporter/netapp"
	"strings"
)

type SystemCollector struct {
	versionDesc *prometheus.Desc
	client      *netapp.Client
}

func NewSystemCollector(client *netapp.Client) *SystemCollector {
	return &SystemCollector{
		client: client,
		versionDesc: prometheus.NewDesc(
			"netapp_system_version",
			"Info about ontap version in labels `version` and `full_version`",
			[]string{"full_version", "version"},
			nil,
		),
	}
}

func (c *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.versionDesc
}

func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	fullVersion, err := c.client.GetSystemVersion()
	version := fullVersion[:strings.Index(fullVersion, ":")]
	if err != nil {
		logger.Error(err)
	}
	ch <- prometheus.MustNewConstMetric(
		c.versionDesc,
		prometheus.GaugeValue,
		0.0,
		fullVersion,
		version,
	)
}
