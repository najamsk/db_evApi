package images

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"time"
	"github.com/satori/go.uuid"
	"github.com/gin-gonic/gin"
	//"golang.org/x/crypto/bcrypt"
	"strings"
	"github.com/najamsk/eventvisor/eventvisor.api/services"
	//"path/filepath"
	//"math/rand"
	//"strconv"
)

var apiConfig *config.Config
var apiRouter *gin.Engine
var roleService services.Role
var userService services.User
var imageService services.Image
var conferenceService services.Conference
var accountService services.Account

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	apiConfig = configuration

	apiRouter = router
	api := router.Group("api/v1/images")
	{
		
		api.POST("/upload", utils.BasicAuth(configuration), uploadImage)
		api.GET("/test", utils.BasicAuth(configuration), testGet)
	}
	apiv2 := router.Group("api/v2/images")
	{
		
		apiv2.POST("/upload", utils.BasicAuthV2(configuration), uploadImage)
		apiv2.GET("/test", utils.BasicAuthV2(configuration), testGet)
	}

}
func printStr(v string) {
	fmt.Println(v)
}

func printObj(v interface{}) {
	vBytes, _ := json.Marshal(v)
	fmt.Println(string(vBytes))
}
func testGet(c *gin.Context) {

	c.JSON(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title": "Test Get",
		},
	)
}

//accepting multipart form data
func uploadImage(c *gin.Context){

	var img *models.Image
	//var user_new models.User
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	

	var loggy = utils.FLogger{}
	loggy.OpenLog()

	var file, err = c.FormFile("ImageFile")
		if(err == nil){
			fmt.Println(file.Filename)
			//c.SaveUploadedFile(file, "uploads/profile_pic/"+userobj.ProfileImg)
			img.Name = file.Filename
			}
		fmt.Println(err)

	var zeroUUID = uuid.Nil;
	//zeroUUID = uuid.FromString("nil");
	var conferenceid=uuid.FromStringOrNil(strings.TrimSpace(c.PostForm("ConferenceID")))
	
	img.EntityID = uuid.FromStringOrNil(strings.TrimSpace(c.PostForm("EntityID")))
	img.EntityType = strings.TrimSpace(c.PostForm("EntityType"))
	img.BasicURL = strings.TrimSpace(c.PostForm("BasicURL"))
	img.ImageURLPrefix = strings.TrimSpace(c.PostForm("ImageURLPrefix"))
	img.FolderPath = strings.TrimSpace(c.PostForm("FolderPath"))
	img.ImageCategory = strings.TrimSpace(c.PostForm("ImageCategory"))
	img.ConferenceID = &conferenceid
	img.IsActive = true;
	fmt.Println("img:",img)

	if(img.Name == ""){
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid image file.")
		error.InnerErrors = append(error.InnerErrors, "image file name is missing.")
		error.ApiStatusCode = http.StatusBadRequest
	}
	if(file == nil){
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid image file.")
		error.InnerErrors = append(error.InnerErrors, "image file is missing.")
		error.ApiStatusCode = http.StatusBadRequest
	}
	if(img.EntityID == zeroUUID){
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid EntityID.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid EntityID.")
		error.ApiStatusCode = http.StatusBadRequest
	}
	if(img.ConferenceID == &zeroUUID){
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid ConferenceID.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid ConferenceID.")
		error.ApiStatusCode = http.StatusBadRequest
	}
	if(img.EntityType == ""){
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid EntityType.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid EntityType.")
		error.ApiStatusCode = http.StatusBadRequest
	}
	if(img.FolderPath == ""){
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid FolderPath.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid FolderPath.")
		error.ApiStatusCode = http.StatusBadRequest
	}
	if(img.BasicURL == ""){
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid BasicURL.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid BasicURL.")
		error.ApiStatusCode = http.StatusBadRequest
	}
	if(img.ImageURLPrefix == ""){
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid ImageURLPrefix.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid ImageURLPrefix.")
		error.ApiStatusCode = http.StatusBadRequest
	}
	if(img.ImageCategory == ""){
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid ImageCategory.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid ImageCategory.")
		error.ApiStatusCode = http.StatusBadRequest
	}
	if(error == nil){

		img, err = imageService.Save(img)
		fmt.Println(err)
		if(err == nil){
		 c.SaveUploadedFile(file, img.FolderPath + img.Name)
		}
	}
	
		
	

	if(error==nil){
		success_message = "Image uploaded successfully."}
	
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"))
	
	defer loggy.CloseLog()

	c.JSON(
		status_code,
		gin.H{
			"imageid": img.ID.String(),
			"error":error, 
			"success_message": success_message, 
		},
	)
}
