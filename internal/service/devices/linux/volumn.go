package linux

import (
	"SuCicada/home/internal/consts"
	"SuCicada/home/internal/util"
	"fmt"
)

type sVolume struct{}

var Volume = &sVolume{}

func getPactlOpts() int {
	options := Config.Control[consts.CONTROL_VOLUME].Options
	pactlOpts := options["pactl"]
	return pactlOpts.(int)
}

func (l *sVolume) Get() (int, error) {
	pactlOpts := getPactlOpts()
	res, err := sshLinux(fmt.Sprintf("pactl get-sink-volume %d", pactlOpts))
	if err != nil {
		return 0, err
	}
	return util.StrToInt(res), nil
}

func (l *sVolume) Set(volume int) error {
	pactlOpts := getPactlOpts()
	_, err := sshLinux(fmt.Sprintf("pactl set-sink-volume %d %d", pactlOpts, volume))
	return err
}
