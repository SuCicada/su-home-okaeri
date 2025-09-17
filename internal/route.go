package internal

import (
	"SuCicada/home/internal/controller"
	"SuCicada/home/internal/service"

	"github.com/gin-gonic/gin"
)

func GetRoute(r *gin.Engine) {
	r.POST("/sleep", service.Sleep)

	// light := r.Group("/light/:type")
	r.POST("/control/:device", service.SetValue)

	
	r.POST("/sms-check/send", controller.SmsCheck.SendVerifyCode)
	r.POST("/sms-check/webhook", controller.SmsCheck.Webhook)
	r.POST("/sms-check/check", controller.SmsCheck.CheckSMS)

}
