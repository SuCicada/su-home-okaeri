package appconfig

type MonitorConfig struct {
	HomeAssistant HomeAssistantConfig `yaml:"homeassistant"`
}

type HomeAssistantConfig struct {
	RestartScript string `yaml:"restart_script"`
}
