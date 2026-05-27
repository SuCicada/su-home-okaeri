package light

import (
	deviceservice "sucicada/home/internal/service/devices"
	commandstructs "sucicada/home/internal/structs/command"
)

func GetBrightness(deviceName string) (int, error) {
	controller, err := deviceservice.GetLightController(deviceName)
	if err != nil {
		return 0, err
	}
	return controller.GetBrightness()
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
