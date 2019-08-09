package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

// Parameter
var (
	sleepTime     = kingpin.Flag("wait", "Wait time").Short('w').Default("300").Int64()
	configFile    = kingpin.Flag("config", "Config file").Short('c').Default("./netapp_filers.yaml").String()
	listenAddress = kingpin.Flag("listen", "Listen address").Short('l').Default("0.0.0.0").String()
	logger        = logrus.New()

	filers []*Filer
)

type myFormatter struct{}

func init() {
	kingpin.Parse()

	logger.Out = os.Stdout
	logger.SetFormatter(new(myFormatter))
	if os.Getenv("DEV") != "" {
		logger.Level = logrus.DebugLevel
		filers = loadFilerFromEnv()
	} else {
		logger.Level = logrus.InfoLevel
		filers = loadFilerFromFile(*configFile)
	}
	for _, f := range filers {
		logger.Printf("Host (%s) loaded", f.Host)
	}
}

func main() {

	p := NewCapacityExporter()
	for _, f := range filers {
		f.Init()
		go p.runGetNetappShare(f, time.Duration(*sleepTime))
		go p.runGetNetappAggregate(f, time.Duration(*sleepTime))
	}

	prometheus.MustRegister(p.volumeCollector)
	prometheus.MustRegister(p.aggregateCollector)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*listenAddress+":9108", nil)
}

func loadFilerFromFile(fileName string) (c []*Filer) {
	var fb []FilerBase
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Fatal("[ERROR] ", err)
	}
	err = yaml.Unmarshal(yamlFile, &fb)
	if err != nil {
		logger.Fatal("[ERROR] ", err)
	}
	for _, b := range fb {
		c = append(c, &Filer{FilerBase: b})
	}
	return
}

func loadFilerFromEnv() (c []*Filer) {
	host := os.Getenv("NETAPP_HOST")
	username := os.Getenv("NETAPP_USERNAME")
	password := os.Getenv("NETAPP_PASSWORD")
	az := os.Getenv("NETAPP_AZ")
	f := NewFiler("test", host, username, password, az)
	c = append(c, f)
	return
}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	s := fmt.Sprintf("%s [%s] %s\t", entry.Time.Format("2006-01-02 15:04:05.000"), entry.Level, entry.Message)
	for k, v := range entry.Data {
		s = s + fmt.Sprintf(" %s=%s", k, v)
	}
	s = s + "\n"
	return []byte(s), nil
}
