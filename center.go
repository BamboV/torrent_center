package torrent_center

import (
	"github.com/bamboV/torrent/trackers/abstract_tracker/client"
	torrentClient "github.com/bamboV/torrent/clients/abstract_torrent_client/client"
)

type Center struct {
	repo TrackerRepository
	trackerClient client.TrackerClient
	torrentClient torrentClient.Client
}

func NewCenter(repository TrackerRepository, client client.TrackerClient, torrentClient torrentClient.Client) Center {
	return Center {
		repo: repository,
		trackerClient: client,
		torrentClient: torrentClient,
	}
}

func(c *Center) GetTracker(name string) Tracker{
	tracker := c.repo.GetTracker(name)
	c.setClient(&tracker)
	return tracker
}

func(c *Center) AddTracker(tracker Tracker) Tracker {
	t := c.repo.CreateTracker(tracker)
	c.setClient(&t)
	return t
}

func(c *Center) GetTrackers() []Tracker {
	trackers := c.repo.GetTrackers()
	for _, value := range trackers {
		c.setClient(&value)
	}
	return trackers
}

func(c *Center) UpdateTracker(tracker Tracker) Tracker {
	t := c.repo.UpdateTracker(tracker)
	c.setClient(&t)
	return t
}

func(c *Center) DeleteTracker(tracker Tracker) {
	c.repo.DeleteTracker(tracker)
}

func(c *Center) Download(magnet string) bool{
	return c.torrentClient.DownloadByMagnet(magnet)
}

func(c *Center) setClient(tracker *Tracker) {
	tracker.client = c.trackerClient
}
