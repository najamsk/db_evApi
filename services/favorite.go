package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	//"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"fmt"
	"github.com/satori/go.uuid"
	//"github.com/jinzhu/gorm"
	//"strings"
)

type Favorite struct {}

func (favorite Favorite) RemoveFavorite(conferenceId uuid.UUID, userId uuid.UUID, entityId uuid.UUID, entityType string) (error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
		
	err := db.Exec(`delete FROM user_favorites where user_id = ? and conference_id=? and  entity_id= ? and entity_type = ?`, userId, conferenceId, entityId, entityType).Error	
	return err
}

func (favorite Favorite) AddFavorite(conferenceId uuid.UUID, userId uuid.UUID, entityId uuid.UUID, entityType string) (error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference

	var query string = ` insert into user_favorites (id, user_id, entity_id, entity_type, conference_id) 
						 SELECT gen_random_uuid(), ?, ?, ?, ? 
						WHERE NOT EXISTS (SELECT * FROM "user_favorites" WHERE "user_id" = ? AND "entity_id" = ? and "entity_type"= ? and "conference_id"=?) `
		
	err := db.Exec(query,userId, entityId, entityType, conferenceId, userId, entityId, entityType, conferenceId).Error	
	return err
}

func (favorite Favorite) AddFavorite_old(conferenceId uuid.UUID, userId uuid.UUID, entityId uuid.UUID, entityType string) (error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
		
	err := db.Exec(`insert into user_favorites(id, user_id, entity_id, entity_type, conference_id) 
	values( gen_random_uuid(), ?, ?, ?, ?) `, userId, entityId, entityType, conferenceId).Error	
	return err
}
func (favorite Favorite) GetFavoriteSpeakers(conferenceId uuid.UUID, userId uuid.UUID) ([]viewmodels.SpeakerVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var speakers []viewmodels.SpeakerVM
		
	var query string = `select users.id, users.first_name, users.last_name, users.email, CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as profile_image,
	users.organization, users.designation, users.twitter, users.facebook, users.linked_in 
	 from users 
	 inner join user_favorites uf on users.id = uf.entity_id and uf.entity_type = 'speakers' and uf.user_id = ? and conference_id = ?
	 inner join conference_speakers cs on users.id = cs.user_id and cs.conference_id = ?
	 left join images on users.id = images.entity_id and images.entity_type = 'user' and images.image_category = 'speaker_profile'  and images.conference_id = ? 
	  where  users.is_active = true 
	  order by cs.sort_order desc, if(users.first_name = '' or users.first_name is null,1,0), lower(users.first_name) `
	   
	  err := db.Raw(query, userId, conferenceId, conferenceId,conferenceId).Scan(&speakers).Error
	return speakers, err
}

func (favorite Favorite) GetFavoriteAttendees(conferenceId uuid.UUID, userId uuid.UUID) ([]viewmodels.AttendeeVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	
	var attendees []viewmodels.AttendeeVM
		
	var query string = ` select users.id, users.first_name, users.last_name, users.email, CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as profile_image,
	users.organization, users.designation, users.twitter, users.facebook, users.linked_in 
	 from users 
	 inner join user_favorites uf on users.id = uf.entity_id and uf.entity_type = 'attendees' and uf.user_id = ? and conference_id = ? 
	 inner join conference_attendees ca on uf.entity_id = ca.user_id and ca.conference_id = ?
	 left join images on users.id = images.entity_id and images.entity_type = 'user' and images.image_category = 'user_profile' 
	  where users.is_active = true and
	  not exists(select * from conference_speakers where conference_id = ? and user_id = users.id)
	  order by ca.sort_order desc, if(users.first_name = '' or users.first_name is null,1,0), lower(users.first_name) `	

	 err := db.Raw(query, userId, conferenceId, conferenceId, conferenceId).Scan(&attendees).Error
	return attendees, err
}

func (favorite Favorite) GetFavoriteSessions(conferenceId uuid.UUID, userId uuid.UUID) ([]viewmodels.SessionVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var sessionDB []viewmodels.SessionVM

	var query string = `
	select s.id, s.title, s.summary, s.details, s.start_date, s.end_date, s.start_date_display, s.end_date_display, 
	s.duration_display, s.venue, s.seats, s.geo_location_lat, s.geo_location_long, s.radius,
	(select avg(rating) from session_ratings where session_id = s.id) as rating , s.address
	from sessions s
	inner join user_favorites uf on s.id = uf.entity_id and uf.entity_type = 'sessions' and uf.user_id = ? and uf.conference_id = ? 
	 where s.is_active = true  
	order by s.sort_order, s.start_date`
	  
	err := db.Raw(query, userId, conferenceId).Scan(&sessionDB).Error

	return sessionDB, err
}