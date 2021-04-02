package viewmodels

import ("time" 
"github.com/satori/go.uuid")
type SessionVM struct {
	ID				 uuid.UUID
	Title            string
	Summary          string
	Details          string
	IsActive         bool
	StartDate        time.Time
	StartDateDisplay string
	EndDate          time.Time
	EndDateDisplay   string
	DurationDisplay  string
	Address          string
	Venue            string
	Seats            int
	GeoLocationLat  float64
	GeoLocationLong float64
	Radius          float64
	Speakers	[]SpeakerVM
	Rating            float64
	Attendees	[]AttendeeVM
	Favourite	 bool
	ConferenceTitle    string
	IsFeatured		bool
	SortOrder		int
	AttendeesCount	int
	Capacity		int
}