package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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

func main() {
	var filers map[string]Filer

	// new prometheus exporter registry and register go process collector
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(prometheus.NewGoCollector())

	// load filers from configuration every minute and register new colloector
	// for new filer
	go func() {
		ticker := time.NewTicker(1 * time.Minute)

		for {
			ff, err := loadFilers(*configFile)
			if err != nil {
				log.WithField("File", *configFile).Errorf("failed to load configuration: %v", err)
				return
			} else {
				// register collector for new filer
				for _, f := range ff {
					if _, ok := filers[f.Host]; !ok {
						err = registerFiler(reg, f)
						if err != nil {
							log.Error(err)
						}
						if filers == nil {
							filers = make(map[string]Filer)
						}
						filers[f.Host] = f
					}
				}
			}

			select {
			case <-ticker.C:
			}
		}
	}()

	port := "9108"
	addr := *listenAddress + ":" + port
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Debugf("open link http://%s/metrics for metrics", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func registerFiler(reg prometheus.Registerer, f Filer) error {
	if f.Name == "" {
		return fmt.Errorf("Filer.Name is not set ")
	}
	if f.AvailabilityZone == "" {
		return fmt.Errorf("Filer.AvailabilityZone is not set ")
	}
	log.WithFields(log.Fields{
		"Name":      f.Name,
		"Host":      f.Host,
		"AZ":        f.AvailabilityZone,
		"Aggregate": f.AggregatePattern,
	}).Info("register filer")
	extraLabels := prometheus.Labels{
		"filer":             f.Name,
		"availability_zone": f.AvailabilityZone,
	}
	prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(f.ScrapeFailures)
	prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(f.FilerDNSFailures)
	prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(f.FilerTimeoutFailures)
	if !*disableAggregate {
		prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(
			collector.NewAggregateCollector(f.Client, f.Name, f.AggregatePattern))
	}
	if !*disableVolume {
		prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(
			collector.NewVolumeCollector(f.Client, f.Name, *volumeFetchPeriod))
	}
	if !*disableSystem {
		prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(
			collector.NewSystemCollector(f.Client, f.Name))
	}
	return nil
}

func init() {
	kingpin.Parse()

	log.SetOutput(os.Stdout)
	log.SetFormatter(new(logFormatter))
	if *debug {
		log.Info("debug mode")
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

type logFormatter struct{}

func (f *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	s := fmt.Sprintf("%s %-5v msg=%q",
		entry.Time.Format("2006-01-02 15:04:05.000"),
		strings.ToUpper(entry.Level.String()),
		entry.Message)
	for k, v := range entry.Data {
		s = s + fmt.Sprintf(" %s=%s", k, v)
	}
	s = s + "\n"
	return []byte(s), nil
}
