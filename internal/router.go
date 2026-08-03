package internal

import (
	httpentry "sucicada/home/internal/entry/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRoute(r *gin.Engine) {
	// r.POST("/sleep", controller.Control.Sleep)
	// r.POST("/control/:device", controller.Control.SetValue)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	httpentry.RegisterAudioRoutes(r)
	httpentry.RegisterMediaRoutes(r)
}
