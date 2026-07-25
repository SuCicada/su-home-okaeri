package linux

import (
	"strings"
	"sucicada/home/internal/cfg"
	devicesstructs "sucicada/home/internal/structs/devices"
	"sucicada/home/internal/util"
)

var Device = devicesstructs.DeviceBase{
	Name: "linux",
	DeviceControl: devicesstructs.DeviceControlUnit{
		Light: &devicesstructs.Control{
			MqttId:     "light_linux",
			Controller: &sLinuxLight{},
		},
		Audio: &devicesstructs.Control{
			Controller: &sLinuxAudio{},
		},
		Media: &devicesstructs.Control{
			MqttId:     "media_linux",
			Controller: &sLinuxMedia{},
		},
	},
}

func sshLinux(cmd string) (string, error) {
	var res, err = util.SSH.SSHRun(cfg.GetSSHConfig(Device.Name), cmd)
	res = strings.TrimSpace(res)
	return res, err
}

func controlOptions(control string) map[string]any {
	options := cfg.GetDeviceConfig(Device.Name).Control[control]["options"]
	return util.Conv.AnyToMap(options)
}
