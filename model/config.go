package model

import "github.com/expr-lang/expr/vm"

// Config is the configuration for the monitor
type Config struct {
	WebURL      string      `json:"web_url" yaml:"web_url"` // The URL of qBittorrent WebUI like https://example.com:8080
	APIKey      string      `json:"api_key" yaml:"api_key"` // API key
	GetTorrents GetTorrents `json:"get_torrents" yaml:"get_torrents"`
	Rules       []*Rule     `json:"rules" yaml:"rules"`
}

type Rule struct {
	Condition string        `json:"condition" yaml:"condition"` // The condition of the rule
	Action    TorrentAction `json:"action" yaml:"action"`       // The action of the rule
	Evaluator *vm.Program   `json:"-" yaml:"-"`
}

type GetTorrents struct {
	Limit   int    `json:"limit" yaml:"limit"`
	Sort    string `json:"sort" yaml:"sort"`
	Reverse bool   `json:"reverse" yaml:"reverse"`
}
