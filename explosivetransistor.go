package explosivetransistor

type Device interface {
	Set(string, int) error
	Get(string) (int, error)
	Name() string
}

type DevicePort int

type Portmap map[string]DevicePort
