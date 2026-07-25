package windows

import (
	devicesstructs "sucicada/home/internal/structs/devices"
)

var Device = devicesstructs.DeviceBase{
	Name: "windows",
	DeviceControl: devicesstructs.DeviceControlUnit{
		Light: &devicesstructs.Control{
			MqttId:     "light_windows",
			Controller: &sWindowsLight{},
		},
		Media: &devicesstructs.Control{
			MqttId:     "media_windows",
			Controller: WindowsMedia,
		},
	},
}
