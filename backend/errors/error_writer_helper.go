package errors

import (
	"fmt"
	"log"
	"net/http"
)

var ErrWriter = func(w http.ResponseWriter, code int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(GetJsonErrMessage(code, msg))
}

//helper for writing error messages
//does not do anything if no error has been found
//errActual is the error which has/has not happened
//if errActual is not null., the following happens:
//the errorCode code is sent to the user
//errExpected is the one which needs to be sent to the user (could be simplified version of errActual)
//if it is nil, errActual will be sent
//logMsg is a string which is logged to the console
var ErrWriterPanic = func(w http.ResponseWriter, code int, errActual error, errExpected error, logMsg string) {
	if errActual == nil {
		return
	}

	log.Println(fmt.Sprintf("%s-%s", logMsg, errActual.Error()))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data := []byte{}
	if errExpected == nil {
		data = GetJsonErrMessage(code, errActual)
	} else {
		data = GetJsonErrMessage(code, errExpected)
	}
	w.Write(data)
	panic(string(data))
}

//checks a few categories
var ErrAnalyzePanic = func(w http.ResponseWriter, err error, logMsg string) {
	if err == nil {
		return
	}
	switch err.(type) {
	case error:
		ErrWriterPanic(w, http.StatusInternalServerError, err, err, logMsg)
	case *DaoError:
		de := err.(*DaoError)
		ErrWriterPanic(w, de.Code, de.ActualErr, de.UserErr, logMsg)
	case *Error:
		de := err.(*Error)
		ErrWriterPanic(w, de.Code, err, err, logMsg)
	}
}
