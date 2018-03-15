package backend

import (
	"net/http"
	"os"

	"github.wdf.sap.corp/I334816/ipl18/backend/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var SetupAndGetRouter = func() http.Handler {
	r := mux.NewRouter()
	setupStatic(r)
	setupRegistered(r)
	return setupLogging(r)
}

func setupStatic(r *mux.Router) {
	//for pages
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	//for registration page
	r.Handle("/register", handler.NotImplemented).Methods("POST")
}

func setupRegistered(r *mux.Router) {
	r.Handle("/profile", handler.NotImplemented).Methods("GET", "POST")
	r.Handle("/buzz", handler.NotImplemented).Methods("GET", "POST")
	r.Handle("/jackpot", handler.NotImplemented).Methods("GET", "POST")
	r.Handle("/voting", handler.NotImplemented).Methods("GET", "POST")
	r.Handle("/leaderboard", handler.NotImplemented).Methods("GET")
	r.Handle("/rules", handler.NotImplemented).Methods("GET")
	r.Handle("/recap", handler.NotImplemented).Methods("GET")
}

func setupLogging(r http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, r)
}
