package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/netapp-api-exporter/pkg/netapp"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

const netappApiVersion = "1.7"

type FilerBase struct {
	Name             string `yaml:"name"`
	Host             string `yaml:"host"`
	AvailabilityZone string `yaml:"availability_zone"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	Version          string `yaml:"version"`
}

type Filer struct {
	FilerBase
	Client         *netapp.Client
	ScrapeFailures *prometheus.CounterVec
}

func NewFiler(f FilerBase) Filer {
	filer := Filer{
		FilerBase: f,
		Client:    netapp.NewClient(f.Host, f.Username, f.Password, f.Version),
		ScrapeFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "netapp_filer_scrape_failure",
				Help: "Number of failed scrapes to netapp filer.",
			},
			[]string{"status"},
		),
	}

	// check if client works properly every 5 miniutes
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		timer := time.NewTimer(time.Millisecond)
		for {
			select {
			case <-ticker.C:
			case <-timer.C:
			}
			statusCode, err := filer.Client.CheckCluster()
			if err != nil {
				log.Errorf("check client: %v", err)
			}
			switch statusCode {
			case 200, 201, 202, 204, 205, 206:
			default:
				filer.ScrapeFailures.With(
					prometheus.Labels{"status": strconv.Itoa(statusCode)},
				).Inc()
			}
		}
	}()

	return filer
}

func loadFilers(configFile string) ([]Filer, error) {
	if os.Getenv("DEV") == "1" {
		log.SetLevel(log.DebugLevel)
		log.Debug("Set log level to DebugLevel")
	}
	if len(configFile) == 0 {
		log.Info("Load filer configuration from env variables")
		return []Filer{loadFilerFromEnv()}, nil
	} else {
		log.Infof("Load filer configuration from %s", configFile)
		return loadFilerFromFile(configFile)
	}
}

func loadFilerFromFile(fileName string) (filers []Filer, err error) {
	var yamlFile []byte
	var filerInfos []*FilerBase
	if yamlFile, err = ioutil.ReadFile(fileName); err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, &filerInfos); err != nil {
		return nil, err
	}
	for _, f := range filerInfos {
		if f.Username == "" || f.Password == "" {
			username, password := getAuthFromEnv()
			f.Username = username
			f.Password = password
		}
		if f.Version == "" {
			f.Version = netappApiVersion
		}
		filers = append(filers, NewFiler(*f))
	}
	return
}

func loadFilerFromEnv() Filer {
	name := os.Getenv("NETAPP_NAME")
	host := os.Getenv("NETAPP_HOST")
	username := os.Getenv("NETAPP_USERNAME")
	password := os.Getenv("NETAPP_PASSWORD")
	az := os.Getenv("NETAPP_AZ")
	version := getEnvWithDefaultValue("Netapp_API_VERSION", netappApiVersion)
	return NewFiler(FilerBase{name, host, az, username, password, version})
}

func getAuthFromEnv() (username, password string) {
	username = os.Getenv("NETAPP_USERNAME")
	password = os.Getenv("NETAPP_PASSWORD")
	return
}

func getEnvWithDefaultValue(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	} else {
		return defaultValue
	}
}
