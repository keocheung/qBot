package model

// Torrent is a qBittorrent torrent
type Torrent struct {
	Hash           string  `json:"hash"`             // Torrent hash
	Catogory       string  `json:"category"`         // Category of the torrent
	Tags           string  `json:"tags"`             // Comma-concatenated tag list of the torrent
	Ratio          float32 `json:"ratio"`            // Torrent share ratio. Max ratio value: 9999.
	MaxRatio       float32 `json:"max_ratio"`        // Maximum share ratio until torrent is stopped from seeding/uploading
	MaxSeedingTime int     `json:"max_seeding_time"` // Maximum seeding time (seconds) until torrent is stopped from seeding
}

// Options is the query options	for GetTorrents
type Options struct {
	Limit   int    `json:"limit"`   // Limit the number of torrents returned
	Sort    string `json:"sort"`    // Sort torrents by given key. They can be sorted using any field of the response's JSON array (which are documented below) as the sort key.
	Reverse bool   `json:"reverse"` // Enable reverse sorting. Defaults to false
}
