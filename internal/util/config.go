package util

import (
	"SuCicada/home/internal/logger"
	"os"
	"strconv"
	"strings"
)

func StrToInt(v string) int {
	i, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		logger.Warn("not number: ", v)
		return 0
	}
	return i
}
func GetInt(key string) int {
	v := os.Getenv(key)
	if v == "" {
		return 0
	}
	return StrToInt(v)
}
