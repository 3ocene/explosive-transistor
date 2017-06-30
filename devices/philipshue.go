package devices

import (
	"encoding/json"
	"errors"
	"fmt"
	et "github.com/ben-turner/explosive-transistor"
	"net/http"
)

type HueConfig struct {
	Ip    string
	Key   string
	Name  string
	Ports et.Portmap
}

type hueDevice struct {
	HueConfig
}

type hueState struct {
	On  bool `json: "on"`
	Bri int  `json: "bri"`
}

func NewHueDevice(c *HueConfig) (et.Device, error) {
	return &hueDevice{*c}, nil
}

func (d *hueDevice) Get(port string) (int, error) {
	res, err := http.Get(fmt.Sprintf("http://%v/api/%v/lights/%d", d.Ip, d.Key, d.Ports[port]))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var parsedRes struct {
		State hueState `json: "state"`
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&parsedRes); err != nil {
		return 0, err
	}

	fmt.Printf("%+v", parsedRes)

	if !parsedRes.State.On {
		return 0, nil
	}

	return parsedRes.State.Bri, nil
}

func (d *hueDevice) Set(port string, val int) error {
	return errors.New("Not implemented")
}

func (d *hueDevice) Name() string {
	return d.HueConfig.Name
}
