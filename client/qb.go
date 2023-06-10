package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"qb-monitor/model"
	"qb-monitor/util/logger"
)

type QbClient interface {
	// GetTorrents returns the torrent list
	GetTorrents(model.Options) ([]model.Torrent, error)
	// SetShareLimits sets the share limit for a torrent list
	SetShareLimits(hashes []string, ratioLimit float32, timeLimit int) error
}

func NewQbClient(webURL, apiKey string) QbClient {
	return &qbClient{
		webURL:     webURL,
		apiKey:     apiKey,
		httpClient: NewHTTPClient(),
	}
}

type qbClient struct {
	webURL     string
	apiKey     string
	httpClient HTTPClient
}

func (c *qbClient) GetTorrents(options model.Options) ([]model.Torrent, error) {
	u, err := url.Parse(c.webURL)
	if err != nil {
		logger.Errorf("url.Parse error: %v", err)
		return nil, err
	}
	u = u.JoinPath("/api/v2/torrents/info")
	params := url.Values{}
	params.Set("limit", strconv.Itoa(options.Limit))
	params.Set("sort", options.Sort)
	params.Set("reverse", strconv.FormatBool(options.Reverse))
	u.RawQuery = params.Encode()
	rsp, err := c.httpClient.Get(u.String(), map[string]string{
		"Cookie": "SID=" + c.apiKey,
	})
	if err != nil {
		logger.Errorf("httpClient.Get error: %v", err)
		return nil, err
	}
	var result = []model.Torrent{}
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		logger.Errorf("json.Unmarshal error: %v", err)
		return nil, err
	}
	return result, nil
}

func (c *qbClient) SetShareLimits(hashes []string, ratioLimit float32, timeLimit int) error {
	u, err := url.Parse(c.webURL)
	if err != nil {
		logger.Infof("url.Parse error: %v", err)
		return err
	}
	u = u.JoinPath("/api/v2/torrents/setShareLimits")
	data := fmt.Sprintf("hashes=%s&ratioLimit=%f&seedingTimeLimit=%d",
		strings.Join(hashes, "|"), ratioLimit, timeLimit)
	_, err = c.httpClient.Post(u.String(), []byte(data), map[string]string{
		"Cookie":       "SID=" + c.apiKey,
		"Content-Type": "application/x-www-form-urlencoded",
	})
	if err != nil {
		logger.Infof("httpClient.Post error: %v", err)
		return err
	}
	return nil
}
