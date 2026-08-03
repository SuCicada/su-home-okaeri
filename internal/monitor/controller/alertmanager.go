package monitorcontroller

import (
	"sucicada/home/internal/monitor/service"
	"sucicada/home/internal/response"
	"sucicada/home/internal/util"

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

	message := service.AlertmanagerAppriseMessage(webhook)
	if service.IsHomeAssistantDown(webhook) {
		message = service.HomeAssistantMessage(webhook)
	}

	if err := send(message); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx)
}
