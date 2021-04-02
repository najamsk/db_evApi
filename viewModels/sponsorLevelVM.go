package viewmodels

import(
	"github.com/satori/go.uuid"
)

type SponsorLevelVM struct {
	ID				 uuid.UUID
	Name             string
	ClientID         uuid.UUID
	ConferenceID     uuid.UUID
	SortOrder int
	Sponsors         []SponsorVM
}