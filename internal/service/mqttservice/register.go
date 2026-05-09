package mqttservice

import (
	"encoding/json"
	"fmt"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/mqttpkg"
	"sucicada/home/internal/structs"
	commandstructs "sucicada/home/internal/structs/command"
	devicesstructs "sucicada/home/internal/structs/devices"
	"sucicada/home/internal/util"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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

func RegisterMqttRouteMedia(control *devicesstructs.Control) {
	if control == nil {
		return
	}
	topics := cfg.GetMqttConfig().Topics[control.MqttId]

	commandTopic := topics["command_topic"]

	statusTopic := topics["status_topic"]
	isMutedTopic := topics["is_muted_topic"]
	volumeTopic := topics["volume_topic"]

	if commandTopic == "" {
		return
	}

	mqttpkg.RegisterRoute(commandTopic, func(client mqtt.Client, message mqtt.Message) {
		logger.Info("Received media command: ", string(message.Payload()))

		err := control.Control.Set(string(message.Payload()))
		if err != nil {
			logger.Error("Failed to execute media command: ", err)
			return
		}

		// internal/devices/windows/media.go
		getRes, err := control.Control.Get()
		if err != nil {
			logger.Error("Failed to get media status: ", err)
			return
		}
		var status structs.MediaStatus
		if getRes, ok := getRes.(structs.MediaStatus); ok {
			status = getRes
		} else {
			logger.Error("Failed to get media status: ", fmt.Sprintf("%v", getRes))
			return
		}

		if statusTopic != "" {
			mqttpkg.Publish(statusTopic, status)
		}
		if isMutedTopic != "" {
			mqttpkg.Publish(isMutedTopic, status.IsMute)
		}
		if volumeTopic != "" {
			mqttpkg.Publish(volumeTopic, status.Volume)
		}
	})
}
