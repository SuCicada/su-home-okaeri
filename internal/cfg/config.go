package cfg

import (
	"SuCicada/home/internal/logger"
	"SuCicada/home/internal/structs/appconfig"
	"os"

	"gopkg.in/yaml.v3"
)

// var (
// 	config *AppConfig
// 	// configOnce sync.Once
// )

// LoadConfig 加载YAML配置文件
func LoadConfig(configPath string) (*appconfig.AppConfig, error) {
	var err error
	// configOnce.Do(func() {
	var config = &appconfig.AppConfig{}

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

var preConfig *appconfig.AppConfig
var CONFIG_PATH = "config.yaml"

func GetConfig() *appconfig.AppConfig {
	var config, err = LoadConfig(CONFIG_PATH)
	if err != nil {
		logger.Error("load config error: ", err)
		return preConfig
	}
	preConfig = config
	return config
}
func GetSSHConfig(device string) appconfig.SSHConfig {
	return GetDeviceConfig(device).SSH
}

func GetDeviceConfig(device string) appconfig.DeviceConfig {
	return GetConfig().Devices[device]
}
