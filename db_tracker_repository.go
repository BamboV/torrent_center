package torrent_center

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DBTrackerRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) DBTrackerRepository {
	return DBTrackerRepository{
		db: db,
	}
}

func (r DBTrackerRepository) GetTrackers() []Tracker {
	trackers := []Tracker{}
	r.db.Find(&trackers)

	return trackers
}

func (r DBTrackerRepository) GetTracker(name string) Tracker {
	tracker := Tracker{}
	r.db.Find(&tracker, Tracker{Name: name})

	return tracker
}

func (r DBTrackerRepository) CreateTracker(t Tracker) Tracker {
	r.db.Create(&t)

	return t
}

func (r DBTrackerRepository) UpdateTracker(t Tracker) Tracker {
	r.db.Save(&t)
	return t
}

func (r DBTrackerRepository) DeleteTracker(t Tracker) {
	r.db.Delete(&t)
}



