package models
import (
	uuid "github.com/satori/go.uuid"
)
type Conferences_contacts struct {
	Base
	PhoneNumber  string
	PhoneNumber2 string
	Email         string
	Web           string
	WebDisplay   string
	ConferenceId uuid.UUID
	ContactWebLink   string
	EmailSupport	string
}