package main

import (
	"fmt"
	"os"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	routes "github.com/najamsk/eventvisor/eventvisor.api/routes"
	//utils "go_cockroachdb_webapi/src/utils"
	//"github.com/joho/godotenv"
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_"net/smtp"
	_"gopkg.in/gomail.v2"

)

func main() {
	//fmt.Println(config.Env)
	fmt.Println("enviroment"+os.Getenv("ENV"))
	fmt.Println(os.Getenv("GIN_MODE"))
	fmt.Println(os.Getenv("APP_ENV"))
	configuration, err := config.New()
	if err != nil {
		fmt.Println(err)
		fmt.Println("can't reac config")
	}
	fmt.Println("configuration.Items.Environment:",configuration.Items.Environment)
	//utils.SendEmail([]string{"muhammad.habib@moftak.com"},[]string{"sayyam.ahmed@moftak.com","muhammad.tanveer@moftak.com"},"my subject2","<b>Hello everyone</b>", configuration);
	
	dbConnection := utils.OpenDbConnection(configuration)
	fmt.Println("db conn main: ", &dbConnection)
	defer dbConnection.Close()
	

	//setting up zero logger to use as middleware for gin. gin-contrib helps with it.
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	//this is same zerolog wrapped in exported func to write to a file. file open and close automatically.
	// utils.LogFile("server worked")
	var loggy = utils.FLogger{}
	loggy.OpenLog()
	loggy.Logger.Info().Msg("server worked using loggy")
	defer loggy.CloseLog()

	

	router := gin.Default()
	fmt.Println("configuration.Items.Image.URLPrefix.Sponsor.Logo: ", configuration.Items.Image.URLPrefix.Sponsor.Logo)
	fmt.Println("configuration.Items.Image.FolderPath.Sponsor.Logo: ", configuration.Items.Image.FolderPath.Sponsor.Logo)
	
	router.Static("/"+configuration.Items.Image.URLPrefix.User.Profile, configuration.Items.Image.FolderPath.User.Profile)
	router.Static("/"+configuration.Items.Image.URLPrefix.User.Poster, configuration.Items.Image.FolderPath.User.Poster)
	router.Static("/"+configuration.Items.Image.URLPrefix.Conference.Poster, configuration.Items.Image.FolderPath.Conference.Poster)
	router.Static("/"+configuration.Items.Image.URLPrefix.Conference.Thumbnail, configuration.Items.Image.FolderPath.Conference.Thumbnail)
	router.Static("/"+configuration.Items.Image.URLPrefix.Sponsor.Logo, configuration.Items.Image.FolderPath.Sponsor.Logo)
	router.Static("/"+configuration.Items.Image.URLPrefix.Startup.Logo, configuration.Items.Image.FolderPath.Startup.Logo)
	router.StaticFile("/liftpakistan19/startupfloorplan", "./uploads/floorplan/liftpk19startupfloorplan.png")
	//actually linking zeroconfg to our gin router.
	router.Use(logger.SetLogger())

	routes.Routes(configuration, router)

	router.Run(":"+string(configuration.Items.PORT))
	//zerolog.TimeFieldFormat = ""

	//Example(configuration)
}
