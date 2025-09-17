package cfg

import (
	"SuCicada/home/internal/logger"
	"os"

	"gopkg.in/yaml.v3"
)

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
