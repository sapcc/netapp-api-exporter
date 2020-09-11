package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/netapp-api-exporter/netapp"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFile               = kingpin.Flag("config", "Config file").Short('c').Default("./config/netapp_filers.yaml").String()
	listenAddress            = kingpin.Flag("listen", "Listen address").Short('l').Default("0.0.0.0").String()
	debug                    = kingpin.Flag("debug", "Debug mode").Short('d').Bool()
	aggregateRetentionPeriod = kingpin.Flag("aggregateRetention", "Aggregate collector retention period").Default("5m").Duration()
	volumeRetentionPeriod    = kingpin.Flag("volumeRetention", "Volume collector retention period").Default("2m").Duration()
	disableAggregate         = kingpin.Flag("no-aggregate", "Disable aggregate collector").Bool()
	disableVolume            = kingpin.Flag("no-volume", "Disable volume collector").Bool()
	disableSystem            = kingpin.Flag("no-system", "Disable system collector").Bool()
	logger                   = logrus.New()
)

type logFormatter struct{}

func init() {
	kingpin.Parse()

	if os.Getenv("DEV") != "" {
		*debug = true
	}

	logger.Out = os.Stdout
	logger.SetFormatter(new(logFormatter))
	if *debug {
		logger.Info("Debug mode")
		logger.Level = logrus.DebugLevel
	} else {
		logger.Level = logrus.InfoLevel
	}
}

func main() {
	// try loading filers every  10 seconds until successful
	var filers []*Filer
	var err error
	for {
		filers, err = loadFilers()
		if err != nil {
			logger.Errorf("Failed to load filer configuration: %v. Retry in 10 seconds...", err)
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}

	reg := prometheus.NewPedanticRegistry()

	for _, f := range filers {
		netappClient := netapp.NewClient(f.Host, f.Username, f.Password, f.Version)
		extraLabels := prometheus.Labels{
			"filer":             f.Name,
			"availability_zone": f.AvailabilityZone,
		}
		logger.Infof("Register collectors for filer: {Name=%s, Host=%s, Username=%s}", f.Name, f.Host, f.Username)
		if !*disableAggregate {
			prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(
				NewAggregateCollector(f.Name, netappClient, *aggregateRetentionPeriod))
		}
		if !*disableVolume {
			prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(
				NewVolumeCollector(f.Name, netappClient, *volumeRetentionPeriod))
		}
		if !*disableSystem {
			prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(NewSystemCollector(f.Name, netappClient))
		}
	}

	port := "9108"
	addr := *listenAddress + ":" + port
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	logger.Fatal(http.ListenAndServe(addr, nil))
}

func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	s := fmt.Sprintf("%s [%s] %s\t", entry.Time.Format("2006-01-02 15:04:05.000"), entry.Level, entry.Message)
	for k, v := range entry.Data {
		s = s + fmt.Sprintf(" %s=%s", k, v)
	}
	s = s + "\n"
	return []byte(s), nil
}
