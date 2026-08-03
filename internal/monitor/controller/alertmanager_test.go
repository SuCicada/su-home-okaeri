package monitorcontroller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/structs/appconfig"
	"sucicada/home/internal/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
	amwebhook "github.com/prometheus/alertmanager/notify/webhook"
	amtemplate "github.com/prometheus/alertmanager/template"
	"github.com/stretchr/testify/assert"
)

func mockConfig() {
	cfg.CONFIG_PATH = "/Users/peng/PROGRAM/GitHub/su-home-okaeri/config.yaml"
	originConfig := cfg.GetConfig()

	file := "/tmp/config.yaml"
	cfg.CONFIG_PATH = file

	config := &appconfig.AppConfig{
		Alert: originConfig.Alert,
		Monitor: appconfig.MonitorConfig{
			HomeAssistant: appconfig.HomeAssistantConfig{
				RestartScript: "echo 'restart'",
			},
		},
	}
	yamlStr, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(file, yamlStr, 0644)
	if err != nil {
		panic(err)
	}
}
func TestMockconfig(t *testing.T) {
	mockConfig()

}
func TestHomeAssistantDownRestartsBeforeSendingMessage(t *testing.T) {
	mockConfig()
	gin.SetMode(gin.TestMode)
	// steps := make([]string, 0, 2)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	// router.POST("/alertmanager/webhook", alertmanagerWebhookHandler(func(message apprise.Message) error {
	// 	steps = append(steps, "send")
	// 	received = message
	// 	return nil
	// }, func() error {
	// 	steps = append(steps, "restart")
	// 	return nil
	// }))
	cfg.GetConfig().Monitor.HomeAssistant.RestartScript = "echo 'restart'"

	request := amwebhook.Message{
		Version: "4",
		Data: &amtemplate.Data{
			Status: "firing",
			CommonLabels: map[string]string{
				"alertname": "HomeAssistantDown",
			},
			Alerts: []amtemplate.Alert{
				{Status: "firing", Labels: map[string]string{"alertname": "HomeAssistantDown"}},
			},
		},

		// Status:  "firing",
		// CommonLabels: map[string]string{
		// "alertname": "HomeAssistantDown",
		// },
		// Alerts: []amwebhook.Alert{{Status: "firing", Labels: map[string]string{"alertname": "HomeAssistantDown"}}},
	}
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/alertmanager/webhook",
		strings.NewReader(util.Conv.ToJsonStr(request)))
	ginContext.Request.Header.Set("Content-Type", "application/json")

	Alertmanager.Webhook(ginContext)

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}
