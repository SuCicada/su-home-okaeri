package service

import (
	"strings"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/devices/linux"

	"github.com/SuCicada/apprise-sdk-go/apprise"
	amwebhook "github.com/prometheus/alertmanager/notify/webhook"
)

var claudeVMServiceName = "claude-vm"
var claudeVMAlertName = "ClaudeVMDown"

func RestartClaudeVM() error {
	restartScript := cfg.GetConfig().Monitor.ClaudeVM.RestartScript
	_, err := linux.SSHLinux(restartScript)
	return err
}

func IsClaudeVMDown(w amwebhook.Message) bool {
	alertNames := []string{w.CommonLabels["alertname"]}
	for _, alert := range w.Alerts {
		alertNames = append(alertNames, alert.Labels["alertname"])
		if alert.Labels["service"] == claudeVMServiceName {
			return true
		}
	}
	if w.CommonLabels["service"] == claudeVMServiceName {
		return true
	}

	for _, name := range alertNames {
		if name == claudeVMAlertName {
			return true
		}
	}
	return false
}

func ClaudeVMMessage(w amwebhook.Message) apprise.Message {
	if strings.EqualFold(w.Status, "resolved") {
		return apprise.Message{
			Title: "Claude VM 已恢复",
			Body:  "claude.sucicada.me:22 已恢复可达。\n本次没有调用重启脚本。",
			Tag:   APPRISE_TAG,
			Type:  apprise.TypeSuccess,
		}
	}

	restartResult := "成功"
	if err := RestartClaudeVM(); err != nil {
		restartResult = "失败：" + err.Error()
	}
	return apprise.Message{
		Title: "Claude VM 不可达",
		Body:  "检测到 claude.sucicada.me:22 TCP 探测失败。\n已经调用重启脚本。\n脚本调用结果：" + restartResult,
		Tag:   APPRISE_TAG,
		Type:  apprise.TypeFailure,
	}
}
