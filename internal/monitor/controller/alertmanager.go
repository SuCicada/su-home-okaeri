package monitorcontroller

import (
	"sucicada/home/internal/monitor/service"
	"sucicada/home/internal/response"
	"sucicada/home/internal/util"

	"github.com/SuCicada/apprise-sdk-go/apprise"
	"github.com/gin-gonic/gin"
	amwebhook "github.com/prometheus/alertmanager/notify/webhook"
)

type cAlertmanager struct{}

var Alertmanager = cAlertmanager{}

func (c cAlertmanager) Webhook(ctx *gin.Context) {
	send := util.Alert.SendApprise

	var webhook amwebhook.Message
	if err := ctx.ShouldBindJSON(&webhook); err != nil {
		response.Bad(ctx, err.Error())
		return
	}
	if webhook.Data == nil || len(webhook.Alerts) == 0 {
		response.Bad(ctx, "alerts must not be empty")
		return
	}

	message := handlerWebhook(webhook)
	if err := send(message); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}

/*
/Users/peng/PROGRAM/GitHub/SuConfig/linux/asus/etc/docker/monitor/prometheus/alert.rules
*/
func handlerWebhook(webhook amwebhook.Message) apprise.Message {
	message := service.AlertmanagerAppriseMessage(webhook)
	if service.IsHomeAssistantDown(webhook) {
		message = service.HomeAssistantMessage(webhook)

	} else if service.IsClaudeVMDown(webhook) {
		message = service.ClaudeVMMessage(webhook)
	}

	return message
}
