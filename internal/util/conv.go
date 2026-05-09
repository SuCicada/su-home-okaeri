package util

import (
	"encoding/json"
	"strconv"
	"strings"
	"sucicada/home/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type uConv struct{}

var Conv = uConv{}

func (c *uConv) StrToInt(v string) int {
	i, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		logger.Warn("not number: ", v)
		return 0
	}
	return i
}
func (c *uConv) ToBytes(v interface{}) []byte {
	var bytes []byte
	switch v := v.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		bytes, _ = json.Marshal(v)
	}

	return bytes
}
func (c *uConv) ToJsonStr(v any) string {
	jsonStr, _ := json.Marshal(v)
	return string(jsonStr)
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

func (c *uConv) AnyToMap(v any) map[string]any {
	var result map[string]any
	mapstructure.Decode(v, &result)
	return result
}

func (c *uConv) MapToAny(mapdata any, target any) error {
	cfg := &mapstructure.DecoderConfig{
		Result:           target,
		TagName:          "json",
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return err
	}
	return decoder.Decode(mapdata)
}
