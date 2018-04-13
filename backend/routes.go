package backend

import (
	"log"
	"net/http"
	"os"

	"github.wdf.sap.corp/I334816/ipl18/backend/dao"
	"github.wdf.sap.corp/I334816/ipl18/backend/handler"
	"github.wdf.sap.corp/I334816/ipl18/backend/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var PDao dao.PredictionDAO
var wsManager *service.WebSocketManager

var SetupAndGetRouter = func() http.Handler {
	log.Println("Setting up routes...")
	r := mux.NewRouter()
	setupRoutes(r)

	//wrap in route logger
	return setupLogging(r)
}

func setupRoutes(r *mux.Router) {
	setupStatic(r)
	//handle ping
	r.PathPrefix("/pub/ping").Methods("GET").Handler(handler.PingHandler)

	pubRouter := r.PathPrefix("/pub").Headers("Content-Type", "application/json").Subrouter()
	setupPublic(pubRouter)

	apiRouter := r.PathPrefix("/api").Subrouter()
	setupApi(apiRouter)
	apiRouter.Use(handler.IsAuthenticated)
	r.PathPrefix("/feeds").Handler(handler.FeedsSocketHandler{wsManager})
}

func setupStatic(r *mux.Router) {
	//for pages
	r.Handle("/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}

func setupPublic(r *mux.Router) {
	r.Handle("/register", handler.RegistrationHandler).Methods("POST")
	r.Handle("/login", handler.LoginHandler).Methods("POST")
}

func setupApi(r *mux.Router) {
	r.Handle("/users/{inumber}", handler.UserGetHandler{}).Methods("GET")
	r.Handle("/users/{inumber}", handler.UserPutHandler{}).Methods("PUT")

	r.Handle("/teams", handler.TeamsGetHandler{}).Methods("GET")
	r.Handle("/teams/{id}", handler.TeamsGetHandler{}).Methods("GET")
	r.Handle("/teams/{id}/players", handler.TeamsGetHandler{}).Methods("GET")
	r.Handle("/teams/{id}/players/{pid}", handler.TeamsGetHandler{}).Methods("GET")

	r.Handle("/players", handler.PlayersGetHandler{}).Methods("GET")
	r.Handle("/players/{id}", handler.PlayersGetHandler{}).Methods("GET")
	r.Handle("/leaders", handler.LeadersGetHandler{}).Methods("GET")

	r.Handle("/bonus", handler.BonusQuestionGetHandler{}).Methods("GET")
	r.Handle("/bonus", handler.BonusPredictionPostHandler{}).Methods("POST")

	r.Handle("/matches", handler.MatchesGetHandler{}).Methods("GET")
	r.Handle("/matches/{id}", handler.MatchesGetHandler{}).Methods("GET")
	r.Handle("/matches/{id}/stats", handler.MatchesGetHandler{}).Methods("GET")

	r.Handle("/predictions", handler.PredictionHandler{PDao}).Methods("POST")
	r.Handle("/predictions/{id}", handler.PredictionHandler{PDao}).Methods("PUT")
	r.Handle("/predictions/{id}", handler.PredictionHandler{PDao}).Methods("GET")
}

func setupLogging(r http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, r)
}

func init() {
	PDao = dao.PredictionDAO{}
	wsManager = service.NewWebSocketManager()
	wsManager.Start()
}
