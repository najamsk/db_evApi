package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

//Constants do we really need to export them? maybe for marshal
type constants struct {
	MoftakDevEmails []string
	Environment string
	TicketExpiryInSec int
	TicketReservationReleaseTimeInSec int
	PORT      string
	BasicAuth struct {
		Username string
		Password string
	}
	RethinkDB struct {
		Address  string
		Database string
		Username string
		Password string
	}
	CockroachDB struct {
		Host  string
		Port  string
		Database string
		User string
		SSL		bool
	}
	SMTP struct {
		Host string
		Port int
		User string
		Password string
		EmailFrom string
	}
	Image struct {
		BasicURL string
		FolderPath struct {
			User struct {
				Profile string
				Poster string
			}
			Conference struct{
				Poster string
				Thumbnail string
			}
			Sponsor struct {
				Logo string
			}
			Startup struct {
				Logo string
			}
		}
		URLPrefix struct{
			User struct {
				Profile string
				Poster string
			}
			Conference struct{
				Poster string
				Thumbnail string
			}
			Sponsor struct {
				Logo string
			}
			Startup struct {
				Logo string
			}
		} 
	}
	ForgotPassword struct {
		ExpiryMinute int
	}
	Contact struct {
		Phone string
		Email string
		Web string
		Phone2 string
		WebDisplay string
	}
	TicketSeller struct {
		Roles []string

	}
	TicketChecker struct {
		Roles []string
	}
}

//Config will hold our constants as items
type Config struct {
	// Database *mgo.Database
	Items constants
	// Database *mgo.Database
}

// New func: NewConfig is used to generate a configuration instance which will be passed around the codebase
func New() (*Config, error) {
	config := Config{}
	constants, err := initViper()
	config.Items = constants
	if err != nil {
		return &config, err
	}
	// dbSession, err := mgo.Dial(config.Constants.Mongo.URL)
	// if err != nil {
	// 	return &config, err
	// }
	// config.Database = dbSession.DB(config.Constants.Mongo.DBName)
	return &config, err
}

// maybe pass config location using command line argument using flag package
// https://stackoverflow.com/questions/35419263/using-a-configuration-file-with-a-compiled-go-program
func initViper() (constants, error) {
	if os.Getenv("eventapi_env") == "staging" {
        viper.SetConfigName("settings.stagconfig")   // Configuration fileName without the .TOML or .YAML extension
	} else if os.Getenv("eventapi_env") == "production"{
		viper.SetConfigName("settings.prodconfig")
	} else {
        viper.SetConfigName("settings.config")   // Configuration fileName without the .TOML or .YAML extension
    }
	
	viper.AddConfigPath("./internal/config") // Search the root directory for the configuration file
	viper.AddConfigPath("./")                // Search the root directory for the configuration file
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		return constants{}, err
	}
	viper.WatchConfig() // Watch for changes to the configuration file and recompile
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.SetDefault("PORT", "8080")
	if err = viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file, %s", err)
	}

	var constants constants
	err = viper.Unmarshal(&constants)
	return constants, err
}
