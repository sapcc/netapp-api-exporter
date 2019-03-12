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

// Parameter
var (
	sleepTime     = kingpin.Flag("wait", "Wait time").Short('w').Default("300").Int64()
	configFile    = kingpin.Flag("config", "Config file").Short('c').Default("./netapp_filers.yaml").String()
	listenAddress = kingpin.Flag("listen", "Listen address").Short('l').Default("0.0.0.0").String()
)

func main() {
	var filers []*Filer

	kingpin.Parse()

	if os.Getenv("DEV") != "" {
		filers = loadFilerFromEnv()
	} else {
		filers = loadFilerFromFile(*configFile)
	}

	p := NewCapacityExporter()

	for _, f := range filers {
		f.Init()
		go p.runGetNetappShare(f, time.Duration(*sleepTime))
	}

	prometheus.MustRegister(p.collector)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*listenAddress+":9108", nil)
}

func loadFilerFromFile(fileName string) (c []*Filer) {
	var fb []FilerBase
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("[ERROR] ", err)
	}
	err = yaml.Unmarshal(yamlFile, &fb)
	if err != nil {
		log.Fatal("[ERROR] ", err)
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
	f := NewFiler("test", host, username, password)
	c = append(c, f)
	return
}
