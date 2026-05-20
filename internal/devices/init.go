package devices

import (
	"sucicada/home/internal/devices/linux"
	"sucicada/home/internal/devices/redmi"
	"sucicada/home/internal/devices/windows"
	deviceservice "sucicada/home/internal/service/devices"
	devicesstructs "sucicada/home/internal/structs/devices"
)

var Devices = []devicesstructs.DeviceBase{
	linux.Device,
	redmi.Device,
	windows.Device,
}

func Init() {
	for _, device := range Devices {
		deviceservice.RegisterDevice(&device)
	}
}
