package util

import (
	"errors"
	"sucicada/home/internal/logger"

	"resty.dev/v3"
)

type uRequest struct{}

var Request uRequest

// func sendGetRequest(apiPath string) (string, error) {
// 	return doRequest("GET", apiPath, nil)
// }

// func sendRequest(apiPath string) (string, error) {
// 	return doRequest("POST", apiPath, nil)
// }

// func sendRequestWithBody(apiPath string, body any) (string, error) {
// 	return doRequest("POST", apiPath, body)
// }

func (u *uRequest) Do(method, url string, body any) (string, error) {
	logger.Info("request url:", method, url)
	var err error

	req := resty.New().
		//SetDebug(true).
		// SetBaseURL(mediaConfig.APIBaseURL).
		R()
	if body != nil {
		req.SetBody(body)
	}

	var res *resty.Response
	switch method {
	case "GET":
		res, err = req.Get(url)
	case "POST":
		res, err = req.Post(url)
	default:
		return "", errors.New("unsupported HTTP method: " + method)
	}

	if err != nil {
		logger.Error("send request error: ", err)
		return "", err
	}
	if res.IsError() {
		logger.Error("controller response is not success: ", res.StatusCode(), res.String())
		return "", errors.New("controller response is not success")
	}
	return res.String(), nil
}
