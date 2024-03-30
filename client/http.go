package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"qbot/internal/meta"
)

type HTTPClient interface {
	Get(url string, headers map[string]string) ([]byte, error)
	Post(url string, body []byte, headers map[string]string) ([]byte, error)
}

func NewHTTPClient() HTTPClient {
	return &httpClient{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			Timeout: 60 * time.Second,
		},
	}
}

type httpClient struct {
	client *http.Client
}

func (c *httpClient) Get(url string, headers map[string]string) ([]byte, error) {
	return c.request(url, http.MethodGet, nil, headers)
}

func (c *httpClient) Post(url string, body []byte, headers map[string]string) ([]byte, error) {
	return c.request(url, http.MethodPost, body, headers)
}

func (c *httpClient) request(url string, method string, body []byte, headers map[string]string) ([]byte, error) {
	var req *http.Request
	var err error
	if len(body) > 0 {
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("%s %s error: %v", method, url, err)
	}
	req.Header.Set("User-Agent", meta.UserAgent)
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s %s error: %v", method, url, err)
	}
	defer rsp.Body.Close()
	rspBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s %s error: %v", method, url, err)
	}
	if rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s %s error: %v %s", method, url, rsp.StatusCode, rspBody)
	}
	return rspBody, nil
}
