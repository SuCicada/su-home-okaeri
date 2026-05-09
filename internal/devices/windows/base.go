package windows

import (
	"sucicada/home/internal/service/devices"
	devicesstructs "sucicada/home/internal/structs/devices"
)

func init() {
	devices.RegisterDevice(&Device)
}

var Device = devicesstructs.DeviceBase{
	Name: "windows",
	DeviceControl: devicesstructs.DeviceControlUnit{
		Media: &devicesstructs.Control{
			MqttId:  "windows_media",
			Control: WindowsMedia,
		},
	},
}

//func ssh(cmd string) (string, error) {
//	return util.SSHRun(cfg.GetSSHConfig(Device.Name), cmd)
//}
