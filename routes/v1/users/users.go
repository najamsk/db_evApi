package users

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
	"golang.org/x/crypto/bcrypt"
	"strings"
	"github.com/najamsk/eventvisor/eventvisor.api/services"
	"path/filepath"
	"math/rand"
	"strconv"
	"os"
)

var apiConfig *config.Config
var apiRouter *gin.Engine
var roleService services.Role
var userService services.User
var imageService services.Image
var conferenceService services.Conference
var accountService services.Account
var tagService services.Tag
var ticketTypeService services.TicketType
var emailTemplateService services.EmailTemplate

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	apiConfig = configuration

	apiRouter = router
	api := router.Group("api/v1/users")
	{
		//uncommnet when json data firstname, lastname, email, password
		api.POST("/register", utils.BasicAuth(configuration), register)

		//uncomment when multipart form data with profile image
		//api.POST("/register", utils.BasicAuth(configuration), registerUser)
		api.POST("/updateprofile", utils.BasicAuth(configuration), updateProfile)
		api.POST("/login", utils.BasicAuth(configuration), login)
		api.POST("/getprofile", utils.BasicAuth(configuration), getProfile)
		api.POST("/password/forgot", utils.BasicAuth(configuration), forgotPassword)
		api.POST("/password/reset", utils.BasicAuth(configuration), resetPassword)
		api.POST("/password/change", utils.BasicAuth(configuration), changePassword)
		api.POST("/tags/getall", utils.BasicAuth(configuration), getTags)
		api.GET("/test", utils.BasicAuth(configuration), testGet)
		
	}
	//only authorization change in v2
	apiv2 := router.Group("api/v2/users")
	{
		//uncommnet when json data firstname, lastname, email, password
		apiv2.POST("/register", utils.BasicAuthV2(configuration), register)

		//uncomment when multipart form data with profile image
		//api.POST("/register", utils.BasicAuth(configuration), registerUser)
		apiv2.POST("/updateprofile", utils.BasicAuthV2(configuration), updateProfile)
		apiv2.POST("/login", utils.BasicAuthV2(configuration), login)
		apiv2.POST("/getprofile", utils.BasicAuthV2(configuration), getProfile)
		apiv2.POST("/password/forgot", utils.BasicAuthV2(configuration), forgotPassword)
		apiv2.POST("/password/reset", utils.BasicAuthV2(configuration), resetPassword)
		apiv2.POST("/password/change", utils.BasicAuthV2(configuration), changePassword)
		apiv2.POST("/tags/getall", utils.BasicAuthV2(configuration), getTags)
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
func registerUser(c *gin.Context){

	var userobj models.User
	var user_new models.User
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var clientid uuid.UUID
	var conferenceId string//uuid.UUID
	var conference_active bool

	var loggy = utils.FLogger{}
	loggy.OpenLog()

	var file, err = c.FormFile("ProfileImg")
		if(err == nil){
			fmt.Println(file.Filename)
			//c.SaveUploadedFile(file, "uploads/profile_pic/"+userobj.ProfileImg)
			userobj.ProfileImg = file.Filename;
			}
		fmt.Println(err)

	clientid, err = uuid.FromString(strings.TrimSpace(c.PostForm("ClientID")))
	
	userobj.FirstName = strings.TrimSpace(c.PostForm("FirstName"))
	userobj.LastName = strings.TrimSpace(c.PostForm("LastName"))
	userobj.Email = strings.TrimSpace(c.PostForm("Email"))
	userobj.Password = strings.TrimSpace(c.PostForm("Password"))
	userobj.Organization = strings.TrimSpace(c.PostForm("Organization"))
	userobj.Designation = strings.TrimSpace(c.PostForm("Designation"))
	userobj.PhoneNumber = strings.TrimSpace(c.PostForm("PhoneNumber"))
	userobj.Bio = strings.TrimSpace(c.PostForm("Bio"))
	userobj.ClientID = clientid
	userobj.IsActive = true;
	
	fmt.Println(userobj)
	if(len(userobj.Email)<1){
		loggy.Logger.Info().Msg("RegisterUser:Error: Please provide valid email address.")
		//isError = true;
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid email address.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid email address.")
		error.ApiStatusCode = http.StatusBadRequest
		
	}

	if(len(userobj.Password)<6){
		loggy.Logger.Info().Msg("RegisterUser:Error: Please provide valid password.")
		//isError = true;
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Passwords must contain at least six characters.")
		error.InnerErrors = append(error.InnerErrors, "Passwords must contain at least six characters.")
	//	status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}
	
	
	if(error==nil){

	var userdb, err = userService.GetUserByEmail(strings.TrimSpace(userobj.Email));
	fmt.Println(err)
	if(len(userdb.Email)>0){
		loggy.Logger.Info().Msg("RegisterUser:Error: User with this email already exists.")
		//isError = true;
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "User with this email already exists.")
		error.InnerErrors = append(error.InnerErrors, "User with this email already exists.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}else{
		var role, err = roleService.GetRoleByName("User")
		if(err != nil){
			loggy.Logger.Info().Msg("GetRoleByName:Error:"+ err.Error())
		}else{
			userobj.Roles = []*models.Role{&role}
		}
	}

	if(error==nil){
		user_new, err = userService.Insert(userobj)
		if(err != nil){
			loggy.Logger.Info().Msg("RegisterUser:Error:"+ err.Error())
			//isError = true;
			if(error == nil){
				error = new(viewmodels.Error)
				}
			error.DisplayErrors = append(error.DisplayErrors, "Error occured while inserting user in database. please contact support team.")
			error.InnerErrors = append(error.InnerErrors, "Error occured while inserting user in database. please contact support team.")
			//status_code = http.StatusBadRequest
			error.ApiStatusCode = http.StatusBadRequest
		}else if(user_new.ProfileImg != ""){

			 var conference, err = conferenceService.GetClientActiveConference(userobj.ClientID); 
				// var conference, err = conferenceService.GetUserClientActiveConference(userdb.ClientID, userinfo.ID);
				if(conference != nil){
					conferenceId = conference.ID.String();
					conference_active = conference.IsActive
					fmt.Println("conference:",conference)

					var err = conferenceService.UpdateConferenceUser(conference.ID, user_new.ID);
					fmt.Println("err:",err)
				}else{
					loggy.Logger.Info().Msg("RegisterUser:Error: no active conference found for clientid:"+ userobj.ClientID.String())
				}

			var imageobj models.Image;
			imageobj.Name = user_new.ProfileImg;
			imageobj.BasicURL = apiConfig.Items.Image.BasicURL;
			imageobj.FolderPath = apiConfig.Items.Image.FolderPath.User.Profile;
			imageobj.ImageURLPrefix = apiConfig.Items.Image.URLPrefix.User.Profile;
			imageobj.EntityID = user_new.ID;
			imageobj.EntityType = "user";
			imageobj.ImageCategory = "user_profile";
			imageobj.IsActive = true;
			//var imageobj_new models.Image;
			imageobj, err = imageService.Insert(imageobj)
			fmt.Println(err)
			if(err == nil){
			 c.SaveUploadedFile(file, apiConfig.Items.Image.FolderPath.User.Profile + user_new.ProfileImg)
			}
		}
	}

	}
	if(error==nil){
		success_message = "You have successfully registered."}
	
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"))
	
	defer loggy.CloseLog()

	c.JSON(
		status_code,
		gin.H{
			//"is_error": isError, 
			"new_id": user_new.ID.String(),
			"error":error, 
			//"warnings": warnings, 
			//"error_message": error_message, 
			"success_message": success_message, 
			"conferenceId": conferenceId,
			"conference_active": conference_active,
		},
	)
}
//accepting json data
func register(c *gin.Context) {

	fmt.Println("register user")
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"))

	var err = error(nil);
	var userobj models.User
	var user_new models.User
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var conferenceId string//uuid.UUID
	var conference_active bool
	var userdb models.User
	var conference *models.Conference
	var role models.Role
	
	c.BindJSON(&userobj)

	var loggy = utils.FLogger{}
	loggy.OpenLog()

	//trim spaces
	userobj.Email = strings.ToLower(strings.TrimSpace(userobj.Email))
	userobj.FirstName = strings.TrimSpace(userobj.FirstName)
	userobj.LastName = strings.TrimSpace(userobj.LastName)
	userobj.Password = strings.TrimSpace(userobj.Password)
	userobj.IsActive = true;
	userobj.MACAddress = strings.TrimSpace(userobj.MACAddress);
	userobj.Platform = strings.TrimSpace(userobj.Platform);
	//userobj.Roles = 

		

	if(len(userobj.Email)<1){
		loggy.Logger.Info().Msg("RegisterUser:Error: Please provide valid email address.")
		//isError = true;
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid email address.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid email address.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}

	if(len(userobj.Password)<6){
		loggy.Logger.Info().Msg("RegisterUser:Error: Please provide valid password.")
		//isError = true;
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Passwords must be 6 or more characters without space.")
		error.InnerErrors = append(error.InnerErrors, "Passwords must contain at least six characters.")
	//	status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}
	
	
	if(error==nil){

	userdb, err = userService.GetUserByEmail(strings.TrimSpace(userobj.Email));
	fmt.Println(err)
	if(len(userdb.Email)>0){
		loggy.Logger.Info().Msg("RegisterUser:Error: User with this email already exists.")
		//isError = true;
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "User with this email already exists.")
		error.InnerErrors = append(error.InnerErrors, "User with this email already exists.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
	}
	// }else{
		
	// }

	if(error==nil){
		user_new, err = userService.Insert(userobj)
		if(err != nil){
			loggy.Logger.Info().Msg("RegisterUser:Error:"+ err.Error())
			//isError = true;
			if(error == nil){
				error = new(viewmodels.Error)
				}
			error.DisplayErrors = append(error.DisplayErrors, "Error occured while inserting user in database. please contact support team.")
			error.InnerErrors = append(error.InnerErrors, "Error occured while inserting user in database. please contact support team.")
			//status_code = http.StatusBadRequest
			error.ApiStatusCode = http.StatusBadRequest
		}else{
				role, err = roleService.GetRoleByName("User")
				if(err != nil){
					loggy.Logger.Info().Msg("GetRoleByName:Error:"+ err.Error())
				}else{
					roleService.InsertUserRole(role.ID, user_new.ID, userobj.ClientID);
					//userobj.Roles = []*models.Role{&role}
				}

				conference, err = conferenceService.GetClientActiveConference(userobj.ClientID); 
				if(err!=nil){
				loggy.Logger.Info().Msg("Register:Error:" + err.Error())}
				// var conference, err = conferenceService.GetUserClientActiveConference(userdb.ClientID, userinfo.ID);
				if(conference == nil){
					loggy.Logger.Info().Msg("Register:Error: no active conference found for clientid:"+ userobj.ClientID.String())
					
				}else{
					conferenceId = conference.ID.String();
					conference_active = conference.IsActive
					fmt.Println("conference:",conference)
					err = conferenceService.UpdateConferenceUser(conference.ID, user_new.ID);
					fmt.Println("err:",err)
				
				//send email	
				type TemplateData struct {
					FirstName string
					LastName string
					Email string
					Phone string
					SupportEmail string
					ConferenceTitle string
					WebsiteLink string
					ContactWebLink string
					Facebook string
					Twitter string }	

				contactdb, _ := conferenceService.GetContact(conference.ID);
				if(contactdb==nil){
					contactdb = &models.Conferences_contacts{};
				}

				data := TemplateData{userobj.FirstName, userobj.LastName, userobj.Email, contactdb.PhoneNumber, contactdb.EmailSupport, conference.Title, contactdb.Web, contactdb.ContactWebLink, conference.SocialMedia.Facebook, conference.SocialMedia.Twitter}
				utils.SendEmailWithDBTemplate("register", data, userobj.ClientID, conference.ID, []string{userobj.Email}, []string{}, apiConfig)
		}
				
		}
	}

	}
	if(error==nil){
		success_message = "You have successfully registered.";
	}
	
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"))
	
	defer loggy.CloseLog()

	c.JSON(
		status_code,
		gin.H{
			//"is_error": isError, 
			"new_id": user_new.ID,
			"error":error, 
			//"warnings": warnings, 
			//"error_message": error_message, 
			"success_message": success_message, 
			"conferenceId": conferenceId,
			"conference_active": conference_active,
		},
	)
}

///////

func updateProfile(c *gin.Context){//21

	var userid  uuid.UUID;
	//var conferenceId  uuid.UUID;
	var err = error(nil);
	var userdb *models.User
	//var user_new models.User
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	//var imageFileName string
	var profileImageURL string

	var loggy = utils.FLogger{}
	loggy.OpenLog()

	userid, err = uuid.FromString(strings.TrimSpace(c.PostForm("UserID")));
	if(err != nil){
		loggy.Logger.Info().Msg("updateProfile:Error: Issue while getting UserID.")
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid UserID.")
		error.InnerErrors = append(error.InnerErrors, "Issue while getting UserID as uuid.")
		error.ApiStatusCode = http.StatusBadRequest
	}

	if(error == nil){

		userdb, err = userService.Get(userid);
		if(err == nil && userdb == nil){

			loggy.Logger.Info().Msg("updateProfile:Error: Record not found in db.")
			if(error == nil){
				error = new(viewmodels.Error)
				}
			error.DisplayErrors = append(error.DisplayErrors, "User does not exist.")
			error.InnerErrors = append(error.InnerErrors, "Record not found in db.")
			error.ApiStatusCode = http.StatusBadRequest
		}else if(err != nil){
			loggy.Logger.Info().Msg("updateProfile:Error: "+ err.Error())
			if(error == nil){
				error = new(viewmodels.Error)
				}
			error.DisplayErrors = append(error.DisplayErrors, "Update profile failed. Please try again.")
			error.InnerErrors = append(error.InnerErrors, err.Error())
			error.ApiStatusCode = http.StatusBadRequest
		}else{

			userdb.Roles = nil;
			if(strings.TrimSpace(c.PostForm("FirstName")) != ""){
				userdb.FirstName = strings.TrimSpace(c.PostForm("FirstName"))}
			if(strings.TrimSpace(c.PostForm("LastName")) != ""){
				userdb.LastName = strings.TrimSpace(c.PostForm("LastName"))}

			userdb.Organization = strings.TrimSpace(c.PostForm("Organization"))
			userdb.Designation = strings.TrimSpace(c.PostForm("Designation"))
			userdb.PhoneNumber = strings.TrimSpace(c.PostForm("PhoneNumber"))
			userdb.Bio = strings.TrimSpace(c.PostForm("Bio"))
			userdb.SocialMedia.Facebook = strings.TrimSpace(c.PostForm("Facebook"))
			userdb.SocialMedia.LinkedIn = strings.TrimSpace(c.PostForm("LinkedIn"))
			userdb.SocialMedia.Twitter = strings.TrimSpace(c.PostForm("Twitter"))

				err = userService.Update(userdb)
				if(err != nil){
					loggy.Logger.Info().Msg("updateProfile:Error: "+ err.Error())
					if(error == nil){
						error = new(viewmodels.Error)
						}
					error.DisplayErrors = append(error.DisplayErrors, "Update profile failed. Please try again.")
					error.InnerErrors = append(error.InnerErrors, err.Error())
					error.ApiStatusCode = http.StatusBadRequest
				}else{

					var file, err_file = c.FormFile("ProfileImg")
					if(err_file != nil && err_file.Error() != "http: no such file"){
						loggy.Logger.Info().Msg("updateProfile: c.FromFile Error: "+ err_file.Error())
						if(error == nil){
							error = new(viewmodels.Error)
							}
						//error.DisplayErrors = append(error.DisplayErrors, "Update profile failed. Please try again.")
						error.InnerErrors = append(error.InnerErrors, err_file.Error())
						error.ApiStatusCode = http.StatusBadRequest
					}else if(file != nil){
						var image_db  = new(models.Image)
						var image_new = new(models.Image) 
						var newImg bool = false;
						image_db , err = imageService.GetImage(userid, "user", "user_profile");
						if(image_db == nil){
							newImg = true;
							image_new.Name = file.Filename;
						 }else{
							*image_new = *image_db;
							image_new.Name = image_new.ID.String() +"_"+ strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(file.Filename); 
						}
								image_new.BasicURL = apiConfig.Items.Image.BasicURL;
								image_new.FolderPath = apiConfig.Items.Image.FolderPath.User.Profile;
								image_new.ImageURLPrefix = apiConfig.Items.Image.URLPrefix.User.Profile;
								image_new.EntityID = userid;
								image_new.EntityType = "user";
								image_new.ImageCategory = "user_profile";
								image_new.IsActive = true;
								//imageobj.ConferenceID = conferenceId;
								image_new, err = imageService.Save(image_new)
						
						if(err != nil){
							loggy.Logger.Info().Msg("updateProfile:image db save Error: "+ err.Error())
							if(error == nil){
								error = new(viewmodels.Error)
								}
							//error.DisplayErrors = append(error.DisplayErrors, "Update profile failed. Please try again.")
							error.InnerErrors = append(error.InnerErrors, err.Error())
							error.ApiStatusCode = http.StatusBadRequest
						}else{
							createDirIfNotExist(apiConfig.Items.Image.FolderPath.User.Profile)
						    err = c.SaveUploadedFile(file, apiConfig.Items.Image.FolderPath.User.Profile + image_new.Name)

						 if(err != nil){ //if image file not saved to disk
							loggy.Logger.Info().Msg("updateProfile:image file save Error: "+ err.Error())
							if(error == nil){
								error = new(viewmodels.Error)
								}
							//error.DisplayErrors = append(error.DisplayErrors, "Update profile failed. Please try again.")
							error.InnerErrors = append(error.InnerErrors, err.Error())
							error.ApiStatusCode = http.StatusBadRequest
							
							//if upload image failed and image is new image then delete this image entry from database.
							//otherwise update with old image data
							if(newImg){
								err = imageService.Delete(image_new.ID)
							}else{
								image_db, err = imageService.Save(image_db)
								profileImageURL = image_db.BasicURL+ image_db.ImageURLPrefix +"/" + image_db.Name
							}
						 }else{
						 	profileImageURL = image_new.BasicURL+ image_new.ImageURLPrefix +"/" + image_new.Name
						 }
						}
					}		
		}
	}
	}
	if(error==nil){
		success_message = "Profile Updated Successfully."}
	
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"))
	
	defer loggy.CloseLog()

	c.JSON(
		status_code,
		gin.H{
			//"is_error": isError, 
			"userid": userid,
			"error":error, 
			"success_message": success_message,
			"profile_image_url": profileImageURL,
			"user": userdb,
		},
	)
}

///////

func login(c *gin.Context) {
	//var userobj models.User
	var err = error(nil);
	var role models.Role;
	var roles []*models.Role;
	var userinfo models.User
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	var conferenceId string//uuid.UUID
	var conference_active bool
	var imageurl string
	//type RequestModel struct {Email string; Password string; ClientID uuid.UUID}
	var req viewmodels.Request;
	loggy.OpenLog()
	
	c.BindJSON(&req)
	
	//fmt.Println("requestbody:",c.Request.Body)
	b, err := json.Marshal(req)
    if err == nil {
        loggy.Logger.Info().Msg("login request:"+string(b))
	}
	
	//fmt.Println("id:",userobj.Email)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	
	if(len(req.Email)<1){
		loggy.Logger.Info().Msg("login:Error: Please provide valid email address.")
		if(error == nil){
		error = new(viewmodels.Error)
		}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid email address.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid email address.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}

	if(len(req.Password)<1 || len(strings.TrimSpace(req.Password))<1){
		loggy.Logger.Info().Msg("login:Error: Please provide valid password.")
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid password.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid password.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}
	
	if(error == nil){
	var userdb, err = userService.GetUserByEmail(req.Email);
	fmt.Println(err)
	if(len(userdb.Email)<1){
		loggy.Logger.Info().Msg("login:Error: Wrong email or password. Please try again.")
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Wrong email or password. Please try again.")
		error.InnerErrors = append(error.InnerErrors, "Wrong email or password. Please try again.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}else if(!userdb.IsActive){

		loggy.Logger.Info().Msg("login:Error: This account is currently not active.")
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "This account is currently not active.")
		error.InnerErrors = append(error.InnerErrors, "User account is not active.")
		error.ApiStatusCode = http.StatusBadRequest

	}else{
		err := bcrypt.CompareHashAndPassword([]byte(userdb.Password), []byte(req.Password))
		if(err != nil){
			loggy.Logger.Info().Msg("login:Error: "+err.Error())
			if(error == nil){
				error = new(viewmodels.Error)
				}
			error.DisplayErrors = append(error.DisplayErrors, "Wrong email or password. Please try again.")
			error.InnerErrors = append(error.InnerErrors, "Password not matched.")
			//status_code = http.StatusBadRequest
			error.ApiStatusCode = http.StatusBadRequest
			
		}else{
			userinfo = userdb

			var img, err_img = imageService.GetImage(userdb.ID, "user", "user_profile");
			if(img != nil && err_img == nil){
				imageurl = img.BasicURL + img.ImageURLPrefix +"/"+ img.Name;
			}

			//get/update user conference

			var conference, err = conferenceService.GetClientActiveConference(req.ClientID); 
			// var conference, err = conferenceService.GetUserClientActiveConference(userdb.ClientID, userinfo.ID);
			if(conference != nil){
				conferenceId = conference.ID.String();
				conference_active = conference.IsActive
				fmt.Println("conference:",conference)

				var err = conferenceService.UpdateConferenceUser(conference.ID, userinfo.ID);
				fmt.Println("err:",err)
			}else{
				if(error == nil){
					error = new(viewmodels.Error)
					}
					error.DisplayErrors = append(error.DisplayErrors, "Nothing to display.")
					error.InnerErrors = append(error.InnerErrors, "no active conference found for clientid:"+req.ClientID.String())
					//status_code = http.StatusBadRequest
					error.ApiStatusCode = http.StatusBadRequest
			}
			// fmt.Println("conference.users:",conference.Users)
			 fmt.Println("err:",err)

			 roles, err = roleService.GetUserRoles(userdb.ID, req.ClientID);
			 fmt.Println("roles:",roles);
			 if(roles == nil || len(roles)<1){

				role, err = roleService.GetRoleByName("User")
				if(err != nil){
					loggy.Logger.Info().Msg("GetRoleByName:Error:"+ err.Error())
				}else{
					roleService.InsertUserRole(role.ID, userdb.ID, req.ClientID);
					//userobj.Roles = []*models.Role{&role}
				}
			 }
		}
	}
	}	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"user": userinfo,
			"conferenceId": conferenceId,
			"conference_active": conference_active,
			"profile_img_url": imageurl,
		},
	)
}

func getProfile(c *gin.Context) {
	var userinfo *models.User
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	var imageurl string
	//var taglist []models.Tag
	loggy.OpenLog()
	//type RequestModel struct {UserId uuid.UUID}
	var req viewmodels.Request;;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var userdb, err = userService.Get(req.UserID);
		
		if(err != nil){
			userinfo = nil;
			loggy.Logger.Info().Msg("getProfile:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "User does not exist.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		}else{
			userinfo = userdb

			userinfo.Tags, err = tagService.GetUserTags(req.UserID, req.ConferenceID);
			userinfo.Roles = nil;
			
			var img, err = imageService.GetImage(req.UserID, "user", "user_profile");
			if(img != nil && err == nil){
				imageurl = img.BasicURL + img.ImageURLPrefix +"/"+ img.Name;
			}
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"user": userinfo,
			"profile_img_url": imageurl,
			//"tags":taglist,
		},
	)
}

func forgotPassword(c *gin.Context) {
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	var req viewmodels.Request;
	var code string;
	loggy.OpenLog()
	
	c.BindJSON(&req)
	
	//fmt.Println("requestbody:",c.Request.Body)
	b, err := json.Marshal(req)
    if err == nil {
        loggy.Logger.Info().Msg("login request:"+string(b))
	}
	
	//fmt.Println("id:",userobj.Email)
	req.Email = strings.TrimSpace(req.Email)
	
	if(len(req.Email)<1){
		loggy.Logger.Info().Msg("forgotPassword:Error: Please provide valid email address.")
		if(error == nil){
		error = new(viewmodels.Error)
		}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid email address.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid email address.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}
	
	if(error == nil){
	var userdb, err = userService.GetUserByEmail(req.Email);
	fmt.Println(err)
	if(len(userdb.Email)<1){
		loggy.Logger.Info().Msg("forgotPassword:Error: Wrong email. Please try again.")
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Wrong email. Please try again.")
		error.InnerErrors = append(error.InnerErrors, "User with this email not found in database.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}else if(!userdb.IsActive){

		loggy.Logger.Info().Msg("forgotPassword:Error: This account is currently not active.")
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "This account is currently not active.")
		error.InnerErrors = append(error.InnerErrors, "User account is not active.")
		error.ApiStatusCode = http.StatusBadRequest

	}else{
		
		rand.Seed(int64(time.Now().Nanosecond()))
		myrand := rand.Intn(999999 - 100000) + 100000
		fmt.Println(myrand)
		code = strconv.FormatInt(int64(myrand), 10);

		var passCode models.ResetPassword;
		passCode.Email = req.Email;
		passCode.Code = code;
		//passCode.UpdatedAt = nil;
		err = accountService.Update(passCode);

		if(err == nil){
			type TemplateData struct {
				FirstName string
				LastName string
				Code string
				ConferenceTitle string
				WebsiteLink string
				ContactWebLink string
				Facebook string
				Twitter string
			}	
			var conference *models.Conference
			conference, err = conferenceService.GetClientActiveConference(req.ClientID); 
			if(err == nil && conference != nil){
				contactdb, _ := conferenceService.GetContact(conference.ID);

				if(contactdb == nil){
					contactdb = &models.Conferences_contacts{};
				}
				data := TemplateData{userdb.FirstName, userdb.LastName, code, conference.Title, contactdb.Web, contactdb.ContactWebLink,conference.SocialMedia.Facebook, conference.SocialMedia.Twitter}
				utils.SendEmailWithDBTemplate("forgot_password", data, req.ClientID, conference.ID, []string{req.Email}, []string{}, apiConfig)
				success_message = "We've sent password reset instructions to your email address.";
			}
		}else{
			loggy.Logger.Info().Msg("forgotPassword:Error: Password reset request not completed.")
			if(error == nil){
				error = new(viewmodels.Error)
				}
			error.DisplayErrors = append(error.DisplayErrors, "Password reset request not completed. Please try again")
			error.InnerErrors = append(error.InnerErrors, "Password reset request not completed. Please try again.")
			//status_code = http.StatusBadRequest
			error.ApiStatusCode = http.StatusBadRequest
		}

    	// Base 64 can be longer than len
    	//return str[:len]
	}
	}	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"code": code,
		},
	)
}

func resetPassword(c *gin.Context) {
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	var req viewmodels.Request;
	//var code string;
	loggy.OpenLog()
	
	c.BindJSON(&req)
	
	//fmt.Println("requestbody:",c.Request.Body)
	b, err := json.Marshal(req)
    if err == nil {
        loggy.Logger.Info().Msg("resetPassword request:"+string(b))
	}
	
	//fmt.Println("id:",userobj.Email)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	
	if(len(req.Email)<1){
		loggy.Logger.Info().Msg("resetPassword:Error: Please provide valid email address.")
		if(error == nil){
		error = new(viewmodels.Error)
		}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid email address.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid email address.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}
	if(len(req.Password)<6){
		loggy.Logger.Info().Msg("resetPassword:Error: Please provide valid password.")
		if(error == nil){
		error = new(viewmodels.Error)
		}
		error.DisplayErrors = append(error.DisplayErrors, "Passwords must be 6 or more characters without space.")
		error.InnerErrors = append(error.InnerErrors, "Passwords must contain at least six characters.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}
	
	if(req.Password != req.ConfirmPassword){
		loggy.Logger.Info().Msg("resetPassword:Error: Password not matched with confirm password.")
		if(error == nil){
		error = new(viewmodels.Error)
		}
		error.DisplayErrors = append(error.DisplayErrors, "Password not matched with confirm password.")
		error.InnerErrors = append(error.InnerErrors, "Password not matched with confirm password.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}
	if(len(req.PasswordToken)<1){
		loggy.Logger.Info().Msg("resetPassword:Error: Please provide valid token.")
		if(error == nil){
		error = new(viewmodels.Error)
		}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid token.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid token.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}

	if(error == nil){
		var expiryMin, err = accountService.GetPasswordTokenExpiryInMin(req.Email, req.PasswordToken); //email string, code string
		fmt.Println("expiryMin:", expiryMin)
		fmt.Println("err:", err)

		if(expiryMin == -1){
			loggy.Logger.Info().Msg("resetPassword:Error: record not found.")
			if(error == nil){
			error = new(viewmodels.Error)
			}
			error.DisplayErrors = append(error.DisplayErrors, "Record not found.")
			error.InnerErrors = append(error.InnerErrors, "Record not found.")
			error.ApiStatusCode = http.StatusBadRequest
		}else if(expiryMin == -2){
			loggy.Logger.Info().Msg("resetPassword:Error: "+err.Error())
			if(error == nil){
			error = new(viewmodels.Error)
			}
			error.DisplayErrors = append(error.DisplayErrors, "Password not updated. Please try again.")
			error.InnerErrors = append(error.InnerErrors, "Error: "+err.Error())
			error.ApiStatusCode = http.StatusBadRequest
		}else if(expiryMin > apiConfig.Items.ForgotPassword.ExpiryMinute){

			loggy.Logger.Info().Msg("resetPassword:Error: forgot password token expired.")
			if(error == nil){
			error = new(viewmodels.Error)
			}
			error.DisplayErrors = append(error.DisplayErrors, "Token expired. Please try again.")
			error.InnerErrors = append(error.InnerErrors, "Token expired. Please try again.")
			error.ApiStatusCode = http.StatusBadRequest
		}else{
			 var err = accountService.UpdatePassword(req.Email, req.Password);
			 if(err != nil){
				  if(error == nil){
					error = new(viewmodels.Error)
					}
					error.DisplayErrors = append(error.DisplayErrors, "Password not updated. Please try again.")
					error.InnerErrors = append(error.InnerErrors, "Password not updated. Please try again.")
					error.ApiStatusCode = http.StatusBadRequest
			 }else{ success_message = "Password updated successfully.";}
		}
	}
	
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			//"code": code,
		},
	)
}

func changePassword(c *gin.Context) {
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	var req viewmodels.Request;
	
	//var code string;
	loggy.OpenLog()
	
	c.BindJSON(&req)
	
	//fmt.Println("requestbody:",c.Request.Body)
	b, err := json.Marshal(req)
    if err == nil {
        loggy.Logger.Info().Msg("changePassword request:"+string(b))
	}

	
	if(req.UserID == uuid.Nil){
		loggy.Logger.Info().Msg("changePassword:Error: Please provide valid UserID.")
		if(error == nil){
		error = new(viewmodels.Error)
		}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid UserID.")
		error.InnerErrors = append(error.InnerErrors, "Error while validating UserID as uuid.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}
	if(len(req.Password)<6){
		loggy.Logger.Info().Msg("changePassword:Error: Please provide valid password.")
		if(error == nil){
		error = new(viewmodels.Error)
		}
		error.DisplayErrors = append(error.DisplayErrors, "Passwords must be 6 or more characters without space.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid password.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}
	if(req.Password != req.ConfirmPassword){
		loggy.Logger.Info().Msg("changePassword:Error: Password not matched with confirm password.")
		if(error == nil){
		error = new(viewmodels.Error)
		}
		error.DisplayErrors = append(error.DisplayErrors, "Password not matched with confirm password.")
		error.InnerErrors = append(error.InnerErrors, "Password not matched with confirm password.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}
	if(len(req.OldPassword)<1){
		loggy.Logger.Info().Msg("changePassword:Error: Please provide valid old password.")
		if(error == nil){
		error = new(viewmodels.Error)
		}
		error.DisplayErrors = append(error.DisplayErrors, "Please provide valid old password.")
		error.InnerErrors = append(error.InnerErrors, "Please provide valid old password.")
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		
	}

	if(error == nil){
		var userdb, err = userService.Get(req.UserID);
			if(userdb == nil && err == nil){
				loggy.Logger.Info().Msg("changePassword:Error: User doesn't exists.")
				if(error == nil){
				error = new(viewmodels.Error)
				}
				error.DisplayErrors = append(error.DisplayErrors, "User doesn't exists.")
				error.InnerErrors = append(error.InnerErrors, "User doesn't exists.")
				error.ApiStatusCode = http.StatusBadRequest
			}else if(err != nil){
				loggy.Logger.Info().Msg("changePassword:Error: User doesn't exists.")
				if(error == nil){
				error = new(viewmodels.Error)
				}
				error.DisplayErrors = append(error.DisplayErrors, "User doesn't exists.")
				error.InnerErrors = append(error.InnerErrors, "Error while getting user info from database.")
				error.ApiStatusCode = http.StatusBadRequest
			}else if(!userdb.IsActive){

				loggy.Logger.Info().Msg("changePassword:Error: This account is currently not active.")
				if(error == nil){
					error = new(viewmodels.Error)
					}
				error.DisplayErrors = append(error.DisplayErrors, "This account is currently not active.")
				error.InnerErrors = append(error.InnerErrors, "User account is not active.")
				error.ApiStatusCode = http.StatusBadRequest
		
			}else{
				err := bcrypt.CompareHashAndPassword([]byte(userdb.Password), []byte(req.OldPassword))
				if(err != nil){
					loggy.Logger.Info().Msg("changePassword:Error: "+err.Error())
					if(error == nil){
						error = new(viewmodels.Error)
						}
					error.DisplayErrors = append(error.DisplayErrors, "Please provide valid old password.")
					error.InnerErrors = append(error.InnerErrors, "Password not matched.")
					//status_code = http.StatusBadRequest
					error.ApiStatusCode = http.StatusBadRequest
					
				}else{
					var err = accountService.UpdatePassword(userdb.Email, req.Password);
					if(err != nil){
						 if(error == nil){
						   error = new(viewmodels.Error)
						   }
						   error.DisplayErrors = append(error.DisplayErrors, "Password not updated. Please try again.")
						   error.InnerErrors = append(error.InnerErrors, "Password not updated. Please try again.")
						   error.ApiStatusCode = http.StatusBadRequest
					}else{ success_message = "Password updated successfully.";}
	     	}
		}
	}
	
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			//"code": code,
		},
	)
}

func getTags(c *gin.Context) {
	var userinfo *models.User
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	var imageurl string
	loggy.OpenLog()
	//type RequestModel struct {UserId uuid.UUID}
	var req viewmodels.Request;;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var userdb, err = userService.Get(req.UserID);

		if(err != nil){
			userinfo = nil;
			loggy.Logger.Info().Msg("getProfile:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "User does not exist.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		}else{
			userinfo = userdb
			
			var img, err = imageService.GetImage(req.UserID, "user", "user_profile");
			if(err == nil){
				imageurl = img.BasicURL + img.ImageURLPrefix +"/"+ img.Name;
			}
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"user": userinfo,
			"profile_img_url": imageurl,
		},
	)
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}



  