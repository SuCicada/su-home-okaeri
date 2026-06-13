package mqttentry

import (
	"encoding/json"
	"errors"
	"strings"
	"sucicada/home/internal/app/light"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/mqttpkg"
	deviceservice "sucicada/home/internal/service/devices"
	commandstructs "sucicada/home/internal/structs/command"
	devicesstructs "sucicada/home/internal/structs/devices"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type lightCommandPayload struct {
	State      string `json:"state"`
	Brightness int    `json:"brightness,omitempty"`
}

func RegisterLightRoutes(r *mqttpkg.Router) {
	for _, device := range deviceservice.GetDevices() {
		control := device.DeviceControl.Light
		if control == nil {
			continue
		}
		if _, ok := control.Controller.(devicesstructs.LightController); !ok {
			continue
		}

		topics := cfg.GetMqttConfig().Topics[control.MqttId]
		commandTopic := topics["command_topic"]
		stateTopic := topics["state_topic"]
		if commandTopic == "" {
			continue
		}

		r.Subscribe(commandTopic, handleLightCommand(device.Name, stateTopic))
		go publishLightState(device.Name, stateTopic)
	}
}

func handleLightCommand(deviceName string, stateTopic string) mqtt.MessageHandler {
	return func(client mqtt.Client, message mqtt.Message) {
		payload := lightCommandPayload{}
		logger.Info("Received light command: ", string(message.Payload()))

		if err := json.Unmarshal(message.Payload(), &payload); err != nil {
			logger.Error("Failed to unmarshal light command: ", err)
			return
		}

		if err := executeLightCommand(deviceName, &payload); err != nil {
			logger.Error("Failed to execute light command: ", err)
			return
		}

		publishLightState(deviceName, stateTopic)
		//if stateTopic != "" {
		//	mqttpkg.Publish(stateTopic, payload)
		//}
	}
}

func executeLightCommand(deviceName string, payload *lightCommandPayload) error {
	switch strings.ToUpper(payload.State) {
	case "OFF":
		//payload.Brightness = 0
		return light.Off(deviceName)
	case "ON":
		if payload.Brightness == 0 {
			return light.On(deviceName)
		}

		//default:
		command := commandstructs.LightCommand{Light: payload.Brightness}
		return light.SetBrightness(deviceName, command)
	}
	return errors.New("invalid light command: " + payload.State)
}

func publishLightState(deviceName string, stateTopic string) {
	if stateTopic == "" {
		return
	}

	brightness, err := light.GetBrightness(deviceName)
	if err != nil {
		logger.Error("Failed to get light brightness: ", err)
		return
	}

	payload := lightCommandPayload{
		State:      "ON",
		Brightness: brightness,
	}
	if brightness <= 0 {
		payload.State = "OFF"
	}

	mqttpkg.Publish(stateTopic, payload)
}
