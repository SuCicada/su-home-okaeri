package raspberrypi

import (
	"fmt"
	"strconv"
	"strings"
	"sucicada/home/internal/consts"
	"sucicada/home/internal/structs"
	commandstructs "sucicada/home/internal/structs/command"
)

type sRaspberryPiLight struct{}

const powerOnBrightness = 100

func (l *sRaspberryPiLight) GetBrightness() (int, error) {
	status, err := l.GetStatus()
	if err != nil {
		return 0, err
	}
	return status.Brightness, nil
}

func (l *sRaspberryPiLight) GetStatus() (structs.LightStatus, error) {
	res, err := runTVCtl("status")
	if err != nil {
		return structs.LightStatus{}, err
	}

	power, err := parsePower(res)
	if err != nil {
		return structs.LightStatus{}, err
	}

	brightness := 0
	if power {
		brightness = powerOnBrightness
	}
	return structs.LightStatus{
		Power:      power,
		Brightness: brightness,
	}, nil
}

func (l *sRaspberryPiLight) SetBrightness(command commandstructs.LightCommand) error {
	if command.Light > 0 {
		return l.On()
	}
	return l.Off()
}

func (l *sRaspberryPiLight) On() error {
	_, err := runTVCtl("on")
	return err
}

func (l *sRaspberryPiLight) Off() error {
	_, err := runTVCtl("off")
	return err
}

func runTVCtl(action string) (string, error) {
	script, err := tvctlScript()
	if err != nil {
		return "", err
	}
	return ssh(fmt.Sprintf("python %s %s", strconv.Quote(script), action))
}

func tvctlScript() (string, error) {
	script, _ := controlOptions(consts.CONTROL_LIGHT)["tvctl_script"].(string)
	script = strings.TrimSpace(script)
	if script == "" {
		return "", fmt.Errorf("raspberrypi light options.tvctl_script is empty")
	}
	return script, nil
}

func parsePower(out string) (bool, error) {
	fields := strings.Fields(strings.ToLower(out))
	if len(fields) == 0 {
		return false, fmt.Errorf("empty tvctl status")
	}

	switch fields[len(fields)-1] {
	case "on", "1", "true":
		return true, nil
	case "off", "0", "false", "unkonwn":
		return false, nil
	default:
		return false, fmt.Errorf("unknown tvctl status: %s", strings.TrimSpace(out))
	}
}
