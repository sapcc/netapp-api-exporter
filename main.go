package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
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

	DNSErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "netapp_filer_dns_error",
			Help: "hostname not resolved",
		},
		[]string{"host"},
	)
	AuthenticationErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "netapp_filer_authentication_error",
			Help: "access netapp filer failed with 401",
		},
		[]string{"host"},
	)
	TimeoutErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "netapp_filer_timeout_error",
			Help: "access netapp filer timeout",
		},
		[]string{"host"},
	)
	UnknownErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "netapp_filer_unknown_error",
			Help: "check filer failed with unknown error",
		},
		[]string{"host"},
	)
)

func main() {
	var filers map[string]Filer

	// new prometheus registry and register global collectors
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(prometheus.NewGoCollector())
	reg.MustRegister(DNSErrorCounter)
	reg.MustRegister(AuthenticationErrorCounter)
	reg.MustRegister(TimeoutErrorCounter)
	reg.MustRegister(UnknownErrorCounter)

	// load filers from configuration and register new colloector for new filer
	go func() {
		// fast ticker for iniital load; will be stopped after 10 times or
		// filers are loaded
		fastTickerCounter := 0
		fastTicker := time.NewTicker(10 * time.Second)
		ticker := time.NewTicker(5 * time.Minute)

		for {
			ff, err := loadFilers(*configFile)
			if err != nil {
				log.WithError(err).Error("load filers failed")
			} else {
				for _, f := range ff {
					if _, ok := filers[f.Host]; ok {
						continue
					}
					l := log.WithFields(log.Fields{
						"Name":             f.Name,
						"Host":             f.Host,
						"AvailabilityZone": f.AvailabilityZone,
						"AggregatePattern": f.AggregatePattern,
					})
					l.Info("check filer")
					if !checkFiler(f, l) {
						continue
					}
					l.Info("register filer")
					err = registerFiler(reg, f)
					if err != nil {
						l.Error(err)
						continue
					}
					if filers == nil {
						filers = make(map[string]Filer)
					}
					filers[f.Host] = f
				}
			}

			select {
			case <-fastTicker.C:
				fastTickerCounter += 1
				if fastTickerCounter == 10 || filers != nil {
					fastTicker.Stop()
				}
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

func checkFiler(f Filer, l *log.Entry) bool {
	var dnsError *net.DNSError
	status, err := f.Client.CheckCluster()
	l = l.WithField("status", strconv.Itoa(status))
	switch status {
	case 200, 201, 202, 204, 205, 206:
	case 401:
		AuthenticationErrorCounter.WithLabelValues(f.Host).Inc()
		l.Error("check filer failed: authentication error")
		return false
	default:
		if err != nil {
			l.WithError(err).Error("check filer failed")
			if errors.As(err, &dnsError) {
				DNSErrorCounter.WithLabelValues(f.Host).Inc()
			} else if errors.Is(err, context.DeadlineExceeded) {
				TimeoutErrorCounter.WithLabelValues(f.Host).Inc()
			} else {
				UnknownErrorCounter.WithLabelValues(f.Host).Inc()
			}
		} else {
			UnknownErrorCounter.WithLabelValues(f.Host).Inc()
			l.Error("check filer failed")
		}
		return false
	}
	return true
}

func registerFiler(reg prometheus.Registerer, f Filer) error {
	if f.Name == "" {
		return fmt.Errorf("Filer.Name not set")
	}
	if f.AvailabilityZone == "" {
		return fmt.Errorf("Filer.AvailabilityZone not set")
	}
	extraLabels := prometheus.Labels{
		"filer":             f.Name,
		"availability_zone": f.AvailabilityZone,
	}
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
	log.SetFormatter(&log.TextFormatter{})
	if *debug {
		log.Info("debug mode")
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
