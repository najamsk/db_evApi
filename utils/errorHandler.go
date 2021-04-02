package utils

import (
	"fmt"
	"net/http"

	viewmodels "github.com/najamsk/deutshebank/eventvisor/ev_api/viewModels"
)

func ErrorHandler(errorModel *viewmodels.Error, err error, displayError string, innerError string, source string) *viewmodels.Error {

	if err != nil {

		var loggy = FLogger{}
		loggy.OpenLog()

		if errorModel == nil {
			errorModel = new(viewmodels.Error)
		}

		fmt.Println("1")
		loggy.Logger.Info().Msg(source + " Error: " + err.Error())
		errorModel.DisplayErrors = append(errorModel.DisplayErrors, displayError)
		errorModel.InnerErrors = append(errorModel.InnerErrors, err.Error())
		errorModel.ApiStatusCode = http.StatusBadRequest

		defer loggy.CloseLog()
	}
	return errorModel
}
func AddError(errorModel *viewmodels.Error, displayError string, innerError string, source string) *viewmodels.Error {

	var loggy = FLogger{}
	loggy.OpenLog()

	if errorModel == nil {
		errorModel = new(viewmodels.Error)
	}

	loggy.Logger.Info().Msg(source + " Error: " + innerError)
	errorModel.DisplayErrors = append(errorModel.DisplayErrors, displayError)
	errorModel.InnerErrors = append(errorModel.InnerErrors, innerError)
	errorModel.ApiStatusCode = http.StatusBadRequest

	defer loggy.CloseLog()
	return errorModel
}
