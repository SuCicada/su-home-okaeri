package mqttentry

import (
	"encoding/json"
	"fmt"
	"sucicada/home/internal/app/media"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/consts"
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
	type devicehandler struct {
		device    string
		getHandle func(deviceName string, topics mediaTopics) mqtt.MessageHandler
	}
	devicehandlers := []devicehandler{
		{
			device:    consts.DEVICE_WINDOWS,
			getHandle: handleMediaWindows,
		},
		{
			device:    consts.DEVICE_LINUX,
			getHandle: handleMediaLinux,
		},
	}

	for _, devicehandler := range devicehandlers {
		device := deviceservice.GetDevice(devicehandler.device)
		control := device.DeviceControl.Media
		isOk := false
		if control != nil {
			topics := resolveMediaTopics(control)
			if topics.Command != "" {
				r.Subscribe(topics.Command, devicehandler.getHandle(device.Name, topics))
				isOk = true
			}
		}

		if !isOk {
			panic(fmt.Sprint("no Register media: ", device, control))
		}
	}
}

func resolveMediaTopics(control *devicesstructs.Control) mediaTopics {
	topics := cfg.GetMqttConfig().Topics[control.MqttId]
	return mediaTopics{
		Command: topics["command_topic"],
		Status:  topics["state_topic"],
	}
}

func handleMediaLinux(deviceName string, topics mediaTopics) mqtt.MessageHandler {
	return func(client mqtt.Client, message mqtt.Message) {
		logger.Info("Received linux media command: ", string(message.Payload()))
		var command structs.MediaCommand
		err := json.Unmarshal(message.Payload(), &command)
		if err != nil {
			logger.Error("Failed to unmarshal media command: ", err)
			return
		}

		err = media.Execute(deviceName, command)
		if err != nil {
			logger.Error("Failed to execute media command: ", err)
			return
		}

		status, err := media.GetStatus(deviceName)
		if err != nil {
			logger.Error("Failed to get media status: ", err)
			return
		}

		if topics.Status != "" {
			status.Thumbnail = nil
			mqttpkg.Publish(topics.Status, status)
		} else {
			logger.Error("no status topic: ", deviceName)
		}
	}
}

// ==============================================================

func handleMediaWindows(deviceName string, topics mediaTopics) mqtt.MessageHandler {
	syncStatus := syncMediaStatus(deviceName, topics)
	startStatusTicker(syncStatus, 10*time.Second)

	return func(client mqtt.Client, message mqtt.Message) {
		logger.Info("Received media command: ", string(message.Payload()))

		var command structs.MediaCommand
		err := json.Unmarshal(message.Payload(), &command)
		if err != nil {
			logger.Error("Failed to unmarshal media command: ", err)
			return
		}

		err = media.Execute(deviceName, command)
		if err != nil {
			logger.Error("Failed to execute media command: ", err)
			return
		}

		// 操作后不要只同步一次，做几次延迟校准
		syncStatus()
		time.Sleep(100 * time.Millisecond)
		syncStatus()
		time.Sleep(500 * time.Millisecond)
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

		if topics.Status != "" {
			status.Thumbnail = nil
			mqttpkg.Publish(topics.Status, status)
		}
	}
}
