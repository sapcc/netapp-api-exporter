package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	url      string
	username string
	password string
)

const (
	_url    = "https://%s/servlets/netapp.servlets.admin.XMLrequest_filer"
	version = "1.7"
)

// Parameter
var (
	sleepTime     = kingpin.Flag("wait", "Wait time").Short('w').Default("300").Int64()
	configFile    = kingpin.Flag("config", "Config file").Short('f').Default("./netapp_filers.yaml").String()
	listenAddress = kingpin.Flag("listen", "Listen address").Short('l').Default("0.0.0.0").String()
)

type filer struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"ip"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func main() {
	var filers []*filer

	kingpin.Parse()

	if os.Getenv("Dev") != "" {
		filers = loadFilerFromEnv()
	} else {
		filers = loadFilerFromFile(*configFile)
	}

	p := NewCapacityExporter()

	for _, f := range filers {
		go p.run(f, time.Duration(*sleepTime))
	}

	prometheus.MustRegister(p.collector)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*listenAddress+":9108", nil)
}

func loadFilerFromFile(fileName string) (c []*filer) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("[ERROR] ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatal("[ERROR] ", err)
	}
	return
}

func loadFilerFromEnv() (c []*filer) {
	c = append(c, &filer{
		Name:     "test",
		Host:     os.Getenv("NETAPP_HOST"),
		Username: os.Getenv("NETAPP_USERNAME"),
		Password: os.Getenv("NETAPP_PASSWORD"),
	})
	return
}
