package utils

import (
	"fmt"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
)

// type Database struct {
// 	*gorm.DB
// }

var DB *gorm.DB

func migrate(db *gorm.DB) {
	var result = db.AutoMigrate(&models.User{}, &models.Role{})
	fmt.Println("db AutoMigrate err: ", result)

}

// Opening a database and save the reference to `Database` struct.
func OpenDbConnection(configuration *config.Config) *gorm.DB {

	username := configuration.Items.CockroachDB.User
	dbName := configuration.Items.CockroachDB.Database
	host := configuration.Items.CockroachDB.Host
	port := configuration.Items.CockroachDB.Port
	ssl := configuration.Items.CockroachDB.SSL
	certpath := "certs-162.252.80.136"

	//password := "habib123";
	var db *gorm.DB
	var err error
	var databaseUrl string

	// without SSL
	if ssl == false {
		databaseUrl = fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable", username, host, port, dbName)
	} else {
		databaseUrl = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=verify-full&sslrootcert=%s/ca.crt&sslcert=%s/client.root.crt&sslkey=%s/client.root.key", username, host, port, dbName, certpath, certpath, certpath)
	}

	var loggy = FLogger{}
	loggy.OpenLog()
	loggy.Logger.Info().Msg(databaseUrl)

	fmt.Println("db err: ", err)
	db, err = gorm.Open("postgres", databaseUrl)

	if err != nil {
		loggy.Logger.Info().Msg("db err:" + err.Error())
		fmt.Println("db err: ", err)
		//os.Exit(-1)
	}
	defer loggy.CloseLog()
	db.DB().SetMaxIdleConns(0)
	db.LogMode(true)
	DB = db
	//migrate(db)
	return DB
}

// Delete the database after running testing cases.
func RemoveDb(db *gorm.DB) error {
	db.Close()
	err := os.Remove(path.Join(".", "app.db"))
	return err
}

// Using this function to get a connection, you can create your connection pool here.
func GetDb() *gorm.DB {
	return DB
}
