package linux

import (
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/service/devices"
	devicesstructs "sucicada/home/internal/structs/devices"
	"sucicada/home/internal/util"
)

func init() {
	devices.RegisterDevice(&Device)
}

var Device = devicesstructs.DeviceBase{
	Name: "linux",
	DeviceControl: devicesstructs.DeviceControlUnit{
		Light: &devicesstructs.Control{
			MqttId:  "linux_light",
			Control: &sLinuxLight{},
		},
		// Media: &devices.Control{
		// 	MqttId:  "linux_media",
		// 	Control: &sLinuxMedia{},
		// },
	},
}

var Config = cfg.GetDeviceConfig(Device.Name)

func sshLinux(cmd string) (string, error) {
	return util.SSHRunRoot(cfg.GetSSHConfig(Device.Name), cmd)
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
