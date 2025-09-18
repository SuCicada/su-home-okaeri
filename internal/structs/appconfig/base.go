package appconfig

type AppConfig struct {
	Devices  map[string]DeviceConfig `yaml:"devices"`
	SMSCheck SMSCheckConfig          `yaml:"sms_check"`
	Alert    AlertConfig             `yaml:"alert"`
}
