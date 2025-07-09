package service

import (
	"SuCicada/home/internal/service/light"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Sleep(c *gin.Context) {
	light.LinuxLight.Set(10)
	light.RedmiLight.Set(20)
}

func SetLight(c *gin.Context) {
	device := c.Param("device")

	req := struct {
		Light int `json:"light" form:"light"`
	}{}
	c.ShouldBind(&req)

	var lightDevice light.LightDevice

	switch device {
	case "linux":
		lightDevice = light.LinuxLight

	case "redmi":
		lightDevice = light.RedmiLight

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device"})
		return
	}

	var result string
	var err error
	if req.Light > 0 {
		result, err = lightDevice.Set(req.Light)
	} else {
		result, err = lightDevice.Toggle()
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "result": result, "error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"ok": true, "result": result})
	}
}
