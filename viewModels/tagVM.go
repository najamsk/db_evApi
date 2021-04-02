package viewmodels

import(
	"github.com/satori/go.uuid"
)

type TagVM struct {
	ID			uuid.UUID
	Title    	 string
	Category     string
	SubCategory  string
	TagType     	 string
	IsActive     bool
	Selected bool
}

