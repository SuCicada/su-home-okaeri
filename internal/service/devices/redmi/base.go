package redmi

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/service/devices"
	"SuCicada/home/internal/util"
)

func init() {
	devices.RegisterDevice(&Device)
}

var Device = devices.DeviceBase{
	Name: "redmi",
	DeviceControl: devices.DeviceControlUnit{
		Light: &devices.Control{
			Control: &sRedmiLight{},
		},
		// Volume: service.Control{
		// Control: &sLinuxVolume{},
		// },
	},
}

func ssh(cmd string) (string, error) {
	return util.SSHRunRoot(cfg.GetSSHConfig(Device.Name), cmd)
}
