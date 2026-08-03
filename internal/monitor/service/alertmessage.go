package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/SuCicada/apprise-sdk-go/apprise"
	amwebhook "github.com/prometheus/alertmanager/notify/webhook"
)

func AlertmanagerAppriseMessage(w amwebhook.Message) apprise.Message {
	status := strings.ToUpper(w.Status)
	if status == "" {
		status = "FIRING"
	}

	name := w.CommonLabels["alertname"]
	if name == "" {
		name = "Alertmanager"
	}
	title := fmt.Sprintf("[%s:%d] %s", status, len(w.Alerts), name)

	messageType := apprise.TypeFailure
	if strings.EqualFold(w.Status, "resolved") {
		messageType = apprise.TypeSuccess
	}

	return apprise.Message{
		Title: title,
		Body:  alertmanagerBody(w),
		Tag:   "mail",
		Type:  messageType,
	}
}

func alertmanagerBody(w amwebhook.Message) string {
	parts := make([]string, 0, len(w.Alerts)+1)
	for _, alert := range w.Alerts {
		text := firstNonEmpty(alert.Annotations["summary"], alert.Annotations["description"], formatLabels(alert.Labels))
		if description := alert.Annotations["description"]; description != "" && description != text {
			text += "\n" + description
		}
		if alert.GeneratorURL != "" {
			text += "\n" + alert.GeneratorURL
		}
		parts = append(parts, "- "+text)
	}
	if w.ExternalURL != "" {
		parts = append(parts, "Alertmanager: "+w.ExternalURL)
	}
	return strings.Join(parts, "\n\n")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return "Alertmanager notification"
}

func formatLabels(labels map[string]string) string {
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%s=%s", key, labels[key]))
	}
	return strings.Join(values, ", ")
}
