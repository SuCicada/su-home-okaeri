package controller

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/response"
	"SuCicada/home/internal/util"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/SuCicada/apprise-sdk-go/apprise"
	"github.com/gin-gonic/gin"
	"resty.dev/v3"
)

type cSmsCheck struct{}

var SmsCheck = cSmsCheck{}

const TEMP_FILE = "/tmp/sms-check.txt"

func (c *cSmsCheck) Webhook(ctx *gin.Context) {

	req := struct {
		Body string `json:"body"`
	}{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	body := req.Body
	log.Println("webhook receive:", body)

	os.WriteFile(TEMP_FILE, []byte(body), 0644)
	response.Success(ctx)
}

const SMS_Name = "短信测试"

func (c *cSmsCheck) SendVerifyCode(ctx *gin.Context) {
	code, err := util.OTP.Generate(cfg.GetConfig().SMSCheck.Secret)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	res, err := resty.New().R().
		SetBody(map[string]string{
			"name": SMS_Name,
			"code": code,
		}).
		Post(cfg.GetConfig().SMSCheck.PushUrl)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if res.StatusCode() != 200 {
		response.Bad(ctx, fmt.Errorf("statusCode: %d", res.StatusCode()))
		return
	}

	response.Success(ctx)
}

func (c *cSmsCheck) CheckSMS(ctx *gin.Context) {

	body, err := os.ReadFile(TEMP_FILE)
	if err != nil {
		util.Alert.SendApprise(apprise.Message{
			Title: "❌[132短信] check failed",
			Body:  err.Error(),
			Tag:   "job",
			Type:  apprise.TypeFailure,
		})
		response.Error(ctx, err)
		return
	}

	text := string(body)
	if DoCheckSMS(text) == false {
		log.Println("检查失败：未收到预期的验证消息或验证超时")
		util.Alert.SendApprise(apprise.Message{
			Title: "❌[132短信] inactive",
			Body:  "流れは問題がある",
			Tag:   "job",
			Type:  apprise.TypeFailure,
		})

		response.Success(ctx, "check over, but failed")
	} else {
		response.Success(ctx)
	}
}

func DoCheckSMS(text string) bool {
	re := regexp.MustCompile(`【Spug推送】(.+)欢迎您，您的验证码为(\d{6})`)
	match := re.FindStringSubmatch(text)
	if len(match) != 3 {
		log.Println("match not match", match)
		return false
	}

	name, code := match[1], match[2]

	if name != SMS_Name {
		log.Println("name not match", name, SMS_Name)
		return false
	}

	secret := cfg.GetConfig().SMSCheck.Secret
	res := util.OTP.Verify(secret, code)
	if !res {
		log.Println("verify error", code)
		return false
	}

	return true
}
