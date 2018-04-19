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

var (
	pDao      dao.PredictionDAO
	uDao      dao.UserDAO
	tDao      dao.TeamDAO
	playerDao dao.PlayerDAO
	bDao      dao.BonusDAO
	mDao      dao.MatchesDAO
	wsManager *service.WebSocketManager
)

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
	r.Handle("/register", handler.PublicUserHandler{uDao}).Methods("POST")
	r.Handle("/login", handler.PublicUserHandler{uDao}).Methods("POST")
}

func setupApi(r *mux.Router) {
	r.Handle("/users/{inumber}", handler.UserGetHandler{}).Methods("GET")
	r.Handle("/users/{inumber}", handler.UserPutHandler{}).Methods("PUT")

	r.Handle("/teams", handler.TeamsGetHandler{playerDao, tDao}).Methods("GET")
	r.Handle("/teams/{id}", handler.TeamsGetHandler{playerDao, tDao}).Methods("GET")
	r.Handle("/teams/{id}/players", handler.TeamsGetHandler{playerDao, tDao}).Methods("GET")
	r.Handle("/teams/{id}/players/{pid}", handler.TeamsGetHandler{playerDao, tDao}).Methods("GET")

	r.Handle("/players", handler.PlayersGetHandler{playerDao}).Methods("GET")
	r.Handle("/players/{id}", handler.PlayersGetHandler{playerDao}).Methods("GET")

	r.Handle("/leaders", handler.LeadersGetHandler{uDao}).Methods("GET")

	r.Handle("/bonus", handler.BonusQuestionGetHandler{bDao}).Methods("GET")
	r.Handle("/bonus", handler.BonusPredictionPostHandler{bDao}).Methods("POST")

	r.Handle("/matches", handler.MatchesGetHandler{mDao}).Methods("GET")
	r.Handle("/matches/userStats", handler.PredictionHandler{pDao}).Methods("GET")
	r.Handle("/matches/{id}", handler.MatchesGetHandler{mDao}).Methods("GET")
	r.Handle("/matches/{id}/stats", handler.MatchesGetHandler{mDao}).Methods("GET")

	r.Handle("/predictions", handler.PredictionHandler{pDao}).Methods("POST")
	r.Handle("/predictions/{id}", handler.PredictionHandler{pDao}).Methods("PUT")
	r.Handle("/predictions/{id}", handler.PredictionHandler{pDao}).Methods("GET")
}

func setupLogging(r http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, r)
}

func init() {
	pDao = dao.PredictionDAO{}
	playerDao = dao.PlayerDAO{}
	tDao = dao.TeamDAO{}
	uDao = dao.UserDAO{}
	bDao = dao.BonusDAO{}
	mDao = dao.MatchesDAO{}
	wsManager = service.NewWebSocketManager()
	wsManager.Start()
}
