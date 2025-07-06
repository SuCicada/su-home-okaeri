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

	if req.Light > 0 {
		lightDevice.Set(req.Light)
	} else {
		lightDevice.Toggle()
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
