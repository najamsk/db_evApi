package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

//FLogger is wrapper to use with filebased logging
type FLogger struct {
	Logger  zerolog.Logger
	FileLog *os.File
}

//OpenLog will open file
func (s *FLogger) OpenLog() {
	createDirIfNotExist("logs")
	currentTime := time.Now()

	filename2 := "./logs/" + currentTime.Format("01-02-2006") + "_2log.log"
	file2, _ := os.OpenFile(filename2, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	s.Logger = zerolog.New(file2).With().Timestamp().Logger()
	s.FileLog = file2
}

//CloseLog will close file
func (s *FLogger) CloseLog() {
	defer s.FileLog.Close()
}

// CreateDirIfNotExist will create logs directory
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

// LogFile export this -- delete this method
func logFile(message string) {
	fmt.Println("oh bahi")

	var loggy FLogger
	loggy.OpenLog()
	loggy.Logger.Info().Msg(message)
	loggy.CloseLog()
	// CreateDirIfNotExist("logs")
	// currentTime := time.Now()

	// filename2 := "./logs/" + currentTime.Format("01-02-2006") + "_2log.log"
	// file2, _ := os.OpenFile(filename2, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	// defer file2.Close()

	// logger := zerolog.New(file2).With().Timestamp().Logger()
	// logger.Info().
	// 	Msg(message)

}
