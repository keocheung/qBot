package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"qbot/model"
	"qbot/util/logger"
)

type QbClient interface {
	// GetTorrents gets the torrent list
	//
	// See https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)#get-torrent-list
	GetTorrents(model.Options) ([]model.Torrent, error)

	// StopTorrents stops torrents with hash.
	//
	// See https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)#pause-torrents
	StopTorrents(hashes []string) error

	// SetShareLimits sets torrent share limit with the torrents' hash.
	//
	// See https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)#set-torrent-share-limit
	SetShareLimits(hashes []string, ratioLimit float32, timeLimit int) error
}

func NewQbClient(webURL, apiKey string) (QbClient, error) {
	u, err := url.Parse(webURL)
	if err != nil {
		logger.Infof("url.Parse error: %v", err)
		return nil, err
	}
	return &qbClient{
		webURL:     u,
		apiKey:     apiKey,
		httpClient: NewHTTPClient(),
	}, nil
}

type qbClient struct {
	webURL     *url.URL
	apiKey     string
	httpClient HTTPClient
}

func (c *qbClient) GetTorrents(options model.Options) ([]model.Torrent, error) {
	u := c.webURL.JoinPath("/api/v2/torrents/info")
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

func (c *qbClient) StopTorrents(hashes []string) error {
	u := c.webURL.JoinPath("/api/v2/torrents/stop")
	param := url.Values{}
	param.Set("hashes", strings.Join(hashes, "|"))
	_, err := c.httpClient.Post(u.String(), []byte(param.Encode()), map[string]string{
		"Cookie":       "SID=" + c.apiKey,
		"Content-Type": "application/x-www-form-urlencoded",
	})
	if err != nil {
		logger.Infof("httpClient.Post error: %v", err)
		return err
	}
	return nil
}

func (c *qbClient) SetShareLimits(hashes []string, ratioLimit float32, timeLimit int) error {
	u := c.webURL.JoinPath("/api/v2/torrents/setShareLimits")
	param := url.Values{}
	param.Set("hashes", strings.Join(hashes, "|"))
	param.Set("ratioLimit", fmt.Sprintf("%.1f", ratioLimit))
	param.Set("seedingTimeLimit", strconv.Itoa(timeLimit))
	param.Set("inactiveSeedingTimeLimit", "-1")
	_, err := c.httpClient.Post(u.String(), []byte(param.Encode()), map[string]string{
		"Cookie":       "SID=" + c.apiKey,
		"Content-Type": "application/x-www-form-urlencoded",
	})
	if err != nil {
		logger.Infof("httpClient.Post error: %v", err)
		return err
	}
	return nil
}
