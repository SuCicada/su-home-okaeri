package windows

import (
	devicesstructs "sucicada/home/internal/structs/devices"
)

var Device = devicesstructs.DeviceBase{
	Name: "windows",
	DeviceControl: devicesstructs.DeviceControlUnit{
		Media: &devicesstructs.Control{
			MqttId:     "media_windows",
			Controller: WindowsMedia,
		},
	},
}
