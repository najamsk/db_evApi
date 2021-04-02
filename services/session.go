package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
	//"strings"
)

type Session struct {}

func (session Session) Get(id uuid.UUID) (*models.Session, error) {
	db := utils.GetDb()
	
	var sessionDB models.Session
	var err = db.Find(&sessionDB, "ID = ?", id).Error
	
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	
	return &sessionDB, nil
}
func (session Session) GetSession(id uuid.UUID, conferenceId uuid.UUID, userId uuid.UUID) (*viewmodels.SessionVM, error) {
	db := utils.GetDb()
	
	var sessionDB viewmodels.SessionVM
	var speakers []viewmodels.SpeakerVM
	var attendees []viewmodels.AttendeeVM
	var query string = `select s.id, s.title, s.summary, s.details, s.start_date, s.end_date, s.start_date_display, s.end_date_display, 
						s.duration_display, s.venue, s.seats, s.geo_location_lat, s.geo_location_long, s.radius,
						(select avg(rating) from session_ratings where session_id = s.id) as rating,
						EXISTS( SELECT * FROM user_favorites where user_id = ? and conference_id=? and  entity_id= s.id and entity_type = 'sessions'
						) as Favourite ,
						s.address,
						(select title from conferences where id = ? limit 1) as conference_title,

						(select count(*) FROM user_favorites uf where uf.conference_id=? and  
						uf.entity_id= s.id and uf.entity_type = 'sessions' 
						and not exists(select * from conference_speakers  cs where cs.conference_id = ? and cs.user_id = uf.user_id) 
						and exists(select * from users where users.id = uf.user_id and users.is_active = true)) as attendees_count,
						s.Capacity

						from sessions s 
						where s.id = ? ;`
	
	err := db.Raw(query, userId, conferenceId,conferenceId, conferenceId, conferenceId, id).Scan(&sessionDB).Error
	
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}else {

		speakersQuery := `select users.id, users.first_name, users.last_name, users.email, CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as profile_image,
		EXISTS( SELECT * FROM user_favorites where user_id = ? and conference_id=? and  entity_id= users.id and entity_type = 'speakers'
		) as Favourite, users.organization, users.designation, ss.sort_order  
								from users 
								inner join session_speakers ss on users.id = ss.user_id 
								left join images on users.id = images.entity_id and images.entity_type = 'user' and images.image_category = 'speaker_profile' 
								and images.conference_id=?
								where SS.session_id=? and users.is_active = true order by ss.sort_order desc;`
			db.Raw(speakersQuery, userId, conferenceId, conferenceId,sessionDB.ID).Scan(&speakers)

			sessionDB.Speakers = speakers;


			attendeesQuery := `select users.id, users.first_name, users.last_name, users.email, CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as profile_image,
			EXISTS( SELECT * FROM user_favorites where user_id = ? and conference_id=? and  entity_id= users.id and entity_type = 'attendees'
			) as Favourite, users.organization, users.designation
			from users 
			inner join user_favorites uf on users.id = uf.user_id and uf.entity_type = 'sessions' and uf.entity_id = ?
			left join images on users.id = images.entity_id and images.entity_type = 'user' and images.image_category = 'user_profile' 
			 where users.is_active = true 
			 and not exists(select * from conference_speakers where conference_id = ? and user_id = users.id)  order by users.first_name;`
			db.Raw(attendeesQuery, userId, conferenceId, sessionDB.ID, conferenceId).Scan(&attendees)

			sessionDB.Attendees = attendees;
	}
	
	return &sessionDB, nil
}

func (session Session) UpdateRating(userid uuid.UUID, sessionid uuid.UUID, rating float64) (error) {
	db := utils.GetDb()

	var sessionRat models.SessionRating
	var err = db.Where("user_id = ? AND session_id = ? ", userid, sessionid).First(&sessionRat).Error;

    if  err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err){
			err = db.Create(&models.SessionRating{UserID: userid, SessionID: sessionid, Rating: rating}).Error  // newUser not user
		}
	}else{
		err = db.Model(&sessionRat).Where("user_id = ? AND session_id = ? ", userid, sessionid).Update("rating", rating).Error
	}
	
	return err
}

func (session Session) GetSessionAvgRating(sessionid uuid.UUID, userId uuid.UUID) (float64) {
	db := utils.GetDb()
	
	//fmt.Println(&db)
	type Result struct {
		AvgRating float64
	   }

	   var result Result
	
	var query string = `select avg(rating) as avg_rating from session_ratings where session_id = ?`
	
	db.Raw(query, sessionid).Scan(&result)
	fmt.Println(result)	
	return result.AvgRating
}

func (session Session) CapacityReached(sessionid uuid.UUID, userId uuid.UUID) (bool, error) {
	db := utils.GetDb()
	
	//fmt.Println(&db)
	type Result struct {
		CapacityReached int
	   }

	   var result Result
	
	var query string = `select 1 as capacity_reached from sessions 
						where Capacity > 
						(select count(1) from user_favorites 
							where entity_type= 'sessions' and entity_id = sessions.id
							and exists (select * from users where users.id = user_favorites.user_id and users.is_active = true) ) 
						and id = ?`
	
	err := db.Raw(query, sessionid).Scan(&result).Error
	fmt.Println("1:")
	if  err != nil {
		if gorm.IsRecordNotFoundError(err){
			return true, nil
		}
	}
	fmt.Println("2:")
	fmt.Println("result:", result)	
	fmt.Println("err:", err)	
	return result.CapacityReached != 1, err
}








