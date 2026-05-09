package linux

import (
	"encoding/json"
	"errors"
	"fmt"
	"sucicada/home/internal/consts"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/util"
)

type sLinuxMedia struct{}

func getPactlOpts() int {
	options := Config.Control[consts.CONTROL_MEDIA]["options"].(map[string]any)
	pactlOpts := options["pactl"]
	return pactlOpts.(int)
}

func (l *sLinuxMedia) Get() (int, error) {
	pactlOpts := getPactlOpts()
	res, err := sshLinux(fmt.Sprintf("pactl get-sink-volume %d", pactlOpts))
	if err != nil {
		return 0, err
	}
	return util.Conv.StrToInt(res), nil
}

func (m *sLinuxMedia) Set(command string) error {
	type CommandRequest struct {
		Command string         `json:"command"`
		Args    map[string]any `json:"args"`
	}

	var commandRequest CommandRequest
	err := json.Unmarshal([]byte(command), &commandRequest)
	if err != nil {
		logger.Error("Failed to unmarshal command: ", err)
		return err
	}

	pactlOpts := getPactlOpts()

	switch commandRequest.Command {
	case "volume_set":
		if volume, ok := commandRequest.Args["volume"]; ok {
			_, err = sshLinux(fmt.Sprintf("pactl set-sink-volume %d %v%%", pactlOpts, volume))
		}
	case "volume_up":
		step := 5
		if s, ok := commandRequest.Args["step"]; ok {
			step = int(s.(float64))
		}
		_, err = sshLinux(fmt.Sprintf("pactl set-sink-volume %d +%d%%", pactlOpts, step))
	case "volume_down":
		step := 5
		if s, ok := commandRequest.Args["step"]; ok {
			step = int(s.(float64))
		}
		_, err = sshLinux(fmt.Sprintf("pactl set-sink-volume %d -%d%%", pactlOpts, step))
	case "mute":
		_, err = sshLinux(fmt.Sprintf("pactl set-sink-mute %d toggle", pactlOpts))
	default:
		logger.Warn("Unknown media command: ", commandRequest.Command)
		return errors.New("unknown media command")
	}

	return err
}
