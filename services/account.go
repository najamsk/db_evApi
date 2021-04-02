package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {}

func (account Account) Update(passwordCode models.ResetPassword) (error) {
	db := utils.GetDb()
	fmt.Println("passwordCode1:", passwordCode)
	var codedb models.ResetPassword
		err := db.Where("email = ?", passwordCode.Email).First(&codedb).Error; 
			
		if err != nil {
			// error handling...
			if gorm.IsRecordNotFoundError(err){
				fmt.Println("insert")
				fmt.Println("passwordCode2:", passwordCode)
				err = db.Create(&passwordCode).Error  // newUser not user
			}
			
		}else{
			fmt.Println("update")
			err = db.Model(passwordCode).Update(&passwordCode).Error
			fmt.Println("err:",err)
		}

		
	return err
}


func (account Account) GetCodeByEmail(email string) (models.ResetPassword, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var passwordCode models.ResetPassword
	
	var err = db.Where("email = ?", email).First(&passwordCode).Error
	fmt.Println(err)
	if err != nil{
		return passwordCode, err
	}
	
	return passwordCode, nil
}

func (account Account) GetPasswordTokenExpiryInMin(email string, code string) (int, error) {
	db := utils.GetDb()
	
	type Result struct {
		ExpiryMin int
	   }

	   var result Result
	
	var query string = `SELECT cast(cast((now() - updated_at) as int)/60 as int) AS expiry_min
						FROM "reset_passwords" where email = ? and code = ?`
	
	err := db.Raw(query, email, code).Scan(&result).Error
	fmt.Println("result:",result)
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return -1, nil
		}else{
		return -2, err}
	}
	
	return result.ExpiryMin, nil
}

func (account Account) UpdatePassword(email string, password string) (error) {
	db := utils.GetDb()
	var userObj models.User
	// Generate "hash" to store from user password
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
		return err	
	}
	
	err = db.Model(&userObj).Where("email = ?", email).Update("password", string(hash)).Error
	
	if err != nil{
		return err
	}
	
	return nil
}