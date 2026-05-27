package linuxutil

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// var execCmd func(cmdArgs []string) ([]byte, error) = nil

func execCmd(cmdArgs []string) ([]byte, error) {
	res, err := sshLinux(strings.Join(cmdArgs, " "))
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

func init() {
	if _, err := execCmd([]string{"pactl", "info"}); err != nil {
		useAmixer = true
	}
}

// GetVolume returns the current volume (0 to 100).
func GetVolume() (int, error) {
	out, err := execCmd(getVolumeCmd())
	if err != nil {
		return 0, err
	}
	return parseVolume(string(out))
}

// SetVolume sets the sound volume to the specified value.
func SetVolume(volume int) error {
	if volume < 0 || 100 < volume {
		return errors.New("out of valid volume range")
	}
	_, err := execCmd(setVolumeCmd(volume))
	return err
}

// IncreaseVolume increases (or decreases) the audio volume by the specified value.
func IncreaseVolume(diff int) error {
	_, err := execCmd(increaseVolumeCmd(diff))
	return err
}

// GetMuted returns the current muted status.
func GetMuted() (bool, error) {
	out, err := execCmd(getMutedCmd())
	if err != nil {
		return false, err
	}
	return parseMuted(string(out))
}

// Mute mutes the audio.
func Mute() error {
	_, err := execCmd(muteCmd())
	return err
}

// Unmute unmutes the audio.
func Unmute() error {
	_, err := execCmd(unmuteCmd())
	return err
}

var useAmixer bool

func cmdEnv() []string {
	return []string{"LANG=C", "LC_ALL=C"}
}

func getVolumeCmd() []string {
	if useAmixer {
		return []string{"amixer", "get", "Master"}
	}
	return []string{"pactl", "list", "sinks"}
}

func getPulseAudioDefaultSink() (string, error) {
	out, err := execCmd([]string{"pactl", "info"})
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")

	defaultSinkStr := "Default Sink: "
	for _, line := range lines {
		s := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(s, defaultSinkStr) {
			return strings.TrimSpace(strings.Replace(s, defaultSinkStr, "", 1)), nil
		}
	}
	return "", errors.New("could not find PulseAudio Default Sink")
}

func parseVolume(out string) (int, error) {
	var volumePattern = regexp.MustCompile(`\d+%`)

	sinkName, sinkNameErr := getPulseAudioDefaultSink()
	isDefaultSink := sinkNameErr != nil

	lines := strings.Split(out, "\n")

	for _, line := range lines {
		s := strings.TrimLeft(line, " \t")

		if !useAmixer && !isDefaultSink && strings.Contains(s, "Name: "+sinkName) {
			isDefaultSink = true
		}

		if useAmixer && strings.Contains(s, "Playback") && strings.Contains(s, "%") ||
			!useAmixer && isDefaultSink && strings.HasPrefix(s, "Volume:") {
			volumeStr := volumePattern.FindString(s)
			return strconv.Atoi(volumeStr[:len(volumeStr)-1])
		}
	}
	return 0, errors.New("no volume found")
}

func setVolumeCmd(volume int) []string {
	if useAmixer {
		return []string{"amixer", "set", "Master", strconv.Itoa(volume) + "%"}
	}
	return []string{"pactl", "set-sink-volume", "@DEFAULT_SINK@", strconv.Itoa(volume) + "%"}
}

func increaseVolumeCmd(diff int) []string {
	var sign string
	if diff >= 0 {
		sign = "+"
	} else if useAmixer {
		diff = -diff
		sign = "-"
	}
	if useAmixer {
		return []string{"amixer", "set", "Master", strconv.Itoa(diff) + "%" + sign}
	}
	return []string{"pactl", "--", "set-sink-volume", "@DEFAULT_SINK@", sign + strconv.Itoa(diff) + "%"}
}

func getMutedCmd() []string {
	if useAmixer {
		return []string{"amixer", "get", "Master"}
	}
	return []string{"pactl", "list", "sinks"}
}

func parseMuted(out string) (bool, error) {
	sinkName, sinkNameErr := getPulseAudioDefaultSink()
	isDefaultSink := sinkNameErr != nil

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		s := strings.TrimLeft(line, " \t")

		if !useAmixer && !isDefaultSink && strings.Contains(s, "Name: "+sinkName) {
			isDefaultSink = true
		}

		if useAmixer && strings.Contains(s, "Playback") && strings.Contains(s, "%") ||
			!useAmixer && isDefaultSink && strings.HasPrefix(s, "Mute: ") {
			if strings.Contains(s, "[off]") || strings.Contains(s, "yes") {
				return true, nil
			} else if strings.Contains(s, "[on]") || strings.Contains(s, "no") {
				return false, nil
			}
		}
	}
	return false, errors.New("no muted information found")
}

func muteCmd() []string {
	if useAmixer {
		return []string{"amixer", "-D", "pulse", "set", "Master", "mute"}
	}
	return []string{"pactl", "set-sink-mute", "@DEFAULT_SINK@", "1"}
}

func unmuteCmd() []string {
	if useAmixer {
		return []string{"amixer", "-D", "pulse", "set", "Master", "unmute"}
	}
	return []string{"pactl", "set-sink-mute", "@DEFAULT_SINK@", "0"}
}
