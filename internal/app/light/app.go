package light

import (
	deviceservice "sucicada/home/internal/service/devices"
	"sucicada/home/internal/structs"
	commandstructs "sucicada/home/internal/structs/command"
)

func GetStatus(deviceName string) (structs.LightStatus, error) {
	controller, err := deviceservice.GetLightController(deviceName)
	if err != nil {
		return structs.LightStatus{}, err
	}
	return controller.GetStatus()
}

func SetBrightness(deviceName string, command commandstructs.LightCommand) error {
	controller, err := deviceservice.GetLightController(deviceName)
	if err != nil {
		return err
	}
	return controller.SetBrightness(command)
}

func On(deviceName string) error {
	controller, err := deviceservice.GetLightController(deviceName)
	if err != nil {
		return err
	}
	return controller.On()
}

func Off(deviceName string) error {
	controller, err := deviceservice.GetLightController(deviceName)
	if err != nil {
		return err
	}
	return controller.Off()
}
