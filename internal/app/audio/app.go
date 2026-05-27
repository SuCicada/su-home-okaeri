package audio

import (
	deviceservice "sucicada/home/internal/service/devices"
	"sucicada/home/internal/structs"
)

func PlayAudio(deviceName string, command structs.AudioPlayRequest) error {
	controller, err := deviceservice.GetAudioController(deviceName)
	if err != nil {
		return err
	}
	return controller.PlayAudio(command)
}
