package models

import "github.com/satori/go.uuid"

// CREATE TABLE startups (id UUID PRIMARY KEY, name string, description string, client_id uuid, conference_id uuid, is_active bool, sort_order int, sponsor_level_id uuid, 
// 	deleted_at timestamptz, created_at timestamptz, updated_at timestamptz, deleted bool, created_by uuid, updated_by uuid, deleted_by uuid, 
// 	twitter string, facebook string, linked_in string, youtube string);

type Startup struct {
	Base
	Name             string
	IsActive         bool
	StartupLevelID   uuid.UUID
	ClientID         uuid.UUID
	ConferenceID     uuid.UUID
	SortOrder int
	Description             string
	SocialMedia  SocialMedia `gorm:"embedded"`
	Tags        []*Tag       `gorm:"many2many:startups_tages;"`
}