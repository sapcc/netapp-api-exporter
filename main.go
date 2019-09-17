package main

import (
	"context"
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
	debug         = kingpin.Flag("debug", "Debug mode").Short('d').Bool()
	logger        = logrus.New()

	filers []*Filer
)

type myFormatter struct{}

func init() {
	kingpin.Parse()

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
	volumeGV := NewVolumeGaugeVec()
	aggrGV := NewAggrGaugeVec()

	for {
		filers = loadFilers()
		if len(filers) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, f := range filers {
		go fetchData(ctx, f)
		go processVolumes(ctx, f, volumeGV)
		go processAggregates(ctx, f, aggrGV)
	}

	prometheus.MustRegister(volumeGV)
	prometheus.MustRegister(aggrGV)
	http.Handle("/metrics", promhttp.Handler())
	logger.Fatal(http.ListenAndServe(*listenAddress+":9108", nil))
}

func fetchData(ctx context.Context, f *Filer) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		f.GetNetappVolume(f.volChan, f.getVolDone)
		f.GetNetappAggregate(f.aggrChan, f.getAggrDone)
		time.Sleep(time.Duration(*sleepTime) * time.Second)
	}
}

func processVolumes(ctx context.Context, f *Filer, gv VolumeGaugeVec) {
	rcvdVolumes := make(map[string]*NetappVolume)
	volumes := make(map[string]bool)
	for {
		select {
		case v := <-f.volChan:
			logger.Debugf("[%s] Volume %s received: %s", f.Name, v.ShareID, v.SizeUsed)
			gv.SetMetric(v)
			rcvdVolumes[v.ShareID] = v
			volumes[v.ShareID] = true
		case <-f.getVolDone:
			for shareID, ok := range volumes {
				if !ok {
					gv.DeleteMetric(rcvdVolumes[shareID])
					delete(rcvdVolumes, shareID)
					delete(volumes, shareID)
					logger.Debugf("[%s] Volume %s deleted", f.Name, shareID)
				}
			}
			for shareID, _ := range volumes {
				volumes[shareID] = false
			}
		case <-ctx.Done():
			return
		}
	}
}

func processAggregates(ctx context.Context, f *Filer, gv AggrGaugeVec) {
	for {
		select {
		case ag := <-f.aggrChan:
			logger.Debugf("[%s] Aggregate %s received", f.Name, ag.Name)
			gv.SetMetric(ag)
		case <-f.getAggrDone:
		case <-ctx.Done():
			return
		}
	}
}

func loadFilerFromFile(fileName string) (c []*Filer) {
	username, password := loadAuthFromEnv()

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
		if b.Username == "" || b.Password == "" {
			c = append(c, NewFiler(b.Name, b.Host, username, password, b.AvailabilityZone))
		} else {
			c = append(c, NewFiler(b.Name, b.Host, b.Username, b.Password, b.AvailabilityZone))
		}
	}
	return
}

func loadFilers() (filers []*Filer) {
	if os.Getenv("DEV") != "" {
		*debug = true
		filers = loadFilerFromEnv()
	} else {
		filers = loadFilerFromFile(*configFile)
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
