package devices

import (
	"errors"
	"github.com/Sirupsen/logrus"
	et "github.com/ben-turner/explosive-transistor"
	"github.com/tarm/serial"
)

const (
	arduinoRead  byte = 0x00
	arduinoWrite byte = 0x01
)

var arduinoPrefix = []byte{0x02, 0x02, 0x02}

type ArduinoConfig struct {
	BaudRate   int
	SerialPort string
	Ports      et.Portmap
	Name       string
}

type arduinoDevice struct {
	ArduinoConfig

	Port *serial.Port
}

// Creates a new explosivetransistor.Device representing an Arduino
func NewArduinoDevice(c *ArduinoConfig) (et.Device, error) {
	conf := &serial.Config{Name: c.SerialPort, Baud: c.BaudRate}
	s, err := serial.OpenPort(conf)
	if err != nil {
		return nil, err
	}

	return &arduinoDevice{
		*c,
		s,
	}, nil
}

func (d *arduinoDevice) write(data []byte) error {
	logrus.WithField("data", data).WithField("device", d.Name()).Info("Writing to arduino")
	l, err := d.Port.Write(data)
	if err != nil {
		return err
	}
	if l != len(data) {
		return errors.New("Wrong number of bytes written")
	}
	return nil
}

// Sets the specified port on the arduino to the specified value
func (d *arduinoDevice) Set(port string, val int) error {
	logrus.WithField("port", port).WithField("value", val).Info("Setting arduino pin")
	msg := append(arduinoPrefix, arduinoWrite, byte(d.Ports[port]), byte(val))
	err := d.write(msg)
	return err
}

// Gets the current state of the specified port
func (d *arduinoDevice) Get(port string) (val int, err error) {
	logrus.WithField("port", port).Info("Getting arduino pin")
	msg := append(arduinoPrefix, arduinoRead, byte(d.Ports[port]))
	err = d.write(msg)
	if err != nil {
		return
	}

	buf := make([]byte, 1)
	_, err = d.Port.Read(buf)
	if err != nil {
		return
	}
	return int(buf[0]), nil
}

func (d *arduinoDevice) Name() string {
	return d.ArduinoConfig.Name
}
