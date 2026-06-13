package devices

import (
	"fmt"
	"sucicada/home/internal/logger"
	devicesstructs "sucicada/home/internal/structs/devices"
)

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
func GetDevice(deviceName string) devicesstructs.DeviceBase {
	return DEVICES[deviceName]
}

func GetMediaController(deviceName string) (devicesstructs.MediaController, error) {
	return getController[devicesstructs.MediaController](
		deviceName,
		"media",
		func(c devicesstructs.DeviceControlUnit) *devicesstructs.Control {
			return c.Media
		},
	)
}

func GetLightController(deviceName string) (devicesstructs.LightController, error) {
	return getController[devicesstructs.LightController](
		deviceName,
		"light",
		func(c devicesstructs.DeviceControlUnit) *devicesstructs.Control {
			return c.Light
		},
	)
}

func GetAudioController(deviceName string) (devicesstructs.AudioController, error) {
	return getController[devicesstructs.AudioController](
		deviceName,
		"audio",
		func(c devicesstructs.DeviceControlUnit) *devicesstructs.Control {
			return c.Audio
		},
	)
}

func getController[T any](
	deviceName string,
	capability string,
	selector func(devicesstructs.DeviceControlUnit) *devicesstructs.Control,
) (T, error) {
	var zero T

	device, ok := DEVICES[deviceName]
	if !ok {
		return zero, fmt.Errorf("device not found: %s", deviceName)
	}

	control := selector(device.DeviceControl)
	if control == nil {
		return zero, fmt.Errorf("device does not support %s: %s", capability, deviceName)
	}

	controller, ok := control.Controller.(T)
	if !ok {
		return zero, fmt.Errorf("device %s controller is not supported: %s", capability, deviceName)
	}

	return controller, nil
}
