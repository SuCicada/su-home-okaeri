/*
light control for mac
openapi: http://mac.sucicada.me:8765
tag: Brightness

GET  /status
POST /on
POST /off
*/

package mac

import (
	"encoding/json"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/structs"
	commandstructs "sucicada/home/internal/structs/command"
	"sucicada/home/internal/util"
)

type sMacLight struct{}

type statusAPIResponse struct {
	DisplayAsleep bool `json:"display_asleep"`
	DisplayOn     bool `json:"display_on"`
}

const fixedBrightness = 100

func (l *sMacLight) GetBrightness() (int, error) {
	return fixedBrightness, nil
}

func (l *sMacLight) GetStatus() (structs.LightStatus, error) {
	res, err := doRequest("GET", "/status", nil)
	if err != nil {
		return structs.LightStatus{}, err
	}

	var resp statusAPIResponse
	if err := json.Unmarshal([]byte(res), &resp); err != nil {
		logger.Error("Failed to unmarshal status: ", err)
		return structs.LightStatus{}, err
	}

	return structs.LightStatus{
		Power:      resp.DisplayOn && !resp.DisplayAsleep,
		Brightness: fixedBrightness,
	}, nil
}

func (l *sMacLight) SetBrightness(command commandstructs.LightCommand) error {
	if command.Light > 0 {
		return l.On()
	}
	return l.Off()
}

func (l *sMacLight) On() error {
	logger.Info("Turning on light")
	_, err := doRequest("POST", "/on", nil)
	return err
}

func (l *sMacLight) Off() error {
	logger.Info("Turning off light")
	_, err := doRequest("POST", "/off", nil)
	return err
}

func doRequest(method, apiPath string, body any) (string, error) {
	type LightConfig struct {
		APIBaseURL string `json:"api_base_url"`
	}
	lightConfMap := cfg.GetDeviceConfig(Device.Name).Control["light"]
	var lightConfig LightConfig
	err := util.Conv.MapToAny(lightConfMap, &lightConfig)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	res, err := util.Request.Do(method, lightConfig.APIBaseURL+apiPath, body)
	if err != nil {
		logger.Error("send request error: ", err)
		return "", err
	}
	return res, nil
}
