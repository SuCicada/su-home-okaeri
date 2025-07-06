package internal

import (
	"SuCicada/home/internal/service"

	"github.com/gin-gonic/gin"
)

func GetRoute(r *gin.Engine) {
	r.POST("/sleep", service.Sleep)

	// light := r.Group("/light/:type")
	r.POST("/light/:device", service.SetLight)
}
