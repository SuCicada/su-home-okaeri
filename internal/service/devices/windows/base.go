package windows

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/service/devices"
	"SuCicada/home/internal/util"
)

func init() {
	devices.RegisterDevice(&Device)
}

var Device = devices.DeviceBase{
	Name: "windows",
	DeviceControl: devices.DeviceControlUnit{
		// Light: &devices.Control{
		// Control: &sRedmiLight{},
		// },
		Volume: &devices.Control{
			Control: &sWindowsVolume{},
		},
	},
}

func ssh(cmd string) (string, error) {
	return util.SSHRun(cfg.GetSSHConfig(Device.Name), cmd)
}
