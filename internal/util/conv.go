package util

import (
	"SuCicada/home/internal/logger"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type uConv struct{}

var Conv = uConv{}

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

func (c *uConv) GetMapFromGinContext(ginC *gin.Context) (map[string]string, error) {

	var result = make(map[string]string)
	switch ginC.Request.Header.Get("Content-Type") {
	case gin.MIMEJSON:
		if err := ginC.ShouldBindJSON(&result); err != nil {
			return nil, err
		}

	case gin.MIMEPOSTForm:
		if err := ginC.Request.ParseForm(); err != nil {
			return nil, err
		}
		for k, v := range ginC.Request.PostForm {
			if len(v) > 0 {
				result[k] = v[0]
			}
		}
	}

	return result, nil
}
