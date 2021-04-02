package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"fmt"
	"github.com/satori/go.uuid"
	"strings"
)

type Tag struct {}


func (tag Tag) GetAll() ([]models.Tag, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var tags []models.Tag
	//var err = db.Where("ID = ?", id).First(&userdb).Error
	var err = db.Find(&tags).Error
	fmt.Println(err)
	if err != nil{
		return tags, err
	}
	
	return tags, nil
}

func (tag Tag) GetUserTags(userid uuid.UUID, conferenceId uuid.UUID) ([]*models.Tag, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var tags []*models.Tag
	//var err = db.Where("ID = ?", id).First(&userdb).Error
	var query string = `select * from tags where id in (select tag_id from users_tages where user_id = ?) and conference_id = ? `
	var err = db.Raw(query, userid, conferenceId).Scan(&tags).Error
	//var err = db.Find(&tags).Error
	fmt.Println(err)
	if err != nil{
		return tags, err
	}
	
	return tags, nil
}

func (tag Tag) GetAllTagsWithSelected(userid uuid.UUID, conferenceId uuid.UUID) ([]*viewmodels.TagVM, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var tags []*viewmodels.TagVM
	//var err = db.Where("ID = ?", id).First(&userdb).Error
	var query string = `select *, EXISTS(select tag_id from users_tages where user_id = ? and tag_id = tags.id) as selected from tags where conference_id = ?`
	var err = db.Raw(query, userid, conferenceId).Scan(&tags).Error
	//var err = db.Find(&tags).Error
	fmt.Println(err)
	if err != nil{
		return tags, err
	}
	
	return tags, nil
}

//first insert data that is not exist already and then delete tags which are not in tagids list
func (tag Tag) UpdateUserTags(userid uuid.UUID, tagids []string) (error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2:",&db)
	fmt.Println(&db)
	
	var query string = `insert into users_tages (user_id, tag_id) SELECT ?, ?
						WHERE NOT EXISTS (SELECT * FROM "users_tages" WHERE "user_id" = ? AND "tag_id" = ?) `
						//fmt.Println(query)
						for i, v := range tagids {
							fmt.Println(i)
							fmt.Println(v)
					
							db.Exec(query, userid, v, userid, v)
						}

	ids := strings.Join(tagids, "','")
	deleteQuery := fmt.Sprintf(`delete from users_tages where user_id = ? and tag_id not in ('%s')`, ids)
	
	
	db.Exec(deleteQuery, userid)
	// if err != nil{
	// 	return tags, err
	// }
	
	return nil
}