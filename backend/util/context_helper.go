package util

import (
	"fmt"
	"net/http"
)

var GetValueFromContext = func(r *http.Request, name string) (string, error) {
	if uname := r.Context().Value(name); uname != nil {
		if val, ok := uname.(string); ok {
			return val, nil
		}
	}
	return "", fmt.Errorf("%s not found in context", name)
}
