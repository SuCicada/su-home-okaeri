package windows

import (
	"SuCicada/home/internal/util"
	"fmt"
)

type sWindowsVolume struct{}

var WindowsVolume = &sWindowsVolume{}

func (l *sWindowsVolume) Get() (int, error) {
	res, err := ssh(`
		windows-volume-controller.exe
	`)
	if err != nil {
		return 0, err
	}
	return util.StrToInt(res), nil
}

func (l *sWindowsVolume) Set(volume int) error {
	_, err := ssh(fmt.Sprintf("windows-volume-controller.exe %d", volume))
	return err
}
