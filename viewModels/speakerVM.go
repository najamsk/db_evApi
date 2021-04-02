package viewmodels

import(
	"github.com/satori/go.uuid"
)

type SpeakerVM struct {
	ID        uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	Organization string
	Designation  string
	PhoneNumber	 string
	PhoneNumber2 string
	Bio			 string
	ProfileImage string
	PosterImage string
	Favourite	 bool
	Twitter  string
	Facebook string
	LinkedIn string
	Youtube  string
	SortOrder	 int
}