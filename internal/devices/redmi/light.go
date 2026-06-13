package redmi

import (
	"errors"
	"fmt"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/consts"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/structs"
	commandstructs "sucicada/home/internal/structs/command"
)

type sRedmiLight struct{}

func (l *sRedmiLight) SetBrightness(command commandstructs.LightCommand) error {
	_, err := ssh(fmt.Sprintf("termux-brightness %d", command.Light))
	return err
}

func (l *sRedmiLight) GetBrightness() (int, error) {
	logger.Warn("not support get redmi light")
	return -1, errors.New("not support get redmi light")
}

func (l *sRedmiLight) GetStatus() (structs.LightStatus, error) {
	_, err := l.GetBrightness()
	if err != nil {
		return structs.LightStatus{}, err
	}
	return structs.LightStatus{}, nil
}

func (l *sRedmiLight) On() error {
	return l.SetBrightness(commandstructs.LightCommand{Light: getRedmiLightHigh()})
}

func (l *sRedmiLight) Off() error {
	return l.SetBrightness(commandstructs.LightCommand{Light: 0})
}

func getRedmiLightHigh() int {
	control := cfg.GetDeviceConfig(Device.Name).Control[consts.CONTROL_LIGHT]
	if high, ok := control["high"]; ok {
		switch value := high.(type) {
		case int:
			return value
		case int64:
			return int(value)
		case float64:
			return int(value)
		}
	}
	return 100
}
