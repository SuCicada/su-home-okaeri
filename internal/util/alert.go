package util

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/logger"
	"context"
	"fmt"

	"github.com/SuCicada/apprise-sdk-go/apprise"
)

type uAlert struct {
}

var Alert uAlert

func (u uAlert) SendApprise(appriseMessage apprise.Message) error {
	fmt.Println("send alert", appriseMessage)
	url := cfg.GetConfig().Alert.Apprise
	err := apprise.Send(context.Background(), url, appriseMessage)
	if err != nil {
		logger.Error("send alert error", err)
		return err
	}
	return nil
}
