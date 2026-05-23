package devices

import (
	"fmt"
	"sucicada/home/internal/logger"
	devicesstructs "sucicada/home/internal/structs/devices"
)

//	func GetDeviceControlConfig(deviceName string) map[string]appconfig.DeviceControl {
//		device := cfg.GetConfig().Devices[deviceName]
//		return device.Control
//	}
//
//	func GetDevice(deviceName string) *devicesstructs.DeviceBase {
//		device, ok := DEVICES[deviceName]
//		if !ok {
//			logger.Warn("Device not found: ", deviceName)
//			return nil
//		}
//		return &device
//	}
var DEVICES = map[string]devicesstructs.DeviceBase{}

func GetDevices() map[string]devicesstructs.DeviceBase {
	devices := make(map[string]devicesstructs.DeviceBase, len(DEVICES))
	for name, device := range DEVICES {
		devices[name] = device
	}
	return devices
}

func RegisterDevice(device *devicesstructs.DeviceBase) {
	DEVICES[device.Name] = *device
	logger.Info("Registered device: ", device.Name)
}

func GetMediaController(deviceName string) (devicesstructs.MediaController, error) {
	device, ok := DEVICES[deviceName]
	if !ok {
		return nil, fmt.Errorf("device not found: %s", deviceName)
	}
	if device.DeviceControl.Media == nil {
		return nil, fmt.Errorf("device does not support media: %s", deviceName)
	}
	media, ok := device.DeviceControl.Media.Control.(devicesstructs.MediaController)
	if !ok {
		return nil, fmt.Errorf("device media controller is not supported: %s", deviceName)
	}
	return media, nil
}
