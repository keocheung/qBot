package model

// Config is the configuration for the monitor
type Config struct {
	WebURL string `json:"web_url"` // The URL of qBittorrent WebUI like https://example.com:8080
	APIKey string `json:"api_key"` // API key
}
