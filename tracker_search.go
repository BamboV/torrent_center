package torrent_center

import "github.com/bamboV/torrent"

type TrackerSearchResult struct {
	Name string `json:"name"`
	Items []torrent.Distribution `json:"items"`
}
