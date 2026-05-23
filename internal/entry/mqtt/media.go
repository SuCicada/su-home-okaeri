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

func RegisterMediaRoutes() {
	for _, device := range deviceservice.GetDevices() {
		registerMediaRoute(device.Name, device.DeviceControl.Media)
	}
}

func registerMediaRoute(deviceName string, control *devicesstructs.Control) {
	if control == nil {
		return
	}
	topics := cfg.GetMqttConfig().Topics[control.MqttId]

	commandTopic := topics["command_topic"]
	statusTopic := topics["status_topic"]
	//isMutedTopic := topics["is_muted_topic"]
	//volumeTopic := topics["volume_topic"]

	if commandTopic == "" {
		return
	}

	var syncStatus = func() {
		status, err := media.GetStatus(deviceName)
		if err != nil {
			logger.Error("Failed to get media status: ", err)
			return
		}

		if statusTopic != "" {
			status.Thumbnail = nil
			mqttpkg.Publish(statusTopic, status)
		}
		//if isMutedTopic != "" {
		//	mqttpkg.Publish(isMutedTopic, status.IsMute)
		//}
		//if volumeTopic != "" {
		//	mqttpkg.Publish(volumeTopic, status.Volume)
		//}
	}

	syncStatus()
	mqttpkg.RegisterRoute(commandTopic, func(client mqtt.Client, message mqtt.Message) {
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
	})

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			syncStatus()
		}
	}()
}
