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

	controller, ok := control.Control.(T)
	if !ok {
		return zero, fmt.Errorf("device %s controller is not supported: %s", capability, deviceName)
	}

	return controller, nil
}

// func GetMediaController(deviceName string) (devicesstructs.MediaController, error) {
// 	device, ok := DEVICES[deviceName]
// 	if !ok {
// 		return nil, fmt.Errorf("device not found: %s", deviceName)
// 	}
// 	if device.DeviceControl.Media == nil {
// 		return nil, fmt.Errorf("device does not support media: %s", deviceName)
// 	}
// 	media, ok := device.DeviceControl.Media.Control.(devicesstructs.MediaController)
// 	if !ok {
// 		return nil, fmt.Errorf("device media controller is not supported: %s", deviceName)
// 	}
// 	return media, nil
// }

// func GetLightController(deviceName string) (devicesstructs.LightController, error) {
// 	device, ok := DEVICES[deviceName]
// 	if !ok {
// 		return nil, fmt.Errorf("device not found: %s", deviceName)
// 	}
// 	if device.DeviceControl.Light == nil {
// 		return nil, fmt.Errorf("device does not support light: %s", deviceName)
// 	}
// 	light, ok := device.DeviceControl.Light.Control.(devicesstructs.LightController)
// 	if !ok {
// 		return nil, fmt.Errorf("device light controller is not supported: %s", deviceName)
// 	}
// 	return light, nil
// }

// func GetAudioController(deviceName string) (devicesstructs.AudioController, error) {
// 	device, ok := DEVICES[deviceName]
// 	if !ok {
// 		return nil, fmt.Errorf("device not found: %s", deviceName)
// 	}
// 	if device.DeviceControl.Audio == nil {
// 		return nil, fmt.Errorf("device does not support audio: %s", deviceName)
// 	}
// 	audio, ok := device.DeviceControl.Audio.Control.(devicesstructs.AudioController)
// 	if !ok {
// 		return nil, fmt.Errorf("device audio controller is not supported: %s", deviceName)
// 	}
// 	return audio, nil
// }
