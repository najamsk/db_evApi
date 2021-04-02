package models

import ("time" 
"github.com/satori/go.uuid")

type Conference struct {
	Base
	Title                string
	Summary              string
	Details              string
	//Poster               string
	//Thumbnail            string
	IsActive             bool
	StartDate            time.Time
	StartDateDisplay     string
	EndDate              time.Time
	EndDateDisplay       string
	DurationDisplay      string
	Address              string
	GeoLocation          GeoLocation `gorm:"embedded"`
	ClientID             uuid.UUID
	Sessions             []Session
	FloorPlanPoster      string
	ExibitionStartupPlan string
	ExibitionSponsorPlan string
	Speakers    []*User `gorm:"many2many:conference_speakers;"`
	Attendees    []*User `gorm:"many2many:conference_attendees;"`
	Organizers    []*User `gorm:"many2many:conference_organizers;"`
	Users  []*User `gorm:"many2many:conferences_users;"`
	Venues []*Venue `gorm:"ForeignKey:ConferenceID"`
	SocialMedia  SocialMedia `gorm:"embedded"`

	// CreatedAt            time.Time
	// UpdatedAt            time.Time
	// DeletedAt            time.Time
}