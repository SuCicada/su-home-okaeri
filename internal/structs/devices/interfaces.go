package devicesstructs

import (
	"sucicada/home/internal/structs"
	commandstructs "sucicada/home/internal/structs/command"
)

type DeviceBase struct {
	Name          string            `yaml:"name"`
	DeviceControl DeviceControlUnit `yaml:"control"`
}

type DeviceControlUnit struct {
	Light *Control `yaml:"light"`
	Media *Control `yaml:"media"`
	Audio *Control `yaml:"audio"`
}

type Control struct {
	Name       string
	Device     *DeviceBase
	MqttId     string
	Controller any
}

type MediaController interface {
	GetStatus() (structs.MediaStatus, error)
	Execute(command structs.MediaCommand) error
}

type LightController interface {
	GetBrightness() (int, error)
	SetBrightness(command commandstructs.LightCommand) error
	GetStatus() (structs.LightStatus, error)
	On() error
	Off() error
}

type AudioController interface {
	PlayAudio(command structs.AudioPlayRequest) error
}
