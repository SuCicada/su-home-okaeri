package service

import (
	"strings"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/devices/linux"

	amwebhook "github.com/prometheus/alertmanager/notify/webhook"

	"github.com/SuCicada/apprise-sdk-go/apprise"
)

type HomeAssistantService struct{}

const APPRISE_TAG = "mail"

// func NewHomeAssistantService(send appriseSender, restart restartAction) *HomeAssistantService {
// return &HomeAssistantService{send: send, restart: restart}
// }
func RestartHomeAssistant() error {
	homeAssistantRestartScript := cfg.GetConfig().Monitor.HomeAssistant.RestartScript
	_, err := linux.SSHLinux(homeAssistantRestartScript)
	return err
}

func IsHomeAssistantDown(w amwebhook.Message) bool {
	alertNames := []string{w.CommonLabels["alertname"]}
	for _, alert := range w.Alerts {
		alertNames = append(alertNames, alert.Labels["alertname"])
	}

	for _, name := range alertNames {
		normalized := strings.NewReplacer("_", "", "-", "", " ", "").Replace(strings.ToLower(name))
		if normalized == "homeassistantdown" || normalized == "homeassistantoffline" {
			return true
		}
	}
	return false
}

func HomeAssistantMessage(w amwebhook.Message) apprise.Message {
	if strings.EqualFold(w.Status, "resolved") {
		return apprise.Message{
			Title: "HomeAssistant 已恢复",
			Body:  "HomeAssistant 已恢复在线。\n本次没有调用重启脚本。",
			Tag:   APPRISE_TAG,
			Type:  apprise.TypeSuccess,
		}
	}

	restartResult := "成功"
	if err := RestartHomeAssistant(); err != nil {
		restartResult = "失败：" + err.Error()
	}
	return apprise.Message{
		Title: "HomeAssistant 掉线",
		Body:  "检测到 HomeAssistant 掉线。\n已经调用重启脚本。\n脚本调用结果：" + restartResult,
		Tag:   APPRISE_TAG,
		Type:  apprise.TypeFailure,
	}
}
