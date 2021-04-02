package models

import "github.com/satori/go.uuid"

type Ambassador struct {
	//ID                int `gorm:"primary_key"`
	Base
	Title             string
	ManagerID         uuid.UUID //manager id from users table
	AmbassadorGroupID uuid.UUID
	UserID            uuid.UUID
	// CreatedAt         time.Time
	// UpdatedAt         time.Time
	// DeletedAt         *time.Time
	IsActive          bool
	//ConferenceID     int //must be set for current active, client can have only one active confercne.
	Conferences      []*Conference `gorm:"many2many:conferences_ambassadors;"`
	ClientID         uuid.UUID
	Sessions         []*Session `gorm:"many2many:session_ambassadors;"`
	DisplayLevelSort int
	DisplayInList    bool
	DisplayTitle     string
}