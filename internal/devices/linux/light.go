package linux

import (
	"fmt"
	"strings"
	"sucicada/home/internal/consts"
	"sucicada/home/internal/logger"
	"sucicada/home/internal/structs"
	commandstructs "sucicada/home/internal/structs/command"
	"sucicada/home/internal/util"
)

type sLinuxLight struct{}

func getOpts() string {
	options := Config.Control[consts.CONTROL_LIGHT]["options"]
	optionsMap := util.Conv.AnyToMap(options)
	if pactlOpts, ok := optionsMap["ddcutil"]; ok {
		return pactlOpts.(string)
	}
	return optionsMap["ddcutil"].(string)
}

func (l *sLinuxLight) GetBrightness() (int, error) {
	opts := getOpts()
	res, err := sshLinux(fmt.Sprintf(`
		ddcutil %s getvcp 10 | grep -i "current value" | awk '{print $9}' | tr -d ','
	`, opts))

	if err != nil {
		logger.Error("Error getting light:", err)
		return 0, err
	}
	return util.Conv.StrToInt(res), nil
}

func (l *sLinuxLight) GetStatus() (structs.LightStatus, error) {
	brightness, err := l.GetBrightness()
	if err != nil {
		return structs.LightStatus{}, err
	}

	opts := getOpts()
	res, err := sshLinux(fmt.Sprintf(
		`ddcutil %s getvcp D6 | grep -i "current value" | awk '{print $9}' | tr -d ','`,
		opts,
	))
	if err != nil {
		return structs.LightStatus{
			Power:      brightness > 0,
			Brightness: brightness,
		}, nil
	}

	powerCode := strings.TrimSpace(res)
	power := powerCode == "01" || (powerCode != "05" && brightness > 0)

	return structs.LightStatus{
		Power:      power,
		Brightness: brightness,
	}, nil
}

func (l *sLinuxLight) SetBrightness(command commandstructs.LightCommand) error {
	_, err := sshLinux(fmt.Sprintf(`
	ddcutil %s setvcp 10 %d
	`, getOpts(), command.Light))

	return err
}

func (l *sLinuxLight) On() error {
	return l.setPowerMode("01")
}

func (l *sLinuxLight) Off() error {
	return l.setPowerMode("05")
}

func (l *sLinuxLight) setPowerMode(value string) error {
	_, err := sshLinux(fmt.Sprintf(`
	ddcutil %s setvcp D6 %s
	`, getOpts(), value))
	return err
}
