package models

import ("time" 
"github.com/satori/go.uuid")
/*
CREATE TABLE venues (id UUID PRIMARY KEY, conference_id UUID, name string, start_date timestamptz, start_date_display string, end_date timestamptz, end_date_display string,
	duration_display string, address string, geo_location_lat  numeric, geo_location_long numeric, radius  numeric, 
	country string, city string, county string, post_code string, state string,
	created_at timestamptz, 
   updated_at timestamptz null, deleted_at timestamptz null, is_active bool, deleted bool, created_by uuid NULL, updated_by uuid null, deleted_by uuid NULL);
   */

type Venue struct {
	Base
	ConferenceID     uuid.UUID
	Name         string
	IsActive     bool
	StartDate            time.Time
	StartDateDisplay     string
	EndDate              time.Time
	EndDateDisplay       string
	DurationDisplay      string
	Address              string
	Country string
	City string 
	County string
	PostCode string
	State string
	DisplayOrder int
	GeoLocation          GeoLocation `gorm:"embedded"`
}