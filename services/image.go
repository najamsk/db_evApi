package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
    _"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
	_"time"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	//"strings"
	//"fmt"
)

type Image struct {}
func (img *Image) Insert(imageobj models.Image) (models.Image, error) {
	db := utils.GetDb()

    err := db.Create(&imageobj).Error
	
	if err != nil{
		return imageobj, err
	}
	
	return imageobj, nil
}

func (img *Image) Save(imageobj *models.Image) (*models.Image, error) {
	db := utils.GetDb()

	// var dbImg models.Image
	// var err = db.Where("entity_id = ? AND entity_type = ? AND image_category = ?", imageobj.EntityID, imageobj.EntityType, imageobj.ImageCategory).First(&dbImg).Error;

    // if  err != nil {
	// 	// error handling...
	// 	if gorm.IsRecordNotFoundError(err){
	// 		err = db.Save(&imageobj).Error  // newUser not user
	// 	}
	// }else{
	// 	err = db.Model(&dbImg).Where("entity_id = ? AND entity_type = ? AND image_category = ?", dbImg.EntityID, dbImg.EntityType, dbImg.ImageCategory).Update("name", imageobj.Name).Error
	// }

	err := db.Save(imageobj).Error
	
	return imageobj, err
}


func (img *Image) GetImage(entityID uuid.UUID, entityType string, imageCategory string) (*models.Image, error) {
	db := utils.GetDb()
	var imageobj models.Image
	//var err = db.Where("entity_id = ? AND entity_type = ? AND image_category = ?", entityid, "user", "user_profile").First(&imageobj).Error;
	var err = db.Where("entity_id = ? AND entity_type = ? AND image_category = ?", entityID, entityType, imageCategory).First(&imageobj).Error;
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	
	return &imageobj, nil
}
func (img *Image) GetImageByConference(entityID uuid.UUID, entityType string, imageCategory string, conferenceID uuid.UUID) (*models.Image, error) {
	db := utils.GetDb()
	var imageobj models.Image
	//var err = db.Where("entity_id = ? AND entity_type = ? AND image_category = ?", entityid, "user", "user_profile").First(&imageobj).Error;
	var err = db.Where("entity_id = ? AND entity_type = ? AND image_category = ? AND conference_id = ?", entityID, entityType, imageCategory, conferenceID).First(&imageobj).Error;
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	
	return &imageobj, nil
}

func (img *Image) GetImages(entityID uuid.UUID, entityType string, imageCategory []string) ([]*models.Image, error) {
	db := utils.GetDb()
	var imageobj []*models.Image

	var err = db.Where("entity_id = ? AND entity_type = ? and image_category in(?)", entityID, entityType, imageCategory).Find(&imageobj).Error;
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	
	return imageobj, nil
}


func (img *Image) Delete(imageId uuid.UUID) (error) {
	db := utils.GetDb()
	
	err := db.Exec(`delete FROM images where id = ? `, imageId).Error	
	return err
}