package internal

import (
	"SuCicada/home/internal/controller"

	"github.com/gin-gonic/gin"
)

func GetRoute(r *gin.Engine) {
	r.POST("/sleep", controller.Control.Sleep)
	r.POST("/control/:device", controller.Control.SetValue)

	r.POST("/sms-check/send", controller.SmsCheck.SendVerifyCode)
	r.POST("/sms-check/webhook", controller.SmsCheck.Webhook)
	r.POST("/sms-check/check", controller.SmsCheck.CheckSMS)

}
