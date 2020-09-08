package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/netapp-api-exporter/netapp"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	configFile    = kingpin.Flag("config", "Config file").Short('c').Default("./config/netapp_filers.yaml").String()
	listenAddress = kingpin.Flag("listen", "Listen address").Short('l').Default("0.0.0.0").String()
	debug         = kingpin.Flag("debug", "Debug mode").Short('d').Bool()
	logger        = logrus.New()
)

type myFormatter struct{}

func init() {
	kingpin.Parse()

	if os.Getenv("DEV") != "" {
		*debug = true
	}

	logger.Out = os.Stdout
	logger.SetFormatter(new(myFormatter))
	if *debug {
		logger.Level = logrus.DebugLevel
	} else {
		logger.Level = logrus.InfoLevel
	}
}

func main() {
	// try loading filers every  10 seconds until successful
	var filers []*netapp.Filer
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
		netappClient := netapp.NewClient(f)
		extraLabels := prometheus.Labels{
			"filer":             f.Name,
			"availability_zone": f.AvailabilityZone,
		}
		logger.Infof("Register collectors for filer: {Name=%s, Host=%s, Username=%s}", f.Name, f.Host, f.Username)
		prometheus.WrapRegistererWith(extraLabels, reg).MustRegister(
			NewAggregateCollector(netappClient, 5*time.Minute),
			NewVolumeCollector(netappClient, 2*time.Minute),
			NewSystemCollector(netappClient),
		)
	}

	port := "9108"
	addr := *listenAddress + ":" + port
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	logger.Fatal(http.ListenAndServe(addr, nil))
}

func loadFilers() ([]*netapp.Filer, error) {
	if os.Getenv("DEV") != "" {
		logger.Info("Load filer configuration from env variables")
		return []*netapp.Filer{loadFilerFromEnv()}, nil
	} else {
		logger.Infof("Load filer configuration from %s", *configFile)
		return loadFilerFromFile(*configFile)
	}
}

func loadFilerFromFile(fileName string) (filers []*netapp.Filer, err error) {
	var yamlFile []byte
	if yamlFile, err = ioutil.ReadFile(fileName); err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, &filers); err != nil {
		return nil, err
	}
	for _, f := range filers {
		if f.Username == "" || f.Password == "" {
			username, password := loadAuthFromEnv()
			f.Username = username
			f.Password = password
		}
		// set netapp api version
		f.Version = "1.7"
	}
	return
}

func loadFilerFromEnv() *netapp.Filer {
	return &netapp.Filer{
		Name:             os.Getenv("NETAPP_NAME"),
		Host:             os.Getenv("NETAPP_HOST"),
		Username:         os.Getenv("NETAPP_USERNAME"),
		Password:         os.Getenv("NETAPP_PASSWORD"),
		AvailabilityZone: os.Getenv("NETAPP_AZ"),
		Version:          GetEnvWithDefaultValue("Netapp_API_VERSION", "1.7"),
	}
}

func loadAuthFromEnv() (username, password string) {
	username = os.Getenv("NETAPP_USERNAME")
	password = os.Getenv("NETAPP_PASSWORD")
	return
}

func GetEnvWithDefaultValue(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	} else {
		return defaultValue
	}
}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	s := fmt.Sprintf("%s [%s] %s\t", entry.Time.Format("2006-01-02 15:04:05.000"), entry.Level, entry.Message)
	for k, v := range entry.Data {
		s = s + fmt.Sprintf(" %s=%s", k, v)
	}
	s = s + "\n"
	return []byte(s), nil
}
