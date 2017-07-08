package torrent_center

import "github.com/bamboV/torrent"

type TrackerSearchResult struct {
	Name string
	Items []torrent.Distribution
}
