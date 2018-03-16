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

	//can be more easily protected
	//or not using Subrouter
	pubRouter := r.Path("/pub/").Headers("Content-Type", "application/json").Subrouter()
	setupPublic(pubRouter)

	apiRouter := r.Path("/api/").Headers("Content-Type", "application/json").Subrouter()
	setupApi(apiRouter)

	return setupLogging(r)
}

func setupStatic(r *mux.Router) {
	//for pages
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}

func setupPublic(r *mux.Router) {
	r.Handle("/pub/register", handler.NotImplemented).Methods("POST")
	r.Handle("/pub/info", handler.NotImplemented).Methods("GET")
}

func setupApi(r *mux.Router) {
	r.Handle("/api/profile", handler.NotImplemented).Methods("GET", "POST")
	r.Handle("/api/buzz", handler.NotImplemented).Methods("GET", "POST")
	r.Handle("/api/jackpot", handler.NotImplemented).Methods("GET", "POST")
	r.Handle("/api/voting", handler.NotImplemented).Methods("GET", "POST")
	r.Handle("/api/leaderboard", handler.NotImplemented).Methods("GET")
	r.Handle("/api/rules", handler.NotImplemented).Methods("GET")
	r.Handle("/api/recap", handler.NotImplemented).Methods("GET")
}

func setupLogging(r http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, r)
}
