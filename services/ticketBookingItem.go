package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
    _"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
	_"time"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type TBookingItem struct {}

func (ticketBooking *TBookingItem) Insert(data models.TicketBookingItem) (models.TicketBookingItem, error) {
	db := utils.GetDb()

	err := db.Create(&data).Error
	
	if err != nil{
		return data, err
	}
	
	return data, nil
}

func (ticketBooking *TBookingItem) Update(data *models.TicketBookingItem) (error) {
	db := utils.GetDb()

	// If you only want to update changed Fields, you could use Update, Updates
	//var err = db.Update(&userobj).Error
	//fmt.Println("err update user:", err)
	// Save will include all fields when perform the Updating SQL, even it is not changed
	 var err = db.Save(&data).Error
	 
	return err
}

func (ticketBooking *TBookingItem) Get(id uuid.UUID) (*models.TicketBookingItem, error) {
	db := utils.GetDb()
	
	var bookingItems models.TicketBookingItem
	var err = db.Find(&bookingItems, "ID = ?", id).Error
	
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	
	return &bookingItems, nil
}

func (ticketBooking *TBookingItem) GetBookingItem(booking_id uuid.UUID) ([]models.TicketBookingItem, error) {
	db := utils.GetDb()
	
	var bookingItems []models.TicketBookingItem
	var err = db.Find(&bookingItems, "ticket_booking_id = ?", booking_id).Error
	
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	
	return bookingItems, nil
}