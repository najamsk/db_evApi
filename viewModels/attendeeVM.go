package viewmodels

import(
	"github.com/satori/go.uuid"
)

type AttendeeVM struct {
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
	Favourite	 bool
}