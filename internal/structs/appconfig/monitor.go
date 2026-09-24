package appconfig

type MonitorConfig struct {
	HomeAssistant HomeAssistantConfig `yaml:"homeassistant"`
	ClaudeVM      ClaudeVMConfig      `yaml:"claude-vm"`
}

type HomeAssistantConfig struct {
	RestartScript string `yaml:"restart_script"`
}

type ClaudeVMConfig struct {
	RestartScript string `yaml:"restart_script"`
}
