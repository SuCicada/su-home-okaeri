package redmi

import (
	"encoding/json"
	"errors"
	"fmt"
	"sucicada/home/internal/logger"
	structs "sucicada/home/internal/structs/command"
)

type sRedmiLight struct{}

func (l *sRedmiLight) Set(command string) error {
	var data structs.LightCommand
	err := json.Unmarshal([]byte(command), &data)
	if err != nil {
		logger.Error("Failed to unmarshal command: ", err)
		return err
	}

	light := data.Light
	_, err = ssh(fmt.Sprintf("termux-brightness %d", light))
	return err
}

func (l *sRedmiLight) Get() (any, error) {
	logger.Warn("not support get redmi light")
	return -1, errors.New("not support get redmi light")
}

func (l *sRedmiLight) Toggle() (string, error) {
	logger.Warn("not support toggle redmi light")
	return "", errors.New("not support toggle redmi light")
}
