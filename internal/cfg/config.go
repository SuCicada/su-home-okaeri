package cfg

import (
	"SuCicada/home/internal/logger"
	"os"

	"gopkg.in/yaml.v3"
)

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
	High int `yaml:"high"`
	Low  int `yaml:"low"`
}

type SSHConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

// 应用配置结构体
type AppConfig struct {
	// SSH struct {
	// 	Linux   SSHHostConfig `yaml:"linux"`
	// 	Redmi   SSHHostConfig `yaml:"redmi"`
	// 	Windows SSHHostConfig `yaml:"windows"`
	// } `yaml:"ssh"`
	Devices map[string]DeviceConfig `yaml:"devices"`
}

// var (
// 	config *AppConfig
// 	// configOnce sync.Once
// )

// LoadConfig 加载YAML配置文件
func LoadConfig(configPath string) (*AppConfig, error) {
	var err error
	// configOnce.Do(func() {
	var config = &AppConfig{}

	yamlFile, _err := os.ReadFile(configPath)
	err = _err
	if err != nil {
		logger.Error("yaml file error: ", err)
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		logger.Error("yaml unmarshal error: ", err)
		return nil, err
	}
	logger.Info("yaml load success: ", configPath)
	// })
	return config, nil
}

var preConfig *AppConfig

func GetConfig() *AppConfig {
	var config, err = LoadConfig("config.yaml")
	if err != nil {
		logger.Error("load config error: ", err)
		return preConfig
	}
	preConfig = config
	return config
}
func GetSSHConfig(device string) SSHConfig {
	return GetDeviceConfig(device).SSH
}

func GetDeviceConfig(device string) DeviceConfig {
	return GetConfig().Devices[device]
}
