package devices

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/logger"
	"SuCicada/home/internal/mqttpkg"
	"SuCicada/home/internal/structs/appconfig"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

	RegisterMqttRoute(device.DeviceControl.Light)
	RegisterMqttRoute(device.DeviceControl.Volume)

}

func RegisterMqttRoute(control *Control) {
	if control == nil {
		return
	}
	topics := cfg.GetMqttConfig().Topics[control.MqttId]

	if topics.CommandTopic == "" {
		return
	}

	type MqttPayload struct {
		State      string `json:"state"` // ON or OFF
		Brightness int    `json:"brightness,omitempty"`
	}

	mqttpkg.RegisterRoute(topics.CommandTopic, func(client mqtt.Client, message mqtt.Message) {
		payload := MqttPayload{}
		logger.Info("Received message: ", string(message.Payload()))

		json.Unmarshal(message.Payload(), &payload)

		control.Control.Set(payload.Brightness)

		// if payload.Brightness == 0 {
		if payload.State == "OFF" {
			payload.Brightness = 0
		}
		// }

		mqttpkg.Publish(topics.StateTopic, payload)
	})
}
