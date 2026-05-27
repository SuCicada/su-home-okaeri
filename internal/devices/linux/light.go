package linux

import (
	"encoding/json"
	"fmt"
	"sucicada/home/internal/consts"
	"sucicada/home/internal/logger"
	commandstructs "sucicada/home/internal/structs/command"
	"sucicada/home/internal/util"
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
	options := Config.Control[consts.CONTROL_LIGHT]["options"]
	optionsMap := util.Conv.AnyToMap(options)
	if pactlOpts, ok := optionsMap["ddcutil"]; ok {
		return pactlOpts.(string)
	}
	return optionsMap["ddcutil"].(string)
}
func (l *sLinuxLight) Get() (any, error) {
	return l.GetBrightness()
}

func (l *sLinuxLight) GetBrightness() (int, error) {
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

func (l *sLinuxLight) Set(command string) error {
	var data commandstructs.LightCommand
	if err := json.Unmarshal([]byte(command), &data); err != nil {
		logger.Error("Failed to unmarshal light command: ", err)
		return err
	}
	return l.SetBrightness(data)
}

func (l *sLinuxLight) SetBrightness(command commandstructs.LightCommand) error {
	_, err := sshLinux(fmt.Sprintf(`
	ddcutil %s setvcp 10 %d
	`, getOpts(), command.Light))

	return err
}

func (l *sLinuxLight) On() error {
	return l.setPowerMode("01")
}

func (l *sLinuxLight) Off() error {
	return l.setPowerMode("05")
}

func (l *sLinuxLight) setPowerMode(value string) error {
	_, err := sshLinux(fmt.Sprintf(`
	ddcutil %s setvcp D6 %s
	`, getOpts(), value))
	return err
}
