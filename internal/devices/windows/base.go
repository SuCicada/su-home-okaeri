package windows

import (
	devicesstructs "sucicada/home/internal/structs/devices"
)

// func init() {
// 	devices.RegisterDevice(&Device)
// }

var Device = devicesstructs.DeviceBase{
	Name: "windows",
	DeviceControl: devicesstructs.DeviceControlUnit{
		Media: &devicesstructs.Control{
			MqttId:  "media_windows",
			Control: WindowsMedia,
		},
	},
}

//func ssh(cmd string) (string, error) {
//	return util.SSHRun(cfg.GetSSHConfig(Device.Name), cmd)
//}
