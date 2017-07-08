package torrent_center

type TrackerRepository interface {
	GetTrackers() []Tracker
	GetTracker(name string) Tracker
	CreateTracker(t Tracker) Tracker
	UpdateTracker(t Tracker) Tracker
	DeleteTracker(t Tracker)
}
