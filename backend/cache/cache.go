package cache

import (
	"log"

	"github.wdf.sap.corp/I334816/ipl18/backend/dao"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

var (
	TeamNameCache  map[string]*models.Team
	TeamSNameCache map[string]*models.Team
	TeamIdCache    map[int]*models.Team

	UserINumberCache map[string]*models.UserBasic

	PlayerIdCache   map[int]*models.Player
	PlayerNameCache map[string]*models.Player

	TeamPlayerCache map[int]map[int]*models.Player

	MatchInfoCache map[int]*models.Match

	tDao dao.TeamDAO
	pDao dao.PlayerDAO
	mDao dao.MatchesDAO
	uDao dao.UserDAO
)

func init() {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Println("panicked", r)
	// 		os.Exit(1)
	// 	}
	// }()

	tDao = dao.TeamDAO{}
	pDao = dao.PlayerDAO{}
	uDao = dao.UserDAO{}
	mDao = dao.MatchesDAO{}

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

	teams, err := tDao.GetAllTeams()
	if err != nil {
		panic("could not get teams")
	}

	for _, team := range teams.Teams {
		TeamIdCache[team.TeamId] = team
		TeamSNameCache[team.ShortName] = team
		TeamNameCache[team.TeamName] = team
		TeamPlayerCache[team.TeamId] = make(map[int]*models.Player)
	}
}

func buildPlayerCaches() {
	log.Println("Building player caches")
	PlayerIdCache = make(map[int]*models.Player)
	PlayerNameCache = make(map[string]*models.Player)

	players, err := pDao.GetAllPlayers()
	if err != nil {
		panic("could not get players")
	}

	for _, player := range players.Players {
		PlayerIdCache[player.PlayerId] = player
		PlayerNameCache[player.Name] = player
		TeamPlayerCache[player.TeamId][player.PlayerId] = player
	}
}

func buildUserCaches() {
	log.Println("Building user caches")
	UserINumberCache = make(map[string]*models.UserBasic)
	users, err := uDao.GetAllUsersBasicInfo()
	if err != nil {
		panic("could not get players")
	}

	for _, user := range users {
		UserINumberCache[user.INumber] = user
	}
}

func buildMatchCaches() {
	log.Println("Building match caches")
	MatchInfoCache = make(map[int]*models.Match)
	matches, err := mDao.GetAllMatches()
	if err != nil {
		panic("could not get players")
	}

	for _, match := range matches.Matches {
		MatchInfoCache[match.MatchId] = match
	}
}
