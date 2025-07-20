package devices

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/logger"
)

func GetDeviceControlConfig(deviceName string) cfg.DeviceControl {
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
