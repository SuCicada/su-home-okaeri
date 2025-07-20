package service

import (
	"SuCicada/home/internal/service/devices"
	"SuCicada/home/internal/service/devices/linux"
	"SuCicada/home/internal/service/devices/redmi"
	"SuCicada/home/internal/service/devices/windows"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Sleep(c *gin.Context) {
	linux.Device.DeviceControl.Light.Control.Set(
		linux.Device.DeviceControl.Light.LOW(),
	)
	redmi.Device.DeviceControl.Light.Control.Set(
		redmi.Device.DeviceControl.Light.LOW(),
	)
	windows.Device.DeviceControl.Volume.Control.Set(
		windows.Device.DeviceControl.Volume.LOW(),
	)
	// redmi.Redmi.Device.Set(10)
}

func SetValue(c *gin.Context) {
	device := c.Param("device")

	// req := struct {
	// Light int `json:"light" form:"light"`
	// }{}
	req := map[string]int{}
	c.ShouldBind(&req)
	// controlConfig := devices.GetDeviceControlConfig(device)
	deviceBase := devices.GetDevice(device)
	if deviceBase == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device"})
		return
	}

	errs := []error{}
	for k, v := range req {
		err := deviceBase.DeviceControl.GetControl(k).Control.Set(v)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// control.Control.Set(req[device])
	// var lightDevice light.DeviceBase

	// switch device {
	// case "linux":
	// 	lightDevice = light.LinuxLight

	// case "redmi":
	// 	lightDevice = light.RedmiLight

	// default:
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device"})
	// 	return
	// }

	// var result string
	// var err error
	// if req.Light > 0 {
	// 	result, err = lightDevice.Set(req.Light)
	// } else {
	// 	result, err = lightDevice.Toggle()
	// }

	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"ok": false, "result": result, "error": err.Error()})
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{"ok": true, "result": result})
	// }
}
