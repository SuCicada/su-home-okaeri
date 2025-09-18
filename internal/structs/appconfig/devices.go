package appconfig

type DeviceConfig struct {
	SSH     SSHConfig     `yaml:"ssh"`
	Control DeviceControl `yaml:"control"`
	// Value struct {
	// Light  int `yaml:"light"`
	// Volume int `yaml:"volume"`
	// } `yaml:"value"`
}
type DeviceControl map[string]Value
type Value struct {
	Options map[string]any `yaml:"options"`
	High    int            `yaml:"high"`
	Low     int            `yaml:"low"`
}

type SSHConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}
