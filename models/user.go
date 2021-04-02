package models

import(
	 "github.com/jinzhu/gorm"
	 "github.com/satori/go.uuid"
	//"path/filepath"
	//"fmt"
)

// User type
type User struct {
	Base
	FirstName    string
	LastName     string
	Email        string
	Password     string
	Organization string
	Designation  string
	ProfileImg   string
	PhoneNumber	 string
	PhoneNumber2 string
	Bio			 string
	SocialMedia  SocialMedia `gorm:"embedded"`
	IsActive     bool
	GeoLocation  GeoLocation `gorm:"embedded"`
	ClientID     uuid.UUID
	Conferences  []*Conference `gorm:"many2many:conferences_users;"`
	Roles        []*Role       `gorm:"many2many:users_roles;"`
	MyAgenda     MyAgenda
	Tags        []*Tag       `gorm:"many2many:users_tages;"`
	MACAddress string
	Platform string
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
		
	uuid, err := uuid.NewV4()
	if err != nil {
	 return err
	}
	//fmt.Println("user.ProfileImg: ", user.ProfileImg)
	//if profile image is not empty then set/replace ProfileImg column value with user primary key which is uuid 
	// if(user.ProfileImg != ""){
	// 	scope.SetColumn("ProfileImg", uuid.String()+ filepath.Ext(user.ProfileImg))
	// }
	return scope.SetColumn("ID", uuid)
   } 

