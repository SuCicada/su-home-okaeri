package linux

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sucicada/home/internal/consts"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/structs"
	"sucicada/home/internal/util"

	"github.com/samber/lo"
)

type sLinuxMedia struct{}

func (l *sLinuxMedia) Get() (any, error) {
	return l.GetStatus()
}

func (l *sLinuxMedia) Set(command string) error {
	var data structs.MediaCommand
	if err := json.Unmarshal([]byte(command), &data); err != nil {
		logger.Error("Failed to unmarshal media command: ", err)
		return err
	}
	return l.Execute(data)
}

//	return {
//	 "volume": 50,
//	 "is_mute": false
//	}
func (l *sLinuxMedia) GetStatus() (structs.MediaStatus, error) {
	currentSink, err := sshLinux(`pactl get-default-sink`)
	if err != nil {
		return structs.MediaStatus{}, err
	}
	//res, err := sshLinux(strings.ReplaceAll(fmt.Sprintf(`
	//pactl --format=json list sinks
	//| jq ".[] | select(.name == \"$%s\") | {volume, mute}"`, currentSink),
	//		"\n", ""))
	res, err := sshLinux(fmt.Sprintf(`pactl --format=json list sinks`))
	if err != nil {
		return structs.MediaStatus{}, err
	}
	type Sink struct {
		Name   string `json:"name"`
		Mute   bool   `json:"mute"`
		Volume map[string]struct {
			ValuePercent string `json:"value_percent"`
		} `json:"volume"`
	}
	sinks := []Sink{}
	err = json.Unmarshal([]byte(res), &sinks)
	if err != nil {
		return structs.MediaStatus{}, err
	}
	sink, ok := lo.Find(sinks, func(s Sink) bool {
		return s.Name == currentSink
	})
	if !ok {
		return structs.MediaStatus{}, errors.New("sink not found")
	}
	// 80%
	volumeStr := sink.Volume["front-left"].ValuePercent
	volume, err := strconv.Atoi(strings.TrimSuffix(volumeStr, "%"))
	if err != nil {
		return structs.MediaStatus{}, err
	}
	return structs.MediaStatus{
		Volume: volume,
		IsMute: sink.Mute,
	}, nil
}

func (l *sLinuxMedia) Execute(command structs.MediaCommand) error {
	switch command.Command {
	// case "play":
	// 	_, err := sshLinux("playerctl play")
	// 	return err
	// case "pause":
	// 	_, err := sshLinux("playerctl pause")
	// 	return err
	case "stop":
		_, err := sshLinux("pkill ffplay")
		return err

	case "mute":
		if command.Args["mute"] == "true" {
			_, err := sshLinux(fmt.Sprintf("pactl set-sink-mute %s true", getPactlSink()))
			return err
		} else {
			_, err := sshLinux(fmt.Sprintf("pactl set-sink-mute %s false", getPactlSink()))
			return err
		}

	case "volume_set":
		volume, err := parseLinuxVolume(command.Args["volume"])
		if err != nil {
			return err
		}
		_, err = sshLinux(fmt.Sprintf(
			"pactl set-sink-volume %s %d%%",
			getPactlSink(),
			volume,
		))
		return err
	default:
		return errors.New("unsupported linux media command: " + command.Command)
	}
}

func getPactlSink() string {
	options := Config.Control[consts.CONTROL_MEDIA]["options"]
	optionsMap := util.Conv.AnyToMap(options)
	if sink, ok := optionsMap["pactl"]; ok {
		return fmt.Sprint(sink)
	}
	return "@DEFAULT_SINK@"
}

func parseLinuxVolume(volume any) (int, error) {
	var value float64
	switch v := volume.(type) {
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
		value = parsed
	case float64:
		value = v
	case int:
		value = float64(v)
	default:
		return 0, fmt.Errorf("volume is not number: %v", volume)
	}

	if value < 0 {
		return 0, errors.New("volume must be greater than or equal to 0")
	}
	if value <= 1 {
		value = value * 100
	}
	return int(value), nil
}
