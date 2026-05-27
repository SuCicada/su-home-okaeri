package internal

import (
	"sucicada/home/internal/controller"
	httpentry "sucicada/home/internal/entry/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRoute(r *gin.Engine) {
	// r.POST("/sleep", controller.Control.Sleep)
	// r.POST("/control/:device", controller.Control.SetValue)

	r.POST("/sms-check/send", controller.SmsCheck.SendVerifyCode)
	r.POST("/sms-check/webhook", controller.SmsCheck.Webhook)
	r.POST("/sms-check/check", controller.SmsCheck.CheckSMS)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	httpentry.RegisterAudioRoutes(r)
	httpentry.RegisterMediaRoutes(r)
}
