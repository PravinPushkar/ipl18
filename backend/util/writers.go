package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var OkWriter = func(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

var StructWriter = func(w http.ResponseWriter, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(msg)
	w.Write(b.Bytes())
}
