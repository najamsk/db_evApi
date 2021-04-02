package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
    _"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
	_"time"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"errors"
	//"strings"
	//"fmt"
)

type TicketBooking struct {}

func (ticketBooking *TicketBooking) Insert(data models.TicketBooking) (models.TicketBooking, error) {
	db := utils.GetDb()

	err := db.Create(&data).Error
	
	if err != nil{
		return data, err
	}
	
	return data, nil
}

func (ticketBooking *TicketBooking) Update(data *models.TicketBooking) (error) {
	db := utils.GetDb()

	// If you only want to update changed Fields, you could use Update, Updates
	//var err = db.Update(&userobj).Error
	//fmt.Println("err update user:", err)
	// Save will include all fields when perform the Updating SQL, even it is not changed
	 var err = db.Save(&data).Error
	 
	return err
}

func (ticketBooking *TicketBooking) Get(id uuid.UUID) (*models.TicketBooking, error) {
	db := utils.GetDb()
	
	var tBooking models.TicketBooking
	var err = db.Find(&tBooking, "ID = ?", id).Error
	
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	
	return &tBooking, nil
}

func (ticketBooking *TicketBooking) UpdateBookingStatus(id uuid.UUID) error{
	db := utils.GetDb()
	type Result struct {ID *uuid.UUID}

	var result []Result

	err := db.Raw(" update ticket_bookings set payment_status = 'paid' where id = ? returning id", id).Scan(&result).Error;

	if(len(result)<1){
		err = errors.New("Booking status not updated.")
	}
	return err;
}
func (ticketBooking *TicketBooking) UpdateBookingItemStatus(id uuid.UUID) error{
	db := utils.GetDb()
	type Result struct {ID *uuid.UUID}

	var result []Result

	err := db.Raw(" update ticket_booking_items set payment_status = 'paid' where ticket_booking_id = ? returning id", id).Scan(&result).Error;
	if(len(result)<1){
		err = errors.New("Booking items status not updated.")
	}
	return err;
}