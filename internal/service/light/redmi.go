package light

import (
	"SuCicada/home/internal/logger"
	"SuCicada/home/internal/util"
	"fmt"
	"os"
)

type sRedmiLight struct{}

var RedmiLight = &sRedmiLight{}

func (l *sRedmiLight) Set(light int) {
	util.SSHRun(util.SSHConfig{
		Host: os.Getenv("REDMI_HOST"),
		Port: util.GetInt("REDMI_PORT"),
		// User:     os.Getenv("REDMI_USER"),
		Password: os.Getenv("REDMI_PASSWORD"),
	}, fmt.Sprintf("termux-brightness %d", light))
}

func (l *sRedmiLight) Get() int {
	logger.Warn("not support get redmi light")
	return -1
}

func (l *sRedmiLight) Toggle() {
	logger.Warn("not support toggle redmi light")
}
