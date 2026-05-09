package devices

import (
	"sucicada/home/internal/logger"
	"sucicada/home/internal/service/mqttservice"
	devicesstructs "sucicada/home/internal/structs/devices"
)

//	func GetDeviceControlConfig(deviceName string) map[string]appconfig.DeviceControl {
//		device := cfg.GetConfig().Devices[deviceName]
//		return device.Control
//	}
//
//	func GetDevices() map[string]devicesstructs.DeviceBase {
//		return DEVICES
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

func RegisterDevice(device *devicesstructs.DeviceBase) {
	DEVICES[device.Name] = *device
	logger.Info("Registered device: ", device.Name)

	mqttservice.RegisterMqttRouteLight(device.DeviceControl.Light)
	mqttservice.RegisterMqttRouteMedia(device.DeviceControl.Media)
}
