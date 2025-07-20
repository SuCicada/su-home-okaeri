package linux

import (
	"SuCicada/home/internal/logger"
	"SuCicada/home/internal/util"
	"fmt"
)

// var LinuxLight = devices.DeviceBase{
// 	Name: "linux",
// 	DeviceControl: devices.DeviceControlUnit{
// 		Light: &devices.Control{
// 			Control: &sLinuxLight{},
// 		},
// 	},
// }

type sLinuxLight struct{}

// func (l *sLinuxLight) sshLinux(cmd string) (string, error) {
// 	return util.SSHRunRoot(util.GetSSHConfig("linux"), cmd)
// }

func (l *sLinuxLight) Get() (int, error) {
	res, err := sshLinux(`
		ddcutil --bus=5 getvcp 10 | grep -i "current value" | awk '{print $9}' | tr -d ','
	`)
	if err != nil {
		logger.Error("Error getting light:", err)
		return 0, err
	}
	return util.StrToInt(res), nil
}

func (l *sLinuxLight) Set(light int) error {
	_, err := sshLinux(fmt.Sprintf("ddcutil --bus=5 setvcp 10 %d", light))
	return err
}
