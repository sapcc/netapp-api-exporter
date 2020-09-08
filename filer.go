package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const netappApiVersion = "1.7"

type Filer struct {
	Name             string `yaml:"name"`
	Host             string `yaml:"host"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	AvailabilityZone string `yaml:"availability_zone"`
	Version          string `yaml:"version"`
}

func loadFilers() ([]*Filer, error) {
	if os.Getenv("DEV") != "" {
		logger.Info("Load filer configuration from env variables")
		return []*Filer{loadFilerFromEnv()}, nil
	} else {
		logger.Infof("Load filer configuration from %s", *configFile)
		return loadFilerFromFile(*configFile)
	}
}

func loadFilerFromFile(fileName string) (filers []*Filer, err error) {
	var yamlFile []byte
	if yamlFile, err = ioutil.ReadFile(fileName); err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, &filers); err != nil {
		return nil, err
	}
	for _, f := range filers {
		if f.Username == "" || f.Password == "" {
			username, password := getAuthFromEnv()
			f.Username = username
			f.Password = password
		}
		if f.Version == "" {
			f.Version = netappApiVersion
		}
	}
	return
}

func loadFilerFromEnv() *Filer {
	return &Filer{
		Name:             os.Getenv("NETAPP_NAME"),
		Host:             os.Getenv("NETAPP_HOST"),
		Username:         os.Getenv("NETAPP_USERNAME"),
		Password:         os.Getenv("NETAPP_PASSWORD"),
		AvailabilityZone: os.Getenv("NETAPP_AZ"),
		Version:          getEnvWithDefaultValue("Netapp_API_VERSION", netappApiVersion),
	}
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
