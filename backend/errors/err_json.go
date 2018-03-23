package errors

import (
	"bytes"
	"encoding/json"

	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

var GetJsonErrMessage = func(code int, msg interface{}) []byte {
	err := models.Error{}
	err.Code = code
	switch msg.(type) {
	case string:
		err.Message = msg.(string)
	case error:
		err.Message = "unknown"
		if e, ok := msg.(error); ok {
			err.Message = e.Error()
		}
		err.Message = (msg.(error)).Error()
	default:
		err.Message = "unknown"
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(err)
	return b.Bytes()
}
