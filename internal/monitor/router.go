package monitor

import (
	"sucicada/home/internal/controller"
	monitorcontroller "sucicada/home/internal/monitor/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.POST("/sms-check/send", controller.SmsCheck.SendVerifyCode)
	r.POST("/sms-check/webhook", controller.SmsCheck.Webhook)
	r.POST("/sms-check/check", controller.SmsCheck.CheckSMS)
	r.POST("/alertmanager/webhook", monitorcontroller.Alertmanager.Webhook)

}
