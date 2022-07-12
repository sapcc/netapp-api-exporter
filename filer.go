package main

import (
	"io/ioutil"
	"os"

	"github.com/sapcc/netapp-api-exporter/pkg/netapp"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

const netappApiVersion = "1.7"

type FilerBase struct {
	Name             string `yaml:"name"`
	Host             string `yaml:"host"`
	AvailabilityZone string `yaml:"availability_zone"`
	AggregatePattern string `yaml:"aggregate_pattern"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	Version          string `yaml:"version"`
}

type Filer struct {
	FilerBase
	Client *netapp.Client
}

func NewFiler(f FilerBase) (Filer, error) {
	c, err := netapp.NewClient(f.Host, f.Username, f.Password, f.Version)
	if err != nil {
		return Filer{}, err
	}
	return Filer{
		FilerBase: f,
		Client:    c,
	}, nil
}

func loadFilers(configFile string) ([]Filer, error) {
	if len(configFile) == 0 {
		log.Debug("load filer configuration from env variables")
		f, err := loadFilerFromEnv()
		if err != nil {
			return nil, err
		}
		return []Filer{f}, nil
	} else {
		log.Debugf("load filer configuration from %s", configFile)
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
		ff, err := NewFiler(*f)
		if err != nil {
			return nil, err
		}
		filers = append(filers, ff)
	}
	return
}

func loadFilerFromEnv() (Filer, error) {
	name := os.Getenv("NETAPP_NAME")
	host := os.Getenv("NETAPP_HOST")
	pattern := os.Getenv("NETAPP_AGGREGATE_PATTERN")
	username := os.Getenv("NETAPP_USERNAME")
	password := os.Getenv("NETAPP_PASSWORD")
	az := os.Getenv("NETAPP_AZ")
	version := getEnvWithDefaultValue("Netapp_API_VERSION", netappApiVersion)
	return NewFiler(FilerBase{name, host, az, pattern, username, password, version})
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
