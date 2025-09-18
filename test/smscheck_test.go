package test

import (
	"SuCicada/home/internal/controller"
	"testing"

	"SuCicada/home/internal/cfg"

	"github.com/stretchr/testify/suite"
)

type TestSmscheck struct {
	suite.Suite
}

func TestSmscheckRun(t *testing.T) {
	cfg.CONFIG_PATH = "../config.yaml"
	suite.Run(t, new(TestSmscheck))
}

func (t *TestSmscheck) TestCheck() {
	res := controller.DoCheckSMS(`【Spug推送】短信测试欢迎您，您的验证码为991814，如非本人操作请忽略。
2025-09-18 23:09:22`)
	t.True(res)
}
