package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
	"strings"
)

type Conference struct {}


func (conference Conference) Get(id uuid.UUID) (models.Conference, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var conferencedb models.Conference
	//var err = db.Where("ID = ?", id).First(&userdb).Error
	//.Preload("Venues")
	var err = db.Preload("Venues", func(db *gorm.DB) *gorm.DB {
		return db.Order("Venues.display_order")
	}).Find(&conferencedb, "ID = ?", id).Error
	fmt.Println(err)
	if err != nil{
		return conferencedb, err
	}
	
	return conferencedb, nil
}

func (conference Conference) GetContact(id uuid.UUID) (*models.Conferences_contacts, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var contactdb models.Conferences_contacts
	//var err = db.Where("ID = ?", id).First(&userdb).Error
	//.Preload("Venues")
	var err =db.Where("conference_id = ?", id).Find(&contactdb).Error
	fmt.Println(err)
	if(gorm.IsRecordNotFoundError(err)){
		return nil, nil
	} 
	
	return &contactdb, err
}

func (conference Conference) GetUserClientActiveConference(clientId uuid.UUID, userId uuid.UUID) (*models.Conference, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var conferencedb models.Conference
	//var err = db.Preload("Users", "ID = ? ", userId).Where("client_id = ? AND is_active = ?", clientId, true).Order("is_active, created_at desc").First(&conferencedb).Error
	var err = db.Preload("Users", "ID = ? ", userId).Where("client_id = ? ", clientId).Order("is_active desc, created_at desc").First(&conferencedb).Error
	fmt.Println(err)
	if(gorm.IsRecordNotFoundError(err)){
		return nil, nil
	} else if err != nil{
		return &conferencedb, err
	}
	
	return &conferencedb, nil
}

func (conference Conference) GetClientActiveConference(clientId uuid.UUID) (*models.Conference, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var conferencedb models.Conference
	
	var err = db.Where("client_id = ? AND is_active = ?", clientId, true).First(&conferencedb).Error

	fmt.Println(err)
	if(gorm.IsRecordNotFoundError(err)){
		return nil, nil
	} else if err != nil{
		return &conferencedb, err
	}
	
	return &conferencedb, nil
}

func (conference Conference) GetClientActiveOrLastConference(clientId uuid.UUID) (*models.Conference, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var conferencedb models.Conference
	//var err = db.Preload("Users", "ID = ? ", userId).Where("client_id = ? AND is_active = ?", clientId, true).Order("is_active, created_at desc").First(&conferencedb).Error
	var err = db.Where("client_id = ? ", clientId).Order("is_active desc, created_at desc").First(&conferencedb).Error
	fmt.Println(err)
	if(gorm.IsRecordNotFoundError(err)){
		return nil, nil
	} else if err != nil{
		return &conferencedb, err
	}
	
	return &conferencedb, nil
}

func (conference Conference) UpdateConferenceUser(conferenceId uuid.UUID, userId uuid.UUID) (error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var conferencedb models.Conference
	var err = db.Where("id = ?", conferenceId).First(&conferencedb).Error
	
	// var err = db.Where("client_id = ? AND is_active = ?", clientId, true).First(&conferencedb).Error
	 db.Model(&conferencedb).Association("Users").Append(models.User{Base: models.Base{ID: userId}})
	// db.Model(&conferencedb).Association("Attendees").Append(models.User{Base: models.Base{ID: userId}})

	var query string = `INSERT INTO "conference_attendees" ("conference_id","user_id") SELECT ?, ?  
	WHERE NOT EXISTS (SELECT * FROM "conference_attendees" WHERE "conference_id" = ? AND "user_id" = ?)
	AND NOT EXISTS (SELECT * FROM "conference_speakers" WHERE "conference_id" = ? AND "user_id" = ?);`
	 
	err = db.Exec(query, conferenceId, userId, conferenceId, userId, conferenceId, userId).Error
	
	
	return err
}

func (conference Conference) GetConferenceAttendees(conferenceId uuid.UUID, userId uuid.UUID, isImageMandatory bool, pageSize int, offSet int, searchText string) ([]viewmodels.AttendeeVM, error) {
	db := utils.GetDb()
	searchText = strings.ToLower(strings.TrimSpace(searchText))
	fmt.Println(&db)
	//var conferencedb models.Conference
	var attendees []viewmodels.AttendeeVM
	//var limitQuery string = ` limit 10 `
	var query string = `select users.id, users.first_name, users.last_name, users.email, CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as profile_image,
						EXISTS( SELECT * FROM user_favorites where user_id = ? and conference_id=? and  entity_id= users.id and entity_type = 'attendees'
						) as Favourite,
						users.organization, users.designation, users.twitter, users.facebook, users.linked_in
						 from users 
						 inner join conference_attendees ca on users.id = ca.user_id 
			 			left join images on users.id = images.entity_id and images.entity_type = 'user' and images.image_category = 'user_profile' 
						 where ca.conference_id=? and users.is_active = true and (images.id is not null or ? = false) 
						 and not exists(select * from conference_speakers where conference_id = ? and user_id = users.id) 
						 and (trim(users.organization) !='' or trim(users.designation) !='')
						 and ( lower(users.first_name) like '%`+searchText+`%' 
						 or lower(users.last_name) like '%`+searchText+`%' 
						 or lower(users.organization) like '%`+searchText+`%'
						 or lower(users.designation) like '%`+searchText+`%'
						 or lower(concat(first_name, '', last_name )) like '%`+strings.Replace(searchText, " ", "", -1)+`%' 
						 or ? = '')
						 order by ca.sort_order desc, if(users.first_name = '' or users.first_name is null,1,0), 
						 lower(users.first_name) `
	
	// if(limit>0){
	// 	query += limitQuery}
	if(offSet>-1){
		query = query +" limit ? offset ? rows"
		//db.Raw(query, userId, conferenceId, conferenceId, pageSize, offSet).Scan(&speakers)
		db.Raw(query, userId, conferenceId, conferenceId, isImageMandatory, conferenceId, searchText, pageSize, offSet).Scan(&attendees)
	}else{
		db.Raw(query, userId, conferenceId, conferenceId, isImageMandatory, conferenceId, searchText).Scan(&attendees)
	}

	
	fmt.Println(attendees)	
	return attendees, nil
}

func (conference Conference) GetConferenceAttendeesCount(conferenceId uuid.UUID, userId uuid.UUID) (int) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	type Result struct {
		TotalAttendees int
	   }

	   var result Result
	
	var query string = `select count(*) as total_attendees 
	from conference_attendees ca 
	inner join users u on ca.user_id = u.id 
	where ca.conference_id=? and u.is_active = true 
	and not exists(select * from conference_speakers where conference_id = ? and user_id = u.id) ;`
	
	db.Raw(query, conferenceId, conferenceId).Scan(&result)
	fmt.Println(result)	
	return result.TotalAttendees
}

func (conference Conference) GetConferenceSessions(userId uuid.UUID, conferenceId uuid.UUID, onlyFeatured bool, pageSize int, offSet int) ([]viewmodels.SessionVM) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var sessions []viewmodels.SessionVM
	//var err = db.Where("ID = ?", id).First(&userdb).Error
	//var err = db.Find(&sessions, "conference_id = ?", conference_id).Error
	sessionQuery := `select *, EXISTS( select * FROM user_favorites where user_id = ? and conference_id=? 
						and  entity_id= sessions.id and entity_type = 'sessions'
						) as Favourite,

						(select count(*) FROM user_favorites uf where uf.conference_id=? and  
						uf.entity_id= sessions.id and uf.entity_type = 'sessions' 
						and not exists(select * from conference_speakers  cs where cs.conference_id = ? and cs.user_id = uf.user_id) 
						and exists (select * from users where users.id = uf.user_id and users.is_active = true)) as attendees_count,
						sessions.Capacity

						from sessions WHERE sessions.conference_id = ? and sessions.is_active = true and (sessions.is_featured = true or ? = false) order by sort_order desc `
		
	if(offSet>-1){
					sessionQuery = sessionQuery +" limit ? offset ? rows"
					db.Raw(sessionQuery, userId, conferenceId, conferenceId, conferenceId, conferenceId, onlyFeatured,  pageSize, offSet).Scan(&sessions)
				}else{
						db.Raw(sessionQuery, userId, conferenceId, conferenceId, conferenceId, conferenceId, onlyFeatured).Scan(&sessions)
					}

	

	// fmt.Println(err)
	// if err != nil{
	// 	return sessions, err
	// }
	
	return sessions
}

func (conference Conference) GetConferenceAgenda(userId uuid.UUID, conferenceId uuid.UUID) ([]viewmodels.ConferenceAgendaVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var agendaList []viewmodels.ConferenceAgendaVM
	var session []viewmodels.SessionVM
	var speakers []viewmodels.SpeakerVM


	var query string = `select distinct cast( cast(sessions.start_date as date) as string) as start_date, 
						conferences.title, conferences.twitter, conferences.facebook, conferences.youtube, conferences.linked_in 
						from "sessions" 
						inner join  conferences on sessions.conference_id = conferences.id and conferences.id = ?
						where conference_id = ? and sessions.is_active = true 
						order by start_date`
	
	db.Raw(query, conferenceId, conferenceId).Scan(&agendaList)
	fmt.Println(agendaList)

	for i, v := range agendaList {
		fmt.Println(i)
		fmt.Println(v.StartDate)

		sessionQuery := `select *, EXISTS( select * FROM user_favorites where user_id = ? and conference_id=? and  entity_id= sessions.id and entity_type = 'sessions'
						) as Favourite,

						(select count(*) FROM user_favorites uf where uf.conference_id=? and  
						uf.entity_id= sessions.id and uf.entity_type = 'sessions' 
						and not exists(select * from conference_speakers  cs where cs.conference_id = ? and cs.user_id = uf.user_id) 
						and exists (select * from users where users.id = uf.user_id and users.is_active = true)) as attendees_count,
						sessions.Capacity
						from sessions WHERE cast(start_date as date) = ? and sessions.is_active = true 
						
						order by sessions.sort_order desc, sessions.start_date`
		db.Raw(sessionQuery, userId, conferenceId, conferenceId, conferenceId, v.StartDate).Scan(&session)
		agendaList[i].Sessions = session;

		for j, v := range session {

			speakersQuery := `select users.id, users.first_name, users.last_name, users.email, CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as profile_image
								from users 
								inner join session_speakers ss on users.id = ss.user_id 
								left join images on users.id = images.entity_id and images.entity_type = 'user' and images.image_category = 'speaker_profile' and images.conference_id=?
								where SS.session_id=? and users.is_active = true;`
			db.Raw(speakersQuery,conferenceId, v.ID).Scan(&speakers)

			agendaList[i].Sessions[j].Speakers = speakers;
		}
		
	}
	
	return agendaList, nil
}

func (conference Conference) GetConferenceDays(userId uuid.UUID, conferenceId uuid.UUID) ([]viewmodels.ConferenceAgendaVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	
	var agendaList []viewmodels.ConferenceAgendaVM
	
    var query string = `select distinct cast( cast(sessions.start_date as date) as string) as start_date
						from "sessions" 
						where conference_id = ? and sessions.is_active = true 
						order by start_date`
	
	db.Raw(query, conferenceId).Scan(&agendaList)
	fmt.Println(agendaList)

	return agendaList, nil
}
func (conference Conference) GetConferenceSessionsByDate(userId uuid.UUID, conferenceId uuid.UUID, date string) ([]viewmodels.SessionVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	
	var session []viewmodels.SessionVM
	var speakers []viewmodels.SpeakerVM

	sessionQuery := `select *, EXISTS( select * FROM user_favorites where user_id = ? and conference_id=? and  entity_id= sessions.id and entity_type = 'sessions'
						) as Favourite,

						(select count(*) FROM user_favorites uf where uf.conference_id=? and  
						uf.entity_id= sessions.id and uf.entity_type = 'sessions' 
						and not exists(select * from conference_speakers  cs where cs.conference_id = ? and cs.user_id = uf.user_id) 
						and exists (select * from users where users.id = uf.user_id and users.is_active = true)) as attendees_count,
						sessions.Capacity
						from sessions WHERE cast(start_date as date) = ? and sessions.is_active = true 
						
						order by sessions.sort_order desc, sessions.start_date`
		db.Raw(sessionQuery, userId, conferenceId, conferenceId, conferenceId, date).Scan(&session)


		for j, v := range session {

			speakersQuery := `select users.id, users.first_name, users.last_name, users.email, CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as profile_image
								from users 
								inner join session_speakers ss on users.id = ss.user_id 
								left join images on users.id = images.entity_id and images.entity_type = 'user' and images.image_category = 'user_profile' 
								where SS.session_id=? and users.is_active = true;`
			db.Raw(speakersQuery, v.ID).Scan(&speakers)

			session[j].Speakers = speakers;
		}	
	return session, nil
}

func (conference Conference) GetConferenceSpeakers(conferenceId uuid.UUID, pageSize int, offSet int) ([]viewmodels.SpeakerVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var speakers []viewmodels.SpeakerVM

	var query string = `select users.id, users.first_name, users.last_name, users.email, CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as profile_image,
						users.organization, users.designation, users.twitter, users.facebook, users.linked_in
						 from users 
						 inner join conference_speakers cs on users.id = cs.user_id 
			 			left join images on users.id = images.entity_id and images.entity_type = 'user' and images.image_category = 'user_profile' 
						 where cs.conference_id=? and users.is_active = true 
						 order by cs.sort_order desc, if(users.first_name = '' or users.first_name is null,1,0), lower(users.first_name);`
	
	if(offSet>-1){
					query = query +" limit ? offset ? rows"
					db.Raw(query, conferenceId, pageSize, offSet).Scan(&speakers)
			}else{
					db.Raw(query, conferenceId).Scan(&speakers)
				}
			
	return speakers, nil
}

func (conference Conference) GetConferenceSponsors(userId uuid.UUID, conferenceId uuid.UUID) ([]viewmodels.SponsorLevelVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var sponserLevels []viewmodels.SponsorLevelVM
	var sponsors []viewmodels.SponsorVM
	//var speakers []viewmodels.SpeakerVM

	var query string = `select id, name , client_id , conference_id , is_active, sort_order 
						from "sponsor_levels"  where conference_id = ? and is_active = true order by sort_order desc`
	
	db.Raw(query, conferenceId).Scan(&sponserLevels)
	fmt.Println(sponserLevels)

	for i, v := range sponserLevels {
		fmt.Println(i)
		fmt.Println(v.ID)

		sponsorsQuery := `select sponsors.id, sponsors.name , sponsors.client_id , sponsors.conference_id , sponsors.sponsor_level_id,  sponsors.is_active, sponsors.sort_order,
						  CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as logo 
						  from "sponsors" 
						  left join images on sponsors.id = images.entity_id and images.entity_type = 'sponsor' and images.image_category = 'logo' and images.conference_id = ?  
						  where sponsor_level_id = ? and sponsors.is_active = true order by sort_order desc`
		db.Raw(sponsorsQuery, conferenceId, v.ID).Scan(&sponsors)
		sponserLevels[i].Sponsors = sponsors;	
	}
	
	return sponserLevels, nil
}

func (conference Conference) GetConferenceStartups(conferenceId uuid.UUID, userId uuid.UUID) ([]viewmodels.StartupVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var startups []viewmodels.StartupVM
	//var sponsors []viewmodels.SponsorVM
	//var speakers []viewmodels.SpeakerVM

	var query string = `select startups.id, startups.name , startups.client_id , startups.conference_id , startups.is_active, startups.sort_order,
						startups.twitter, startups.facebook, startups.linked_in, startups.youtube, 
						CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as logo 
						from "startups"  
						left join images on startups.id = images.entity_id and images.entity_type = 'startup' and images.image_category = 'logo' and images.conference_id = ?  
						where startups.conference_id = ? and startups.is_active = true order by startups.sort_order desc`
	
	db.Raw(query, conferenceId, conferenceId).Scan(&startups)
	fmt.Println(startups)


	return startups, nil
}

func (conference Conference) GetConferenceStartupDetail(startupId uuid.UUID, conferenceId uuid.UUID, userId uuid.UUID) (*viewmodels.StartupVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var startup viewmodels.StartupVM
	//var sponsors []viewmodels.SponsorVM
	//var speakers []viewmodels.SpeakerVM

	var query string = `select startups.id, startups.name , startups.client_id , startups.conference_id , startups.is_active, startups.sort_order,
						startups.twitter, startups.facebook, startups.linked_in, startups.youtube, 
						CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as logo, startups.description 
						from "startups"  
						left join images on startups.id = images.entity_id and images.entity_type = 'startup' and images.image_category = 'logo' and images.conference_id = ?  
						where startups.conference_id = ? and startups.id = ? and startups.is_active = true order by startups.sort_order desc`
	
	err := db.Raw(query, conferenceId, conferenceId, startupId).Scan(&startup).Error
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	fmt.Println(startup)


	return &startup, nil
}

func (conference Conference) GetConferenceSponsorDetail(startupId uuid.UUID, conferenceId uuid.UUID, userId uuid.UUID) (*viewmodels.StartupVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var startup viewmodels.StartupVM
	//var sponsors []viewmodels.SponsorVM
	//var speakers []viewmodels.SpeakerVM

	var query string = `select sponsors.id, sponsors.name , sponsors.client_id , sponsors.conference_id , sponsors.is_active, sponsors.sort_order,
						sponsors.twitter, sponsors.facebook, sponsors.linked_in, sponsors.youtube, 
						CONCAT(images.basic_url, images.image_url_prefix, '/', images.name) as logo, sponsors.description 
						from "sponsors"  
						left join images on sponsors.id = images.entity_id and images.entity_type = 'sponsor' and images.image_category = 'logo' and images.conference_id = ?  
						where sponsors.conference_id = ? and sponsors.id = ? and sponsors.is_active = true  order by sponsors.sort_order desc`
	
	err := db.Raw(query, conferenceId, conferenceId, startupId).Scan(&startup).Error
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
	fmt.Println(startup)


	return &startup, nil
}


