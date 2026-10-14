package mqttentry

import (
	"encoding/json"
	"sucicada/home/internal/app/media"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/mqttpkg"

	deviceservice "sucicada/home/internal/service/devices"
	"sucicada/home/internal/structs"
	devicesstructs "sucicada/home/internal/structs/devices"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mediaTopics struct {
	Command string
	Status  string
}

func RegisterMediaRoutes(r *mqttpkg.Router) {
	for _, device := range deviceservice.GetDevices() {
		control := device.DeviceControl.Media
		if control == nil {
			continue
		}
		if _, ok := control.Controller.(devicesstructs.MediaController); !ok {
			continue
		}

		topics := resolveMediaTopics(control)
		if topics.Command == "" {
			continue
		}

		r.Subscribe(topics.Command, handleMediaCommand(device.Name, topics, control.PollInterval))
	}
}

func resolveMediaTopics(control *devicesstructs.Control) mediaTopics {
	topics := cfg.GetMqttConfig().Topics[control.MqttId]
	return mediaTopics{
		Command: topics["command_topic"],
		Status:  topics["state_topic"],
	}
}

func handleMediaCommand(deviceName string, topics mediaTopics, pollInterval time.Duration) mqtt.MessageHandler {
	syncStatus := syncMediaStatus(deviceName, topics)
	if pollInterval > 0 {
		startStatusTicker(syncStatus, pollInterval)
	}

	return func(client mqtt.Client, message mqtt.Message) {
		logger.Info("Received media command: ", deviceName, " ", string(message.Payload()))

		var command structs.MediaCommand
		if err := json.Unmarshal(message.Payload(), &command); err != nil {
			logger.Error("Failed to unmarshal media command: ", err)
			return
		}

		if err := media.Execute(deviceName, command); err != nil {
			logger.Error("Failed to execute media command: ", err)
			return
		}

		if pollInterval > 0 {
			syncWithMultiTry(syncStatus)
			return
		}
		syncStatus()
	}
}

func syncMediaStatus(deviceName string, topics mediaTopics) func() {
	return func() {
		status, err := media.GetStatus(deviceName)
		if err != nil {
			logger.Error("Failed to get media status: ", err)
			return
		}

		if topics.Status == "" {
			logger.Error("no status topic: ", deviceName)
			return
		}
		status.Thumbnail = nil
		mqttpkg.Publish(topics.Status, status)
	}
}
