package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
)

func HttpPost(url string, data interface{}, headers *map[string]string) (string, int, error) {
	request := gorequest.New()
	request = request.Post(url)
	if headers != nil {
		for k, v := range *headers {
			request.Set(k, v)
		}
	}
	if data != nil {
		request.Send(data)
	}
	resp, body, errs := request.
		Retry(
			3,
			time.Second*5,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		).
		End()
	if errs != nil {
		log.Printf("HttpPost Full errors: %v", errs)
		if resp == nil {
			return "", 0, errs[0]
		}
		return "", resp.StatusCode, errs[0]
	}
	return body, resp.StatusCode, nil
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
	request = request.Get(url)
	if headers != nil {
		for k, v := range *headers {
			request.Set(k, v)
		}
	}
	resp, body, errs := request.
		Retry(
			3,
			time.Second*5,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		).End()
	if errs != nil {
		log.Printf("HttpGet Full errors: %v", errs)
		if resp == nil {
			return "", 0, errs[0]
		}
		return "", resp.StatusCode, errs[0]
	}
	return body, resp.StatusCode, nil
}

func HttpDelete(url string, data interface{}, headers *map[string]string) (string, int, error) {
	request := gorequest.New()
	request = request.Delete(url)
	if headers != nil {
		for k, v := range *headers {
			request.Set(k, v)
		}
	}
	if data != nil {
		request.Send(data)
	}
	resp, body, errs := request.
		Retry(
			3,
			time.Second*5,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		).
		End()
	if errs != nil {
		log.Printf("HttpDelete Full errors: %v", errs)
		if resp == nil {
			return "", 0, errs[0]
		}
		return "", resp.StatusCode, errs[0]
	}
	return body, resp.StatusCode, nil
}

func BuildQueryStringFromMap(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	return values.Encode()
}

func ProcessResponse(respBody string, result interface{}) error {
	err := json.Unmarshal([]byte(respBody), result)
	if err != nil {
		log.Printf("Error processing response: %s", err)
		return err
	}
	return nil
}
