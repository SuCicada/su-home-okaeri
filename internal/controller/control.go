package controller

import (
	"SuCicada/home/internal/logger"
	"SuCicada/home/internal/service/devices"
	"SuCicada/home/internal/util"
	"fmt"
	"strconv"

	"SuCicada/home/internal/service/devices/linux"
	"SuCicada/home/internal/service/devices/redmi"
	"SuCicada/home/internal/service/devices/windows"

	"net/http"

	"github.com/gin-gonic/gin"
)

type cControl struct {
}

var Control = cControl{}

func (c *cControl) Sleep(ginC *gin.Context) {
	linux.Device.DeviceControl.Light.Control.Set(
		linux.Device.DeviceControl.Light.LOW(),
	)
	redmi.Device.DeviceControl.Light.Control.Set(
		redmi.Device.DeviceControl.Light.LOW(),
	)
	windows.Device.DeviceControl.Volume.Control.Set(
		windows.Device.DeviceControl.Volume.LOW(),
	)
}

func (c *cControl) SetValue(ginC *gin.Context) {
	device := ginC.Param("device")

	// req := struct {
	// Light int `json:"light" form:"light"`
	// }{}
	//err := ginC.ShouldBind(&req)
	req, err := util.Conv.GetMapFromGinContext(ginC)
	if err != nil {
		ginC.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// c.Request.PostForm 是 map[string][]string

	deviceBase := devices.GetDevice(device)
	if deviceBase == nil {
		ginC.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device"})
		return
	}

	logger.Info(fmt.Sprintf("SetValue: %s, %v", device, req))

	errs := []error{}
	for k, v := range req {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		err = deviceBase.DeviceControl.GetControl(k).Control.Set(vInt)
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
		ginC.JSON(http.StatusBadRequest, gin.H{"error": errs})
		return
	}

	ginC.JSON(http.StatusOK, gin.H{"ok": true})

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"ok": false, "result": result, "error": err.Error()})
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{"ok": true, "result": result})
	// }
}
