package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	et "github.com/ben-turner/explosive-transistor"
	"github.com/ben-turner/explosive-transistor/devices"
	"github.com/gorilla/mux"
	"github.com/kr/pretty"
	"net/http"
	"os"
	"strconv"
)

type api struct {
	Config

	devices map[string]et.Device
}

func newApi(c *Config) (a api, err error) {
	a = api{
		*c,
		make(map[string]et.Device),
	}

	for key, deviceConfig := range c.Devices {
		switch deviceConfig.Type {
		case "arduino":
			a.devices[key], err = devices.NewArduinoDevice(&devices.ArduinoConfig{
				BaudRate:   deviceConfig.Config["baudRate"].(int),
				SerialPort: deviceConfig.Config["serialPort"].(string),
				Ports:      deviceConfig.Portmap,
				Name:       key,
			})
		default:
			err = fmt.Errorf("No device of type %v available", deviceConfig.Type)
		}
		if err != nil {
			return
		}
	}
	return
}

func (a api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	dev, ok := a.devices[vars["device"]]
	if !ok {
		fmt.Fprint(w, "Device not found")
		logrus.WithField("device", vars["device"]).Info("Request to unavailable device")
		return
	}

	switch vars["action"] {
	case "get":
		val, err := dev.Get(vars["port"])
		if err != nil {
			fmt.Fprint(w, "Could not talk to device")
			logrus.WithError(err).
				WithField("action", vars["action"]).
				WithField("device", vars["device"]).
				WithField("port", vars["port"]).
				WithField("action", vars["action"]).
				Info("Failed to set device port")
			return
		}
		fmt.Fprintf(w, "%v.%v = %d", dev.Name(), vars["port"], val)
		return
	case "set":
		val, err := strconv.Atoi(vars["value"])
		if err != nil {
			fmt.Fprint(w, "Could not parse query value")
			logrus.WithError(err).
				WithField("action", vars["action"]).
				WithField("device", vars["device"]).
				WithField("port", vars["port"]).
				WithField("action", vars["action"]).
				Info("Could not parse query value")
			return
		}

		if err := dev.Set(vars["port"], val); err != nil {
			fmt.Fprint(w, "Could not talk to device")
			logrus.WithError(err).
				WithField("action", vars["action"]).
				WithField("device", vars["device"]).
				WithField("port", vars["port"]).
				WithField("action", vars["action"]).
				Info("Failed to get device port")
			return
		}
		fmt.Fprint(w, "ok")
	default:
		fmt.Fprint(w, "Action does not exist")
		logrus.WithField("action", vars["action"]).
			WithField("device", vars["device"]).
			WithField("port", vars["port"]).
			WithField("action", vars["action"]).
			Info("Bad action called")
		return
	}
}

func run() int {
	logrus.Info("Starting Explosive Transistor API")

	conf, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Info("Failed to load config")
		return 1
	}

	fmt.Fprintf(os.Stderr, "%# v", pretty.Formatter(*conf))

	a, err := newApi(conf)
	if err != nil {
		logrus.WithError(err).Info("Could not create API object")
		return 1
	}

	r := mux.NewRouter()
	r.Handle("/{device}/{action}/{port}/{value}", a)

	logrus.WithError(http.ListenAndServe(a.Api.Address, r)).Info("API server exited")
	return 0
}

func main() {
	os.Exit(run())
}
