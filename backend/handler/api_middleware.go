package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/PravinPushkar/ipl18/backend/errors"
)

var IsAuthenticated = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v", r)
		token := r.Header.Get("Authorization")
		if claims, err := tokenManager.GetClaims(token); err != nil {
			log.Println("IsAuthenticated:", err.Error())
			errors.ErrWriter(w, http.StatusForbidden, "token not valid")
			return
		} else {
			if inumber, ok := claims["inumber"]; ok {
				ctx := context.WithValue(r.Context(), "inumber", inumber)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				log.Println("IsAuthenticated: inumber not found in claims")
				errors.ErrWriter(w, http.StatusForbidden, "token not valid")
				return
			}
		}
	})
}
