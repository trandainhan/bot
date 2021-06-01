package utils

import (
	"github.com/parnurzeal/gorequest"
	"net/http"
	"net/url"
	"time"
)

func HttpPost(url string, data interface{}, headers *map[string]string) (string, int, error) {
	request := gorequest.New()
	if headers != nil {
		for k, v := range *headers {
			request.Set(k, v)
		}
	}
	if data != nil {
		request.Send(data)
	}
	resp, result, errs := request.
		Post(url).
		Retry(
			5,
			time.Second*5,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
			http.StatusForbidden,
		).
		End()
	if len(errs) > 0 {
		return "", resp.StatusCode, errs[0]
	}

	return result, resp.StatusCode, nil
}

func BuildUrlWithParams(baseURL string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	u.RawQuery = values.Encode()
	return u.String(), nil
}

func HttpGet(url string, headers *map[string]string) (string, int, error) {
	request := gorequest.New()
	if headers != nil {
		for k, v := range *headers {
			request.Set(k, v)
		}
	}
	resp, body, errs := request.
		Get(url).
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		End()
	if len(errs) > 0 {
		return "", resp.StatusCode, errs[0]
	}
	return body, resp.StatusCode, nil
}
