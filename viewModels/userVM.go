
package viewmodels

import(
	"github.com/satori/go.uuid"
)

type UserVM struct {
	ID        uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	Password     string
	Organization string
	Designation  string
	PhoneNumber	 string
	PhoneNumber2 string
	Bio			 string
	ClientID     uuid.UUID
}