package cfg

type SMSCheckConfig struct {
	PushUrl string `yaml:"push_url"`
	Secret  string `yaml:"secret"`
}
