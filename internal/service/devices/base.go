package devices

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/logger"
	"SuCicada/home/internal/structs/appconfig"
)

func GetDeviceControlConfig(deviceName string) appconfig.DeviceControl {
	device := cfg.GetConfig().Devices[deviceName]
	return device.Control
}
func GetDevice(deviceName string) *DeviceBase {
	device, ok := devices[deviceName]
	if !ok {
		logger.Warn("Device not found: ", deviceName)
		return nil
	}
	return &device
}

var devices = map[string]DeviceBase{}

func RegisterDevice(device *DeviceBase) {
	devices[device.Name] = *device
	logger.Info("Registered device: ", device.Name)
}
