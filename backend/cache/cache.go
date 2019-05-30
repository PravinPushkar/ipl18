package cache

import (
	"log"
	"sync"

	"github.com/PravinPushkar/ipl18/backend/models"
)

var (
	TeamNameCache  map[string]*models.Team
	TeamSNameCache map[string]*models.Team
	TeamIdCache    map[int]*models.Team
	Teams          *models.Teams

	UserINumberCache map[string]*models.User

	PlayerIdCache   map[int]*models.Player
	PlayerNameCache map[string]*models.Player
	Players         *models.PlayersModel

	TeamPlayerCache map[int]map[int]*models.Player

	MatchInfoCache map[int]*models.Match

	Lock *sync.RWMutex
)

func init() {
	Lock = &sync.RWMutex{}
	buildTeamCaches()
	buildPlayerCaches()
	buildUserCaches()
	buildMatchCaches()
	log.Println("Done")
}

func buildTeamCaches() {
	log.Println("Building team caches")
	TeamIdCache = make(map[int]*models.Team)
	TeamSNameCache = make(map[string]*models.Team)
	TeamNameCache = make(map[string]*models.Team)
	TeamPlayerCache = make(map[int]map[int]*models.Player)
}

func buildPlayerCaches() {
	log.Println("Building player caches")
	PlayerIdCache = make(map[int]*models.Player)
	PlayerNameCache = make(map[string]*models.Player)
}

func buildUserCaches() {
	log.Println("Building user caches")
	UserINumberCache = make(map[string]*models.User)
}

func buildMatchCaches() {
	log.Println("Building match caches")
	MatchInfoCache = make(map[int]*models.Match)
}
