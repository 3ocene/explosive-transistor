package main

import (
	"errors"
	et "github.com/ben-turner/explosive-transistor"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type DeviceConfig struct {
	Type    string                 `yaml: "type"`
	Config  map[string]interface{} `yaml: "config"`
	Portmap et.Portmap             `yaml: "portmap"`
}

type ApiConfig struct {
	Address string `yaml: "address"`
}

type Config struct {
	Devices map[string]DeviceConfig `yaml: "devices"`
	Api     ApiConfig               `yaml: "api"`
}

func loadConfig() (*Config, error) {
	configFilename := os.Getenv("EXPLOSIVE_TRANSISTOR_CONFIG")
	if configFilename == "" {
		fileStats, err := os.Stat("./config.yml")
		if err != nil {
			return nil, err
		}
		if fileStats.IsDir() {
			return nil, errors.New("Config file is a directory")
		}
		configFilename = "./config.yml"
	}

	configFile, err := os.Open(configFilename)
	if err != nil {
		return nil, err
	}

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	yaml.Unmarshal(configBytes, c)

	return c, nil
}
