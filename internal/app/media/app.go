package media

import (
	deviceservice "sucicada/home/internal/service/devices"
	"sucicada/home/internal/structs"
)

func GetStatus(deviceName string) (structs.MediaStatus, error) {
	controller, err := deviceservice.GetMediaController(deviceName)
	if err != nil {
		return structs.MediaStatus{}, err
	}
	return controller.GetStatus()
}

func Execute(deviceName string, command structs.MediaCommand) error {
	controller, err := deviceservice.GetMediaController(deviceName)
	if err != nil {
		return err
	}
	return controller.Execute(command)
}
