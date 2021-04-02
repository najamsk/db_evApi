package viewmodels

import(
	"github.com/satori/go.uuid"
)

type SponsorVM struct {
	ID				 uuid.UUID
	Name             string
	IsActive         bool
	SponsorLevelID   uuid.UUID
	ClientID         uuid.UUID
	ConferenceID     uuid.UUID
	SortOrder int
	Logo             string
	Twitter  string
	Facebook string
	LinkedIn string
	Youtube  string
}