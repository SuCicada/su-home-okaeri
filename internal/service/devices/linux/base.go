package linux

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/service/devices"
	"SuCicada/home/internal/util"
)

func init() {
	devices.RegisterDevice(&Device)
}

var Device = devices.DeviceBase{
	Name: "linux",
	DeviceControl: devices.DeviceControlUnit{
		Light: &devices.Control{
			Control: &sLinuxLight{},
		},
		// Volume: service.Control{
		// Control: &sLinuxVolume{},
		// },
	},
}

func sshLinux(cmd string) (string, error) {
	return util.SSHRunRoot(cfg.GetSSHConfig("linux"), cmd)
}

// func (l *sLinuxLight) Get() (int, error) {
// 	res, err := l.sshLinux(`
// 		ddcutil --bus=5 getvcp 10 | grep -i "current value" | awk '{print $9}' | tr -d ','
// 	`)
// 	if err != nil {
// 		logger.Error("Error getting light:", err)
// 		return 0, err
// 	}
// 	return util.StrToInt(res), nil
// }

// // func (l *sLinuxLight) Toggle() (string, error) {
// // 	light, err := l.Get()
// // 	if err != nil {
// // 		return "", err
// // 	}
// // 	if light < HIGH_LIGHT {
// // 		return l.Set(HIGH_LIGHT)
// // 	} else {
// // 		return l.Set(LOW_LIGHT)
// // 	}
// // }

// func (l *sLinuxLight) Set(light int) (string, error) {
// 	return l.sshLinux(fmt.Sprintf("ddcutil --bus=5 setvcp 10 %d", light))
// }
