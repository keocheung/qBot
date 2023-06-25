package model

import (
	"strings"
)

// Torrent is a qBittorrent torrent
type Torrent struct {
	Hash           string  `json:"hash"`             // Torrent hash
	Name           string  `json:"name"`             // Torrent name
	Category       string  `json:"category"`         // Category of the torrent
	Tags           Tags    `json:"tags"`             // Comma-concatenated tag list of the torrent
	Ratio          float32 `json:"ratio"`            // Torrent share ratio. Max ratio value: 9999
	MaxRatio       float32 `json:"max_ratio"`        // Maximum share ratio until torrent is stopped from seeding/uploading
	MaxSeedingTime int     `json:"max_seeding_time"` // Maximum seeding time (seconds) until torrent is stopped from seeding
}

type TorrentAction struct {
	MaxRatio *float32 `json:"max_ratio" yaml:"max_ratio"`
}

type Tags []string

func (t *Tags) UnmarshalJSON(b []byte) error {
	*t = strings.Split(string(b[1:len(b)-1]), ", ")
	return nil
}

// Options is the query options	for GetTorrents
type Options struct {
	Limit   int    `json:"limit"`   // Limit the number of torrents returned
	Sort    string `json:"sort"`    // Sort torrents by given key. They can be sorted using any field of the response's JSON array (which are documented below) as the sort key.
	Reverse bool   `json:"reverse"` // Enable reverse sorting. Defaults to false
}
