package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/satori/go.uuid"
	"github.com/gin-gonic/gin"
	"github.com/najamsk/eventvisor/eventvisor.api/services"
	//"os"
	//"io"
)

var apiConfig *config.Config
var apiRouter *gin.Engine
var roleService services.Role
var userService services.User
var imageService services.Image
var conferenceService services.Conference
var speakerService services.Speaker

// Routes will be accessable to add routes
func Routes(configuration *config.Config, router *gin.Engine) {
	apiConfig = configuration

	apiRouter = router
	api := router.Group("api/v1/conferences")
	{
		
		api.POST("/getclientactiveconference", utils.BasicAuth(configuration), getClientActiceConference)
		api.POST("/attendees", utils.BasicAuth(configuration), getConferenceAttendees)
		api.GET("/test", utils.BasicAuth(configuration), testGet)
		api.POST("/get", utils.BasicAuth(configuration), getConference)
		api.POST("/agenda", utils.BasicAuth(configuration), getConferenceAgenda)
		api.POST("/sponsors", utils.BasicAuth(configuration), getSponsors)
		api.POST("/startups", utils.BasicAuth(configuration), getStartups)
		api.POST("/startup/detail", utils.BasicAuth(configuration), getStartupById)
		api.POST("/sponsor/detail", utils.BasicAuth(configuration), getSponsorById)
		api.POST("/detail", utils.BasicAuth(configuration), getConferenceAgendaV2)
		
	}
	apiv2 := router.Group("api/v2/conferences")
	{
		
		apiv2.POST("/getclientactiveconference", utils.BasicAuthV2(configuration), getClientActiceConference)
		apiv2.POST("/attendees", utils.BasicAuthV2(configuration), getConferenceAttendees)
		apiv2.GET("/test", utils.BasicAuthV2(configuration), testGet)
		apiv2.POST("/get", utils.BasicAuthV2(configuration), getConference)
		apiv2.POST("/agenda", utils.BasicAuthV2(configuration), getConferenceAgenda)
		apiv2.POST("/sponsors", utils.BasicAuthV2(configuration), getSponsors)
		apiv2.POST("/startups", utils.BasicAuthV2(configuration), getStartups)
		apiv2.POST("/startup/detail", utils.BasicAuthV2(configuration), getStartupById)
		apiv2.POST("/sponsor/detail", utils.BasicAuthV2(configuration), getSponsorById)
		apiv2.POST("/detail", utils.BasicAuthV2(configuration), getConferenceAgendaV2)
		apiv2.POST("/getcontact", utils.BasicAuthV2(configuration), getContact)
		
	}

}


func getContact(c *gin.Context) {
  	//var err = error(nil)
	//var errorVM *viewmodels.Error
	var status_code = http.StatusOK
	var req viewmodels.Request;
	var contactdb *models.Conferences_contacts
	
	c.BindJSON(&req)

	contactdb, _ = conferenceService.GetContact(req.ConferenceID);

	if(contactdb == nil){
		contactdb = &models.Conferences_contacts{}
	}

	c.JSON(
		status_code,
		gin.H{
			"Phone2": contactdb.PhoneNumber2,
			"WebDisplay":contactdb.WebDisplay,
			"email": contactdb.Email,
			"phone": contactdb.PhoneNumber,
			"web": contactdb.Web,
		},
	)
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

func getClientActiceConference(c *gin.Context) {
	var conferenceinfo *models.Conference
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	type RequestModel struct {ClientID uuid.UUID;}
	var req RequestModel;
		
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var conferencedb, err = conferenceService.GetClientActiveConference(req.ClientID);

		if(err != nil){
			conferenceinfo = nil;
			loggy.Logger.Info().Msg("getProfile:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Conference does not exist.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		//status_code = http.StatusBadRequest
		error.ApiStatusCode = http.StatusBadRequest
		}else{
			conferenceinfo = conferencedb
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"conference": conferenceinfo,
		},
	)
}

func getConferenceAttendees(c *gin.Context) {
	//var conferenceinfo *models.Conference
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; UserId uuid.UUID}
	var req viewmodels.Request;
	var pageSize int = 50
//	var offSet int = 0
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)

	if(req.PageNo<1){
		req.OffSet = -1
	}
	
	var attendees, err = conferenceService.GetConferenceAttendees(req.ConferenceID, req.UserID, true, pageSize, req.OffSet, req.SearchText);
	fmt.Println(err)
			
	defer loggy.CloseLog()

	if(len(attendees) == pageSize && req.PageNo>0){
		req.PageNo = req.PageNo +1
	}
	
	req.OffSet = req.OffSet + len(attendees);
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"attendees": attendees,
			"next_page_no": req.PageNo,
			"attendees_count": len(attendees),
			"off_set":req.OffSet,
		},
	)
}

func getConference(c *gin.Context) {
	var conferenceinfo *models.Conference
	var errorVM *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceID uuid.UUID;}
	var req viewmodels.Request;
	//var req RequestModel;
	var sessions []viewmodels.SessionVM
	var attendies []viewmodels.AttendeeVM
	var speakers []viewmodels.SpeakerVM
	var imgPosterURL string
	var imgThumbnailURL string 
	var totalAttendies int = 0
	//var totalSpeakers int = 0
	var err = error(nil)
	//var userdb *models.User
	var conferencedb models.Conference
	var roles []*models.Role;
		
	//c.BindJSON(&userinfo) 
	c.BindJSON(&req)

		//userdb, err = userService.Get(req.UserID);
		roles, err = roleService.GetUserRoles(req.UserID, req.ClientID);
		var isTicketSeller = utils.IsInRole(roles, apiConfig.Items.TicketSeller.Roles);
		var isTicketChecker = utils.IsInRole(roles, apiConfig.Items.TicketChecker.Roles);
	
		conferencedb, err = conferenceService.Get(req.ConferenceID);

		if(err != nil){
			conferenceinfo = nil;
			loggy.Logger.Info().Msg("getConference:Error: "+err.Error())
		if(errorVM == nil){
			errorVM = new(viewmodels.Error)
			}
			errorVM.DisplayErrors = append(errorVM.DisplayErrors, "Conference does not exist.")
			errorVM.InnerErrors = append(errorVM.InnerErrors, err.Error())
			errorVM.ApiStatusCode = http.StatusBadRequest
		}else{
			conferenceinfo = &conferencedb
			sessions = conferenceService.GetConferenceSessions(req.UserID, req.ConferenceID, true, 2, 0)
			attendiesTemp, err_att := conferenceService.GetConferenceAttendees(req.ConferenceID, req.UserID, true, 10, 0, "")
			totalAttendies = conferenceService.GetConferenceAttendeesCount(req.ConferenceID, req.UserID)
			fmt.Println(err_att)
			fmt.Println(len(attendiesTemp))
			if(len(attendiesTemp)<4){
				attendies = nil
			}else{
				attendies = attendiesTemp
			}

			speakers, err = speakerService.GetSpeakersByConference(req.ConferenceID, req.UserID, 5,0, "")
						
			var img, err = imageService.GetImages(req.ConferenceID, "conference", []string{"poster","thumbnail"});
			for i, v := range img {
				fmt.Println("i:",i)
				if(v.ImageCategory=="poster"){
					imgPosterURL = v.BasicURL + v.ImageURLPrefix +"/"+ v.Name;
				}
				if(v.ImageCategory=="thumbnail"){
					imgThumbnailURL = v.BasicURL + v.ImageURLPrefix +"/"+ v.Name;
				}
				if(len(imgPosterURL)>0 && len(imgThumbnailURL)>0){
					break;
				}
			}
			// //fmt.Println("img:",img)
			 fmt.Println("err:",err)
			// var img, err = imageService.GetImage(req.ConferenceID, "conference", "poster");
			// if(err == nil){
			// 	imgPosterURL = img.BasicURL + img.ImageURLPrefix +"/"+ img.Name;
			// }

			// img, err = imageService.GetImage(req.ConferenceID, "conference", "thumbnail");
			// if(err == nil){
			// 	imgThumbnailURL = img.BasicURL + img.ImageURLPrefix +"/"+ img.Name;
			// }
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":errorVM, 
			"success_message": success_message, 
			"conference": conferenceinfo,
			"sessions": sessions,
			"imgPosterURL": imgPosterURL,
			"imgThumbnailURL": imgThumbnailURL,
			"attendees" : attendies,//[0:10], //.slice(0, 1)
			"attendees_total": totalAttendies,
			"speakers":speakers,
			//"speakers_total": totalSpeakers,
			"ticket_seller": isTicketSeller,
			"ticket_checker": isTicketChecker,
		},
	)
}

func getConferenceAgenda(c *gin.Context) {
	var conferenceAgenda []viewmodels.ConferenceAgendaVM
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	type RequestModel struct {UserID uuid.UUID; ConferenceID uuid.UUID}
	var req RequestModel;
	// var sessions []models.Session
	// var imgPosterURL string
	// var imgThumbnailURL string 
		
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var conferencedb, err = conferenceService.GetConferenceAgenda(req.UserID, req.ConferenceID);

		if(err != nil){
			conferenceAgenda = nil;
			loggy.Logger.Info().Msg("getConference:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Conference does not exist.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		error.ApiStatusCode = http.StatusBadRequest
		}else{
			conferenceAgenda = conferencedb
			// sessions, err = conferenceService.GetConferenceSessions(req.ConferenceID)

			// var img, err = imageService.GetImage(req.ConferenceID, "conference", "poster");
			// if(err == nil){
			// 	imgPosterURL = img.BasicURL + img.ImageURLPrefix +"/"+ img.Name;
			// }

			// img, err = imageService.GetImage(req.ConferenceID, "conference", "thumbnail");
			// if(err == nil){
			// 	imgThumbnailURL = img.BasicURL + img.ImageURLPrefix +"/"+ img.Name;
			// }
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"conference": conferenceAgenda,
			//"sessions": sessions,
			//"imgPosterURL": imgPosterURL,
			//"imgThumbnailURL": imgThumbnailURL,
		},
	)
}

func getConferenceAgendaV2(c *gin.Context) {
	var conferenceDays []viewmodels.ConferenceAgendaVM
	var errorModel *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	type RequestModel struct {UserID uuid.UUID; ConferenceID uuid.UUID; Date string}
	var req RequestModel;
	var title string
	var err = error(nil)
	var conferencedb models.Conference
	var session []viewmodels.SessionVM
	
	c.BindJSON(&req)
	
		conferencedb, err = conferenceService.Get(req.ConferenceID);
		if(err == nil){
			title = conferencedb.Title;
			conferenceDays, err = conferenceService.GetConferenceDays(req.UserID, req.ConferenceID);
			if(len(req.Date)>0){
				session, err = conferenceService.GetConferenceSessionsByDate(req.UserID, req.ConferenceID, req.Date);
			}else if(len(conferenceDays)>0){
				session, err = conferenceService.GetConferenceSessionsByDate(req.UserID, req.ConferenceID, conferenceDays[0].StartDate);
			}
			
		}

		if(err != nil){
			loggy.Logger.Info().Msg("getConference:Error: "+err.Error())
		if(errorModel == nil){
			errorModel = new(viewmodels.Error)
			}
			errorModel.DisplayErrors = append(errorModel.DisplayErrors, "Conference does not exist.")
			errorModel.InnerErrors = append(errorModel.InnerErrors, err.Error())
			errorModel.ApiStatusCode = http.StatusBadRequest
		}
	
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error": errorModel, 
			"success_message": success_message, 
			"title": title,
			"conference_days": conferenceDays,
			"sessions":session,
		},
	)
}

func getSponsors(c *gin.Context) {
	//var conferenceinfo *models.Conference
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; UserId uuid.UUID}
	var req viewmodels.Request;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var sponsors, err = conferenceService.GetConferenceSponsors(req.UserID, req.ConferenceID);
		fmt.Println(err)
			
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"sponsors": sponsors,
		},
	)
}

//GetConferenceStartups
func getStartups(c *gin.Context) {
	//var conferenceinfo *models.Conference
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; UserId uuid.UUID}
	var req viewmodels.Request;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var startups, err = conferenceService.GetConferenceStartups(req.ConferenceID, req.UserID);
		fmt.Println(err)
			
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"startups": startups,
		},
	)
}

func getStartupById(c *gin.Context) {
	//var conferenceinfo *models.Conference
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	//type RequestModel struct {ConferenceId uuid.UUID; UserId uuid.UUID}
	var req viewmodels.Request;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var startup, err = conferenceService.GetConferenceStartupDetail(req.StartupID, req.ConferenceID, req.UserID);
		fmt.Println(err)
		if(err != nil){
			loggy.Logger.Info().Msg("getStartupById:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Error while getting startup detail. Please try again.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		error.ApiStatusCode = http.StatusBadRequest
		}else if(startup == nil){
			if(error == nil){
				error = new(viewmodels.Error)
				}
			error.DisplayErrors = append(error.DisplayErrors, "Startup does not exists.")
			error.InnerErrors = append(error.InnerErrors, "Startup does not exists.")
			error.ApiStatusCode = http.StatusBadRequest
		}
			
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"startup": startup,
		},
	)
}

func getSponsorById(c *gin.Context) {
	//var conferenceinfo *models.Conference
	var error *viewmodels.Error
	var success_message string
	var status_code = http.StatusOK
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	type RequestModel struct {ConferenceID uuid.UUID; UserID uuid.UUID;SponsorID uuid.UUID; MACAddress string; Platform	string}
	var req RequestModel;
	
	//c.BindJSON(&userinfo)
	c.BindJSON(&req)
	
		var sponsor, err = conferenceService.GetConferenceSponsorDetail(req.SponsorID, req.ConferenceID, req.UserID);
		if(err != nil){
			loggy.Logger.Info().Msg("getSponsorById:Error: "+err.Error())
		if(error == nil){
			error = new(viewmodels.Error)
			}
		error.DisplayErrors = append(error.DisplayErrors, "Error while getting sponsor detail. Please try again.")
		error.InnerErrors = append(error.InnerErrors, err.Error())
		error.ApiStatusCode = http.StatusBadRequest
		}else if(sponsor == nil){
			if(error == nil){
				error = new(viewmodels.Error)
				}
			error.DisplayErrors = append(error.DisplayErrors, "Sponsor does not exists.")
			error.InnerErrors = append(error.InnerErrors, "Sponsor does not exists.")
			error.ApiStatusCode = http.StatusBadRequest
		}
			
	defer loggy.CloseLog()
	
	c.JSON(
		status_code,
		gin.H{
			"error":error, 
			"success_message": success_message, 
			"sponsor": sponsor,
		},
	)
}
