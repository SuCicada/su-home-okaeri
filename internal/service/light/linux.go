package light

import (
	"SuCicada/home/internal/util"
	"fmt"
	"os"
)

type sLinuxLight struct{}

var LinuxLight = &sLinuxLight{}

const HIGH_LIGHT = 100
const LOW_LIGHT = 10

func (l *sLinuxLight) sshLinux(cmd string) string {
	return util.SSHRunRoot(util.SSHConfig{
		Host:     os.Getenv("LINUX_HOST"),
		User:     os.Getenv("LINUX_USER"),
		Password: os.Getenv("LINUX_PASSWORD"),
	}, cmd)
}
func (l *sLinuxLight) Get() int {
	res := l.sshLinux(`
		ddcutil --bus=5 getvcp 10 | grep -i "current value" | awk '{print $9}' | tr -d ','
	`)
	return util.StrToInt(res)
}

func (l *sLinuxLight) Toggle() {
	light := l.Get()
	if light < HIGH_LIGHT {
		l.Set(HIGH_LIGHT)
	} else {
		l.Set(LOW_LIGHT)
	}
}

func (l *sLinuxLight) Set(light int) {
	l.sshLinux(fmt.Sprintf("ddcutil --bus=5 setvcp 10 %d", light))
}
