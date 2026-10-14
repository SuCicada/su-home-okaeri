package devices

import (
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/devices/linux"
	"sucicada/home/internal/devices/mac"
	"sucicada/home/internal/devices/redmi"
	"sucicada/home/internal/devices/windows"
	"sucicada/home/internal/logger"
	deviceservice "sucicada/home/internal/service/devices"
	devicesstructs "sucicada/home/internal/structs/devices"
)

var Devices = []devicesstructs.DeviceBase{
	linux.Device,
	redmi.Device,
	windows.Device,
	mac.Device,
}

func Init() {
	enabled := cfg.GetConfig().Devices
	for _, device := range Devices {
		if _, ok := enabled[device.Name]; !ok {
			logger.Info("Skip device not in config: ", device.Name)
			continue
		}
		deviceservice.RegisterDevice(&device)
	}
}
