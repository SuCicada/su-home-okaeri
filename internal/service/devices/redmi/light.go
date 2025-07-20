package redmi

import (
	"SuCicada/home/internal/logger"
	"errors"
	"fmt"
)

type sRedmiLight struct{}

func (l *sRedmiLight) Set(light int) error {
	_, err := ssh(fmt.Sprintf("termux-brightness %d", light))
	return err
}

func (l *sRedmiLight) Get() (int, error) {
	logger.Warn("not support get redmi light")
	return -1, errors.New("not support get redmi light")
}

func (l *sRedmiLight) Toggle() (string, error) {
	logger.Warn("not support toggle redmi light")
	return "", errors.New("not support toggle redmi light")
}
