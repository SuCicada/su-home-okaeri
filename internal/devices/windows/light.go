/*
light / brightness control for windows
openapi: http://windows.sucicada.me:5278/swagger/index.html
tag: Brightness

GET  /api/brightness
GET  /api/brightness/screen-state
POST /api/brightness/set   body: { level }
POST /api/brightness/on
POST /api/brightness/off
*/

package windows

import (
	"encoding/json"
	"strings"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/structs"
	commandstructs "sucicada/home/internal/structs/command"
)

type sWindowsLight struct{}

type int32APIResponse struct {
	Success bool `json:"success"`
	Data    int  `json:"data"`
}

type stringAPIResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
}

func (l *sWindowsLight) GetBrightness() (int, error) {
	res, err := sendGetRequest("/api/brightness")
	if err != nil {
		return 0, err
	}
	var resp int32APIResponse
	if err := json.Unmarshal([]byte(res), &resp); err != nil {
		logger.Error("Failed to unmarshal brightness: ", err)
		return 0, err
	}
	return resp.Data, nil
}

func (l *sWindowsLight) GetStatus() (structs.LightStatus, error) {
	brightness, err := l.GetBrightness()
	if err != nil {
		return structs.LightStatus{}, err
	}

	res, err := sendGetRequest("/api/brightness/screen-state")
	if err != nil {
		return structs.LightStatus{
			Power:      brightness > 0,
			Brightness: brightness,
		}, nil
	}

	var resp stringAPIResponse
	if err := json.Unmarshal([]byte(res), &resp); err != nil {
		logger.Error("Failed to unmarshal screen-state: ", err)
		return structs.LightStatus{
			Power:      brightness > 0,
			Brightness: brightness,
		}, nil
	}

	power := strings.EqualFold(strings.TrimSpace(resp.Data), "On")
	return structs.LightStatus{
		Power:      power,
		Brightness: brightness,
	}, nil
}

func (l *sWindowsLight) SetBrightness(command commandstructs.LightCommand) error {
	body := map[string]any{"level": command.Light}
	_, err := sendRequestWithBody("/api/brightness/set", body)
	return err
}

func (l *sWindowsLight) On() error {
	_, err := sendRequest("/api/brightness/on")
	return err
}

func (l *sWindowsLight) Off() error {
	_, err := sendRequest("/api/brightness/off")
	return err
}
