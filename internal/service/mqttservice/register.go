package mqttservice

import (
	"encoding/json"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/mqttpkg"
	deviceservice "sucicada/home/internal/service/devices"
	commandstructs "sucicada/home/internal/structs/command"
	devicesstructs "sucicada/home/internal/structs/devices"
	"sucicada/home/internal/util"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func RegisterRoutes() {
	for _, device := range deviceservice.GetDevices() {
		RegisterMqttRouteLight(device.DeviceControl.Light)
	}
}

func RegisterMqttRouteLight(control *devicesstructs.Control) {
	if control == nil {
		return
	}

	// type MqttTopicLight struct {
	// CommandTopic string `yaml:"command_topic"`
	// StateTopic   string `yaml:"state_topic"`
	// }

	topics := cfg.GetMqttConfig().Topics[control.MqttId]

	var commandTopic string = topics["command_topic"]
	var stateTopic string = topics["state_topic"]

	if commandTopic == "" {
		return
	}

	type MqttPayload struct {
		State      string `json:"state"` // ON or OFF
		Brightness int    `json:"brightness,omitempty"`
	}

	mqttpkg.RegisterRoute(commandTopic, func(client mqtt.Client, message mqtt.Message) {
		payload := MqttPayload{}
		logger.Info("Received message: ", string(message.Payload()))

		json.Unmarshal(message.Payload(), &payload)

		command := util.Conv.ToJsonStr(commandstructs.LightCommand{Light: payload.Brightness})
		control.Control.Set(command)

		if payload.State == "OFF" {
			payload.Brightness = 0
		}

		mqttpkg.Publish(stateTopic, payload)
	})
}
