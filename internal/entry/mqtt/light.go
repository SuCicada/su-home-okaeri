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
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type lightCommandPayload struct {
	State      string `json:"state"`
	Brightness int    `json:"brightness,omitempty"`
}

func RegisterLightRoutes(r *mqttpkg.Router) {
	for _, device := range deviceservice.GetDevices() {
		logger.Info("Registering light routes for device: ", device.Name)
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
	publishStatus := func() { publishLightState(deviceName, stateTopic) }
	startStatusTicker(publishStatus, 10*time.Second)

	return func(client mqtt.Client, message mqtt.Message) {
		// paho 默认按顺序分发消息，回调返回后才会 ACK。
		// tvctl 要十几秒，再在这里 Publish().Wait() 会把整条 MQTT 接收卡住。
		payloadBytes := append([]byte(nil), message.Payload()...)
		go func() {
			payload := lightCommandPayload{}
			logger.Info("Received light command: ", string(payloadBytes))

			if err := json.Unmarshal(payloadBytes, &payload); err != nil {
				logger.Error("Failed to unmarshal light command: ", err)
				return
			}

			if err := executeLightCommand(deviceName, &payload); err != nil {
				logger.Error("Failed to execute light command: ", err)
				return
			}

			syncWithMultiTry(publishStatus)
		}()
	}
}

func executeLightCommand(deviceName string, payload *lightCommandPayload) error {
	logger.Info("Executing light command: ", payload.State, payload.Brightness)
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

	status, err := light.GetStatus(deviceName)
	if err != nil {
		logger.Error("Failed to get light status: ", err)
		return
	}

	state := "OFF"
	if status.Power {
		state = "ON"
	}
	payload := lightCommandPayload{
		State:      state,
		Brightness: status.Brightness,
	}
	logger.Debug("Publishing light state: ", payload)

	mqttpkg.Publish(stateTopic, payload)
}
