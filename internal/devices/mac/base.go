package mac

import (
	devicesstructs "sucicada/home/internal/structs/devices"
)

var Device = devicesstructs.DeviceBase{
	Name: "mac",
	DeviceControl: devicesstructs.DeviceControlUnit{
		Light: &devicesstructs.Control{
			MqttId:     "light_mac",
			Controller: &sMacLight{},
		},
	},
}
