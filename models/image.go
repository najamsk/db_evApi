package models
import(
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
   "path/filepath"
   "fmt")

	// CREATE TABLE images (id UUID PRIMARY KEY, name string, basic_url string, image_url_prefix string, folder_path string, 
	// 	entity_id uuid, entity_type string, image_category string, conference_id uuid, is_active bool, 
	// deleted_at timestamptz, created_at timestamptz, updated_at timestamptz, deleted bool, created_by uuid, updated_by uuid, deleted_by uuid);
type Image struct {
	Base
	Name		string
	BasicURL	string
	ImageURLPrefix	string
	FolderPath	string
	EntityID	uuid.UUID
	EntityType	string
	ImageCategory string
	IsActive	bool
	ConferenceID     *uuid.UUID
}

func (img *Image) BeforeCreate(scope *gorm.Scope) error {
		
	uuid, err := uuid.NewV4()
	if err != nil {
	 return err
	}
	fmt.Println("img.Name: ", img.Name)
	//if profile image is not empty then set/replace ProfileImg column value with user primary key which is uuid 
	if(img.Name != ""){
		scope.SetColumn("name", uuid.String()+ filepath.Ext(img.Name))
	}
	return scope.SetColumn("ID", uuid)
   } 