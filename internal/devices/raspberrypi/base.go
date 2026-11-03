package raspberrypi

import (
	"strings"
	"sucicada/home/internal/cfg"
	devicesstructs "sucicada/home/internal/structs/devices"
	"sucicada/home/internal/util"
	"sucicada/home/internal/util/tinyssh"
)

var Device = devicesstructs.DeviceBase{
	Name: "raspberrypi",
	DeviceControl: devicesstructs.DeviceControlUnit{
		Light: &devicesstructs.Control{
			MqttId:     "light_raspberrypi",
			Controller: &sRaspberryPiLight{},
		},
	},
}

func ssh(cmd string) (string, error) {
	res, err := tinyssh.SSH.SSHRun(cfg.GetSSHConfig(Device.Name), cmd)
	return strings.TrimSpace(res), err
}

func controlOptions(control string) map[string]any {
	controlCfg := cfg.GetDeviceConfig(Device.Name).Control[control]
	if controlCfg == nil {
		return map[string]any{}
	}
	return util.Conv.AnyToMap(controlCfg["options"])
}
