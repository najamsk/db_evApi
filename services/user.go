package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"fmt"
    "github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
	_"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/satori/go.uuid"
)

type User struct {}

func (user User) Insert(userobj models.User) (models.User, error) {
	db := utils.GetDb()

	// Generate "hash" to store from user password
    hash, err := bcrypt.GenerateFromPassword([]byte(userobj.Password), bcrypt.DefaultCost)
    if err != nil {
		return userobj, err	
	}

	userobj.Password = string(hash)
	
	err = db.Create(&userobj).Error
	
	if err != nil{
		return userobj, err
	}
	
	return userobj, nil
}

func (user User) Update(userobj *models.User) (error) {
	db := utils.GetDb()

	// If you only want to update changed Fields, you could use Update, Updates
	//var err = db.Update(&userobj).Error
	//fmt.Println("err update user:", err)
	// Save will include all fields when perform the Updating SQL, even it is not changed
	 var err = db.Save(&userobj).Error
	 fmt.Println("err save user:", err)
	
	return err
}

func (user User) Get(id uuid.UUID) (*models.User, error) {
	db := utils.GetDb()
	
	var userdb models.User
	var err = db.Preload("Roles").Find(&userdb, "ID = ?", id).Error
	
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	
	return &userdb, nil
}

func (user User) GetUserByEmail(email string) (models.User, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var userdb models.User
	//clientId, err_uuid := uuid.FromString("8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e")query
	//fmt.Println(err_uuid)
	//var err = db.Where("email = ?", email).First(&userdb).Where("client_id = ?", clientId).Related(&userdb.Conferences,  "Conferences").Error
	//var err = db.Preload("Conferences", " is_active = ? AND client_id", true, clientid).Where("email = ?", email).First(&userdb).Error
	var err = db.Where("lower(email) = lower(?)", email).First(&userdb).Error
	fmt.Println(err)
	if err != nil{
		return userdb, err
	}
	
	return userdb, nil
}

