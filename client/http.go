package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"qbot/util/logger"
)

type HTTPClient interface {
	Get(url string, headers map[string]string) ([]byte, error)
	Post(url string, body []byte, headers map[string]string) ([]byte, error)
}

var client = &http.Client{}

func NewHTTPClient() HTTPClient {
	return &httpClient{}
}

type httpClient struct {
}

func (c *httpClient) Get(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Errorf("get %s error: %v", url, err)
		return nil, err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	rsp, err := client.Do(req)
	if err != nil {
		logger.Errorf("get %s error: %v", url, err)
		return nil, err
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logger.Errorf("get %s error: %v", url, err)
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get %s error: %v %s", url, rsp.StatusCode, body)
	}
	return body, nil
}

func (c *httpClient) Post(url string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logger.Infof("post %s error: %v", url, err)
		return nil, err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	rsp, err := client.Do(req)
	if err != nil {
		logger.Infof("post %s error: %v", url, err)
		return nil, err
	}
	defer rsp.Body.Close()
	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logger.Infof("post %s error: %v", url, err)
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post %s error: %v %s", url, rsp.StatusCode, rspBody)
	}
	return rspBody, nil
}
