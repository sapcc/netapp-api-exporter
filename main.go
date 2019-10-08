package main

import (
	"context"
	"fmt"
	"github.com/pepabo/go-netapp/netapp"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

// Parameter
var (
	sleepTime     = kingpin.Flag("wait", "Wait time").Short('w').Default("300").Float64()
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
			logger.Debugf("[%s] NetappAggregate %s received", f.Name, ag.Name)
			gv.SetMetric(ag)
		case <-f.getAggrDone:
		case <-ctx.Done():
			return
		}
	}
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
		if b.Username == "" || b.Password == "" {
			username, password := loadAuthFromEnv()
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

//volumes    []*NetappVolume
//aggregates []*NetappAggregate

type FilerCollector struct {
	*Filer
	volumeList []*NetappVolume
	mux        sync.Mutex

	volumeSizeTotal                    *prometheus.Desc
	volumeSizeAvail                    *prometheus.Desc
	volumeSizeUsed                     *prometheus.Desc
	volumeSizeUsedBySnapshots          *prometheus.Desc
	volumeSizeAvailForSnapshots        *prometheus.Desc
	volumeSizeReservedForSnapshots     *prometheus.Desc
	volumePercentageUsed               *prometheus.Desc
	volumePercentageCompressionSaved   *prometheus.Desc
	volumePercentageDeduplicationSaved *prometheus.Desc
	volumePercentageTotalSaved         *prometheus.Desc
	aggregateSizeUsed                  *prometheus.Desc
	aggregateSizeTotal                 *prometheus.Desc
	aggregateSizeAvail                 *prometheus.Desc
	aggregateSizeTotalReserved         *prometheus.Desc
	aggregatePercentUsed               *prometheus.Desc
	aggregatePhysicalUsed              *prometheus.Desc
	aggregatePercentPhysicalUsed       *prometheus.Desc

	up, scrapeDuration prometheus.Gauge
	scrapesTotal       prometheus.Counter
	lastScrapeTime     time.Time
}

func NewFilerCollector(f *Filer) prometheus.Collector {
	volumeLabels := []string{
		"filer",
		"vserver",
		"volume",
		"project_id",
		"share_id",
	}
	aggregateLabels := []string{
		"availability_zone",
		"filer",
		"node",
		"aggregate",
	}
	c := &FilerCollector{
		//Filer: f,
		volumeSizeTotal: prometheus.NewDesc(
			"netapp_capacity_svm",
			"Netapp Volume Metrics: total size",
			volumeLabels,
			prometheus.Labels{"metric": "size_total"},
		),
		aggregateSizeAvail: prometheus.NewDesc(
			"netapp_capacity_aggregate",
			"Netapp Aggregate Metrics",
			aggregateLabels,
			nil,
		),
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "netapp",
			Subsystem: f.Name,
			Name:      "up",
			Help:      "'1' if the last scrape of filer was successful, '0' otherwise.",
		}),
		scrapeDuration: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "netapp",
			Subsystem: f.Name,
			Name:      "scrape_duration_seconds",
			Help:      "The duration it took to scrape filer.",
		}),
		scrapesTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "netapp",
			Subsystem: f.Name,
			Name:      "scrape_total",
			Help:      "The total number of filer scrapes.",
		}),
		lastScrapeTime: time.Time{},
	}

	return c
}

func (c *FilerCollector) Collect(chan<- prometheus.Metric) {
	var err error
	if time.Now().Sub(c.lastScrapeTime).Seconds() > *sleepTime {
		c.lastScrapeTime = time.Now()
		c.mux.Lock()
		c.volumeList, err = c.Filer.GetNetappVolume()
		c.mux.Unlock()
		if err != nil {
			c.up.Set(0)
			return
		}
		for _, v := range c.volumeList {
			labels := []string{c.Filer.Name, v.Vserver, v.Volume, v.ProjectID, v.ShareID}

			sizeTotal, _ := strconv.ParseFloat(v.SizeTotal, 64)
			sizeAvailable, _ := strconv.ParseFloat(v.SizeAvailable, 64)
			sizeUsed, _ := strconv.ParseFloat(v.SizeUsed, 64)
			sizeUsedBySnapshots, _ := strconv.ParseFloat(v.SizeUsedBySnapshots, 64)
			sizeAvailableForSnapshots, _ := strconv.ParseFloat(v.SizeAvailableForSnapshots, 64)
			snapshotReserveSize, _ := strconv.ParseFloat(v.SnapshotReserveSize, 64)
			percentageSizeUsed, _ := strconv.ParseFloat(v.PercentageSizeUsed, 64)
			percentageCompressionSpaceSaved, _ := strconv.ParseFloat(v.PercentageCompressionSpaceSaved, 64)
			percentageDeduplicationSpaceSaved, _ := strconv.ParseFloat(v.PercentageDeduplicationSpaceSaved, 64)
			percentageTotalSpaceSaved, _ := strconv.ParseFloat(v.PercentageTotalSpaceSaved, 64)
			ch <- prometheus.MustNewConstMetric(c.volumeSizeTotal, prometheus.GaugeValue, sizeTotal, labels...)
			ch <- prometheus.MustNewConstMetric(c.volumeSizeAvail, prometheus.GaugeValue, sizeAvailable, labels...)
			ch <- prometheus.MustNewConstMetric(c.volumeSizeUsed, prometheus.GaugeValue, sizeUsed, labels...)
			ch <- prometheus.MustNewConstMetric(c.volumeSizeUsedBySnapshots, prometheus.GaugeValue, sizeUsedBySnapshots, labels...)
			ch <- prometheus.MustNewConstMetric(c.volumeSizeAvailForSnapshots, prometheus.GaugeValue, sizeAvailableForSnapshots, labels...)
			ch <- prometheus.MustNewConstMetric(c.volumeSizeReservedForSnapshots, prometheus.GaugeValue, snapshotReserveSize, labels...)
			ch <- prometheus.MustNewConstMetric(c.volumePercentageUsed, prometheus.GaugeValue, percentageSizeUsed, labels...)
			ch <- prometheus.MustNewConstMetric(c.volumePercentageTotalSaved, prometheus.GaugeValue, percentageTotalSpaceSaved, labels...)
			ch <- prometheus.MustNewConstMetric(c.volumePercentageCompressionSaved, prometheus.GaugeValue, percentageCompressionSpaceSaved, labels...)
			ch <- prometheus.MustNewConstMetric(c.volumePercentageDeduplicationSaved, prometheus.GaugeValue, percentageDeduplicationSpaceSaved, labels...)
		}
	}
}

func (c *FilerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.volumeSizeTotal
	ch <- c.volumeSizeAvail
	ch <- c.volumeSizeUsed
	ch <- c.volumeSizeUsedBySnapshots
	ch <- c.volumeSizeAvailForSnapshots
	ch <- c.volumeSizeReservedForSnapshots
	ch <- c.volumePercentageUsed
	ch <- c.volumePercentageCompressionSaved
	ch <- c.volumePercentageDeduplicationSaved
	ch <- c.volumePercentageTotalSaved
	ch <- c.aggregateSizeUsed
	ch <- c.aggregateSizeTotal
	ch <- c.aggregateSizeAvail
	ch <- c.aggregateSizeTotalReserved
	ch <- c.aggregatePercentUsed
	ch <- c.aggregatePhysicalUsed
	ch <- c.aggregatePercentPhysicalUsed
	ch <- c.up.Desc()
	ch <- c.scrapesTotal.Desc()
	ch <- c.scrapeDuration.Desc()
}
