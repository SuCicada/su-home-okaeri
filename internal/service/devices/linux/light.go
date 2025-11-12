package linux

import (
	"SuCicada/home/internal/consts"
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

//	func (l *sLinuxLight) sshLinux(cmd string) (string, error) {
//		return util.SSHRunRoot(util.GetSSHConfig("linux"), cmd)
//	}
//
//	func init() {
//		mqttcontroller.RegisterRoute(Device.Name, sLinuxLight.Set)
//	}
func getOpts() string {
	options := Config.Control[consts.CONTROL_LIGHT].Options
	if pactlOpts, ok := options["ddcutil"]; ok {
		return pactlOpts.(string)
	}
	return ""
}
func (l *sLinuxLight) Get() (int, error) {
	opts := getOpts()
	res, err := sshLinux(fmt.Sprintf(`
		ddcutil %s getvcp 10 | grep -i "current value" | awk '{print $9}' | tr -d ','
	`, opts))

	if err != nil {
		logger.Error("Error getting light:", err)
		return 0, err
	}
	return util.Conv.StrToInt(res), nil
}

func (l *sLinuxLight) Set(light int) error {
	_, err := sshLinux(fmt.Sprintf(`
	ddcutil %s setvcp 10 %d
	`, getOpts(), light))

	return err
}
