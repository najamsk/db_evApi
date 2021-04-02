package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	//"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"fmt"
	"github.com/satori/go.uuid"
	//"github.com/jinzhu/gorm"
	"strings"
)

type Speaker struct {}


func (speaker Speaker) GetSpeakersByConference(conferenceId uuid.UUID, userId uuid.UUID, pageSize int, offSet int, searchText string) ([]viewmodels.SpeakerVM, error) {
	db := utils.GetDb()
	
	searchText = strings.ToLower(strings.TrimSpace(searchText))

	var speakers []viewmodels.SpeakerVM

	var query string = `select users.id, users.first_name, users.last_name, users.email, 
	CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as profile_image,
						EXISTS( SELECT * FROM user_favorites where user_id = ? and conference_id=? and  entity_id= users.id and entity_type = 'speakers'
						) as Favourite,
						users.organization, users.designation, users.twitter, users.facebook, users.linked_in, cs.sort_order  
						 from users 
						 inner join conference_speakers cs on users.id = cs.user_id 
			 			left join images on users.id = images.entity_id and images.entity_type = 'user' and images.image_category = 'speaker_profile'  and images.conference_id = ? 
						 where cs.conference_id=? and users.is_active = true 
						 and ( lower(users.first_name) like '%`+searchText+`%' 
						 or lower(users.last_name) like '%`+searchText+`%' 
						 or lower(users.organization) like '%`+searchText+`%'
						 or lower(users.designation) like '%`+searchText+`%' 
						 or lower(concat(users.first_name, '', users.last_name )) like '%`+strings.Replace(searchText, " ", "", -1)+`%'  
						 or ? = '')
						 order by cs.sort_order desc, if(users.first_name = '' or users.first_name is null,1,0),users.first_name `
						 
	if(offSet>-1){
		query = query +" limit ? offset ? rows"
		db.Raw(query, userId, conferenceId, conferenceId, conferenceId, searchText,  pageSize, offSet).Scan(&speakers)
	}else{
		db.Raw(query, userId, conferenceId, conferenceId, conferenceId, searchText).Scan(&speakers)
	}
	
	
	fmt.Println(speakers)
	// if(gorm.IsRecordNotFoundError(err)){
	// 	return nil, nil
	// } else if err != nil{
	// 	return &conferencedb, err
	// }
	
	return speakers, nil
}

func (speaker Speaker) GetSpeakerProfile(conferenceId uuid.UUID, userId uuid.UUID, speakerId uuid.UUID) (viewmodels.SpeakerVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var speakerProfile viewmodels.SpeakerVM

	var query string = `select users.id, users.first_name, users.last_name, users.email,
						users.twitter, users.facebook, users.linked_in, users.youtube,
						EXISTS( SELECT * FROM user_favorites where user_id = ? and conference_id=? and  entity_id= users.id and entity_type = 'speakers'
						) as Favourite,
						users.organization, users.designation, users.twitter, users.facebook, users.linked_in,
						(select CONCAT(images.basic_url, images.image_url_prefix, '/', images.name)  from images where images.conference_id = ? and images.entity_id =users.id  and images.entity_type = 'user' and images.image_category = 'speaker_profile' limit 1) as profile_image,
						(select CONCAT(images.basic_url, images.image_url_prefix, '/', images.name)  from images where images.conference_id = ? and images.entity_id =users.id  and images.entity_type = 'user' and images.image_category = 'poster' limit 1) as poster_image,
						users.bio
						 from users 
						 inner join conference_speakers cs on users.id = cs.user_id 
			 			where cs.conference_id=? and users.id = ?;`
	//left join images on users.id = images.entity_id and images.entity_type = 'user' and images.image_category = 'user_profile' 
	db.Raw(query, userId, conferenceId, conferenceId, conferenceId, conferenceId, speakerId).Scan(&speakerProfile)
	fmt.Println(speakerProfile)
	// if(gorm.IsRecordNotFoundError(err)){
	// 	return nil, nil
	// } else if err != nil{
	// 	return &conferencedb, err
	// }
	
	return speakerProfile, nil
}

func (speaker Speaker) GetSpeakerSessions(conferenceId uuid.UUID, userId uuid.UUID, speakerId uuid.UUID) ([]viewmodels.SessionVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var speakerSessions []viewmodels.SessionVM

	var query string = `select s.id, s.title, s.summary, s.details, s.start_date, s.end_date, s.start_date_display, s.end_date_display, 
						s.duration_display, s.venue, s.seats, s.geo_location_lat, s.geo_location_long, s.radius ,
						EXISTS( SELECT * FROM user_favorites where user_id = ? and conference_id=? and  entity_id= s.id and entity_type = 'sessions'
						) as Favourite,
						s.address,
						(select count(*) FROM user_favorites uf where uf.conference_id=? and  
						uf.entity_id= s.id and uf.entity_type = 'sessions' 
						and not exists(select * from conference_speakers  cs where cs.conference_id = ? and cs.user_id = uf.user_id) 
						and exists(select * from users where users.id = uf.user_id and users.is_active = true)) as attendees_count,
						s.Capacity
						from sessions s
						inner join session_speakers ss on s.id = ss.session_id and ss.user_id = ? 
						where s.conference_id = ? order by s.sort_order desc, s.start_date ;`
	
	db.Raw(query, userId, conferenceId, conferenceId, conferenceId, speakerId, conferenceId).Scan(&speakerSessions)
	fmt.Println(speakerSessions)
	// if(gorm.IsRecordNotFoundError(err)){
	// 	return nil, nil
	// } else if err != nil{
	// 	return &conferencedb, err
	// }
	
	return speakerSessions, nil
}

func (speaker Speaker) RemoveFavoriteSpeaker(conferenceId uuid.UUID, userId uuid.UUID, speakerId uuid.UUID) (error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
		
	err := db.Exec(`delete FROM user_favorites where user_id = ? and conference_id=? and  entity_id= ? and entity_type = 'speakers'`, userId, conferenceId, speakerId).Error	
	return err
}

func (speaker Speaker) AddFavoriteSpeaker(conferenceId uuid.UUID, userId uuid.UUID, speakerId uuid.UUID) (error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
		
	err := db.Exec(`insert into user_favorites(id, user_id, entity_id, entity_type, conference_id) values( gen_random_uuid(), ?, ?, ?, ?) `, userId, conferenceId, "speakers", speakerId).Error	
	return err
}