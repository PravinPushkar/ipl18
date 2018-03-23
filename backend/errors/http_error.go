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

var ErrWriterPanic = func(w http.ResponseWriter, code int, errActual error, errExpected error, logMsg string) {
	if errActual == nil {
		return
	}

	log.Println(fmt.Sprintf("%s-%s", logMsg, errActual.Error()))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data := GetJsonErrMessage(code, errExpected)
	w.Write(data)
	dataStr := string(data)
	panic(dataStr)
}
