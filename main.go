package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/netapp-api-exporter/pkg/collector"
	"gopkg.in/alecthomas/kingpin.v2"

	log "github.com/sirupsen/logrus"
)

var (
	configFile        = kingpin.Flag("config", "Config file").Short('c').Default("").String()
	listenAddress     = kingpin.Flag("listen", "Listen address").Short('l').Default("0.0.0.0").String()
	debug             = kingpin.Flag("debug", "Debug mode").Short('d').Bool()
	volumeFetchPeriod = kingpin.Flag("volume-fetch-period", "Period of asynchronously fetching volumes").Short('v').Default("2m").Duration()
	disableAggregate  = kingpin.Flag("no-aggregate", "Disable aggregate collector").Bool()
	disableVolume     = kingpin.Flag("no-volume", "Disable volume collector").Bool()
	disableSystem     = kingpin.Flag("no-system", "Disable system collector").Bool()
)

type logFormatter struct{}

func init() {
	kingpin.Parse()

	log.SetOutput(os.Stdout)
	log.SetFormatter(new(logFormatter))
	if *debug {
		log.Info("Debug mode")
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	// try loading filers every  10 seconds until successful
	var filers []Filer
	var err error
	for {
		filers, err = loadFilers(*configFile)
		if err != nil {
			log.Errorf("Failed to load filer configuration: %v. Retry in 10 seconds...", err)
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}

	reg := prometheus.NewPedanticRegistry()

	// register go process collector
	reg.MustRegister(prometheus.NewGoCollector())

	for _, f := range filers {
		extraLabels := prometheus.Labels{
			"filer":             f.Name,
			"availability_zone": f.AvailabilityZone,
		}
		log.Infof("Register collectors for filer: {Name=%s, Host=%s, Username=%s}", f.Name, f.Host, f.Username)
		prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(f.ScrapeFailures)
		if !*disableAggregate {
			prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(
				collector.NewAggregateCollector(f.Name, f.Client))
		}
		if !*disableVolume {
			prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(
				collector.NewVolumeCollector(f.Name, f.Client, *volumeFetchPeriod))
		}
		if !*disableSystem {
			prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(
				collector.NewSystemCollector(f.Name, f.Client))
		}
	}

	port := "9108"
	addr := *listenAddress + ":" + port
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Debugf("Open link http://%s/metrics for metrics", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (f *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	s := fmt.Sprintf("%s [%s] %s\t", entry.Time.Format("2006-01-02 15:04:05.000"), entry.Level, entry.Message)
	for k, v := range entry.Data {
		s = s + fmt.Sprintf(" %s=%s", k, v)
	}
	s = s + "\n"
	return []byte(s), nil
}
