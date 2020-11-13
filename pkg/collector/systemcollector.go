package collector

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/netapp-api-exporter/pkg/netapp"
	log "github.com/sirupsen/logrus"
)

type SystemCollector struct {
	filerName   string
	versionDesc *prometheus.Desc
	client      *netapp.Client
}

func NewSystemCollector(filerName string, client *netapp.Client) *SystemCollector {
	return &SystemCollector{
		filerName: filerName,
		client:    client,
		versionDesc: prometheus.NewDesc(
			"netapp_filer_system_version",
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
	if err != nil {
		log.Error(err)
		return
	}
	version := fullVersion[:strings.Index(fullVersion, ":")]
	ch <- prometheus.MustNewConstMetric(
		c.versionDesc,
		prometheus.GaugeValue,
		0.0,
		fullVersion,
		version,
	)
}
