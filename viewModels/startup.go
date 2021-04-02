package viewmodels

import(
	"github.com/satori/go.uuid"
)

type StartupVM struct {
	ID				 uuid.UUID
	Name             string
	Description      string
	IsActive         bool
	SponsorLevelID   uuid.UUID
	ClientID         uuid.UUID
	ConferenceID     uuid.UUID
	Twitter  string
	Facebook string
	LinkedIn string
	Youtube  string
	SortOrder	 int
	Logo             string
}