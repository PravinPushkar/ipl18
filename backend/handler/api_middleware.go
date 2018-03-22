package handler

import (
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

var IsAuthenticated = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if tokenManager.IsValidToken(token) != nil {
			util.ErrWriter(w, http.StatusForbidden, "token not valid")
			return
		}
		next.ServeHTTP(w, r)
	})
}
