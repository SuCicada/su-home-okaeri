package light

import (
	"SuCicada/home/internal/logger"
	"SuCicada/home/internal/util"
	"fmt"
	"os"
)

type sLinuxLight struct{}

var LinuxLight = &sLinuxLight{}

const HIGH_LIGHT = 100
const LOW_LIGHT = 10

func (l *sLinuxLight) sshLinux(cmd string) (string, error) {
	return util.SSHRunRoot(util.SSHConfig{
		Host:     os.Getenv("LINUX_HOST"),
		User:     os.Getenv("LINUX_USER"),
		Password: os.Getenv("LINUX_PASSWORD"),
	}, cmd)
}

func (l *sLinuxLight) Get() (int, error) {
	res, err := l.sshLinux(`
		ddcutil --bus=5 getvcp 10 | grep -i "current value" | awk '{print $9}' | tr -d ','
	`)
	if err != nil {
		logger.Error("Error getting light:", err)
		return 0, err
	}
	return util.StrToInt(res), nil
}

func (l *sLinuxLight) Toggle() (string, error) {
	light, err := l.Get()
	if err != nil {
		return "", err
	}
	if light < HIGH_LIGHT {
		return l.Set(HIGH_LIGHT)
	} else {
		return l.Set(LOW_LIGHT)
	}
}

func (l *sLinuxLight) Set(light int) (string, error) {
	return l.sshLinux(fmt.Sprintf("ddcutil --bus=5 setvcp 10 %d", light))
}
