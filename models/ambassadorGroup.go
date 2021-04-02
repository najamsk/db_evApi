package models

import "github.com/satori/go.uuid"

type AmbassadorGroup struct {
	Base
	Title       string
	Ambassadors []Ambassador
	IsActive    bool
	//Conferences      []*Conference `gorm:"many2many:organizer_conferences;"`
	ClientID         uuid.UUID
	//Sessions         []*Session `gorm:"many2many:session_ambassador;"`
	DisplayLevelSort int
	DisplayInList    bool
	DisplayTitle     string

	// CreatedAt   time.Time
	// UpdatedAt   time.Time
	// DeletedAt   time.Time
}