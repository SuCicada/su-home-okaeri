package redmi

import (
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/service/devices"
	devicesstructs "sucicada/home/internal/structs/devices"
	"sucicada/home/internal/util/tinyssh"
)

func init() {
	devices.RegisterDevice(&Device)
}

var Device = devicesstructs.DeviceBase{
	Name: "redmi",
	DeviceControl: devicesstructs.DeviceControlUnit{
		Light: &devicesstructs.Control{
			Controller: &sRedmiLight{},
		},
		// Volume: service.Control{
		// Control: &sLinuxVolume{},
		// },
	},
}

func ssh(cmd string) (string, error) {
	return tinyssh.SSH.SSHRunRoot(cfg.GetSSHConfig(Device.Name), cmd)
}
