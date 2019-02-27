package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
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
	waitTime      = kingpin.Flag("wait", "Wait time").Short('w').Default("300").Int()
	configFile    = kingpin.Flag("config", "Config file").Short('f').Default("./netapp_filers.yaml").String()
	listenAddress = kingpin.Flag("listen", "Listen address").Short('l').Default("0.0.0.0").String()
)

var (
	netappCapacity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netapp",
			Subsystem: "capacity",
			Name:      "svm",
			Help:      "netapp SVM capacity",
		},
		[]string{
			"filer",
			"svm",
			"volume",
			"metric",
		},
	)
)

type filer struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"ip"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func main() {
	kingpin.Parse()

	var filers []filer
	if os.Getenv("Dev") != "" {
		filers = loadFilerFromEnv()
	} else {
		filers = loadFilerFromFile(*configFile)
	}

	p := prometheusExporter(filers)

	p.run()
	// vserverInfo := &netapp.VServerInfo{
	// 	VserverName:   "1",
	// 	UUID:          "1",
	// 	State:         "1",
	// 	AggregateList: &[]string{"x"},
	// }

	// volumeQuery:= &netapp.VolumeInfo{

	// }

	// volumeInfo := &netapp.VolumeInfo{
	// 	VolumeIDAttributes: &netapp.VolumeIDAttributes{
	// 		Name:              "x",
	// 		OwningVserverName: "x",
	// 	},
	// }

}

func loadFilerFromFile(fileName string) (c []filer) {
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

func loadFilerFromEnv() (c []filer) {
	c = append(c, filer{
		Name:     "test",
		Host:     os.Getenv("NETAPP_HOST"),
		Username: os.Getenv("NETAPP_USERNAME"),
		Password: os.Getenv("NETAPP_PASSWORD"),
	})

	return
}
