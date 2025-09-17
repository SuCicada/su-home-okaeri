package util

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type uOTP struct{}

var OTP uOTP

func (u *uOTP) Generate(secret string) (string, error) {
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    60 * 15,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	return passcode, err
}

func (u *uOTP) Verify(secret, passcode string) bool {
	valid, err := totp.ValidateCustom(passcode, secret, time.Now(), totp.ValidateOpts{
		Period:    60 * 15,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	if err != nil {
		return false
	}
	return valid
}
