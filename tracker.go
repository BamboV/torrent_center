package torrent_center

import (
	"github.com/bamboV/torrent/trackers/abstract_tracker/client"
	"github.com/bamboV/torrent"
)

type Tracker struct {
	client client.TrackerClient `gorm:"-"`
	Name string `gorm:"primary_key" json:"name"`
	Url string `json:"url"`
}

func (t *Tracker) Get(id int) torrent.Distribution{
	result, _ :=  t.client.GetTorrent(t.Url, id)

	return result
}

func (t *Tracker) Search(phrase string) []torrent.Distribution{
	result, _ := t.client.Search(t.Url, phrase)

	return result
}
