package models

import(
	 "github.com/jinzhu/gorm"
	 "github.com/satori/go.uuid"
	//"path/filepath"
	//"fmt"
)

// CREATE TABLE user_favorites (
//     id uuid,
//     conference_id uuid NOT NULL,
//     user_id uuid NOT NULL,
//     entity_id uuid NOT null,
//     entity_type text NOT null,
//     PRIMARY KEY (id)
// );

type UserFavorite struct {
	ID          uuid.UUID `gorm:"primary_key"`
	ConferenceID      uuid.UUID `gorm:"not null"`
	UserID      uuid.UUID `gorm:"not null"`
	EntityID    uuid.UUID `gorm:"not null"`
	EntityType	string
}

// BeforeCreate will set a UUID rather than numeric ID.
func (usrfav *UserFavorite) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
	 return err
	}
	return scope.SetColumn("ID", uuid)
   }