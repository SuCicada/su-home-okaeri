/*
media control for windows
openapi: http://windows.sucicada.me:5278/swagger/index.html

*/

package windows

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sucicada/home/internal/cfg"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/structs"
	"sucicada/home/internal/util"

	"resty.dev/v3"
)

type sWindowsMedia struct{}

var WindowsMedia = &sWindowsMedia{}

func (l *sWindowsMedia) Get() (any, error) {
	return l.GetStatus()
}

func (l *sWindowsMedia) GetStatus() (structs.MediaStatus, error) {
	type StatusRes struct {
		Data structs.MediaStatus `json:"data"`
	}
	res, err := sendGetRequest("/api/media/status")
	if err != nil {
		return structs.MediaStatus{}, err
	}
	var statusRes StatusRes
	err = json.Unmarshal([]byte(res), &statusRes)
	return statusRes.Data, err
}

func (l *sWindowsMedia) Set(command string) error {
	var commandRequest structs.MediaCommand
	err := json.Unmarshal([]byte(command), &commandRequest)
	if err != nil {
		logger.Error("Failed to unmarshal command: ", err)
		return err
	}
	return l.Execute(commandRequest)
}

func (l *sWindowsMedia) Execute(commandRequest structs.MediaCommand) error {
	var res string
	var err error
	switch commandRequest.Command {
	case "play":
		res, err = sendRequest("/api/media/play")
	case "pause", "stop":
		res, err = sendRequest("/api/media/pause")
	case "next":
		res, err = sendRequest("/api/media/next")
	case "previous":
		res, err = sendRequest("/api/media/previous")
	case "mute":
		res, err = sendRequest("/api/volume/mute")
	case "volume_set":
		var volume, err = parseVolume(commandRequest.Args["volume"])
		if err != nil {
			return err
		}
		if volume < 0 || volume > 1 {
			return errors.New("volume must be between 0 and 1")
		}
		level := int(volume * 100)
		body := map[string]any{"level": level}
		logger.Info("windows media volume set: ", body)
		res, err = sendRequestWithBody("/api/volume/set", body)
	case "volume_up":
		res, err = sendRequest("/api/volume/up")
	case "volume_down":
		res, err = sendRequest("/api/volume/down")
	default:
		logger.Warn("Unknown media command: ", commandRequest.Command)
		return errors.New("unknown media command")
	}

	if res != "" {
		logger.Info("windows media command response: ", res)
	}
	if err != nil {
		logger.Error("Failed to execute media command: ", err)
	}
	return err
}

func parseVolume(volume any) (float64, error) {
	switch value := volume.(type) {
	case string:
		return strconv.ParseFloat(value, 64)
	case float64:
		return value, nil
	case int:
		return float64(value), nil
	default:
		return 0, fmt.Errorf("volume is not number: %v", volume)
	}
}

type ControllerResponse struct {
	Success bool `json:"success"`
	// Code    int    `json:"code"`
	// Message string `json:"message"`
	// Data    any    `json:"data"`
}

func sendGetRequest(apiPath string) (string, error) {
	return doRequest("GET", apiPath, nil)
}

func sendRequest(apiPath string) (string, error) {
	return doRequest("POST", apiPath, nil)
}

func sendRequestWithBody(apiPath string, body any) (string, error) {
	return doRequest("POST", apiPath, body)
}

func doRequest(method, apiPath string, body any) (string, error) {
	type MediaConfig struct {
		APIBaseURL string `json:"api_base_url"`
	}
	mediaConfMap := cfg.GetDeviceConfig(Device.Name).Control["media"]
	var mediaConfig MediaConfig
	err := util.Conv.MapToAny(mediaConfMap, &mediaConfig)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	logger.Info("windows media APIBaseURL:", mediaConfig.APIBaseURL)
	logger.Info("windows media req:", method, apiPath)

	req := resty.New().
		//SetDebug(true).
		SetBaseURL(mediaConfig.APIBaseURL).
		R()

	if body != nil {
		req.SetBody(body)
	}

	var res *resty.Response
	switch method {
	case "GET":
		res, err = req.Get(apiPath)
	case "POST":
		res, err = req.Post(apiPath)
	default:
		return "", errors.New("unsupported HTTP method: " + method)
	}

	if err != nil {
		logger.Error("send request error: ", err)
		return "", err
	}
	var controllerResponse ControllerResponse
	json.Unmarshal(res.Bytes(), &controllerResponse)
	if res.IsError() || !controllerResponse.Success {
		logger.Error("controller response is not success: ", res.StatusCode(), res.String())
		return "", errors.New("controller response is not success")
	}
	return res.String(), nil
}
