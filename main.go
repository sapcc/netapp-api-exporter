package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"os"

	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

// Parameter
var (
	configFile    = kingpin.Flag("config", "Config file").Short('c').Default("./netapp_filers.yaml").String()
	listenAddress = kingpin.Flag("listen", "Listen address").Short('l').Default("0.0.0.0").String()
	debug         = kingpin.Flag("debug", "Debug mode").Short('d').Bool()
	logger        = logrus.New()

	filers []Filer
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
	for _, f := range filers {
		logger.Printf("Host (%s) loaded", f.Host)
	}
}

func main() {
	// try loading filers every 5 seconds until successful
	for {
		filers = loadFilers()
		if len(filers) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	reg := prometheus.NewPedanticRegistry()

	for _, f := range filers {
		logger.Println("Register filer: Name:", f.FilerBase.Name, "Host:", f.FilerBase.Host, "Username:", f.FilerBase.Username, "AvailabilityZone:", f.FilerBase.AvailabilityZone)
		cc := NewFilerCollector(f)
		labels := prometheus.Labels{
			"filer":             f.Name,
			"availability_zone": f.AvailabilityZone,
		}
		prometheus.WrapRegistererWith(labels, reg).MustRegister(cc)
	}

	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	logger.Fatal(http.ListenAndServe(*listenAddress+":9108", nil))
}

func loadFilerFromFile(fileName string) (c []Filer) {
	var filers []FilerBase
	if yamlFile, err := ioutil.ReadFile(fileName); err != nil {
		logger.Fatal("read file ", fileName, err)
	} else {
		if err := yaml.Unmarshal(yamlFile, &filers); err != nil {
			logger.Fatal("unmarshal yaml struct", err)
		}
	}

	for _, f := range filers {
		if f.Username == "" || f.Password == "" {
			username, password := loadAuthFromEnv()
			f.Username = username
			f.Password = password
		}
		c = append(c, NewFiler(f))
	}
	return
}

func loadFilers() (filers []Filer) {
	if os.Getenv("DEV") != "" {
		filers = loadFilerFromEnv()
	} else {
		filers = loadFilerFromFile(*configFile)
	}
	return
}

func loadFilerFromEnv() (c []Filer) {
	name := os.Getenv("NETAPP_NAME")
	host := os.Getenv("NETAPP_HOST")
	username := os.Getenv("NETAPP_USERNAME")
	password := os.Getenv("NETAPP_PASSWORD")
	az := os.Getenv("NETAPP_AZ")
	f := NewFiler(FilerBase{name, host, username, password, az})
	c = append(c, f)
	return
}

func loadAuthFromEnv() (username, password string) {
	username = os.Getenv("NETAPP_USERNAME")
	password = os.Getenv("NETAPP_PASSWORD")
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
