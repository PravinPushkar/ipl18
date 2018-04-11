package scraper

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.wdf.sap.corp/I334816/ipl18/backend/dao"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

type cacheModel struct {
	match  *models.Match
	pCache map[int]*models.PredictionsModel
}

func (c *cacheModel) String() string {
	return fmt.Sprintln(c.match, c.pCache)
}

var tl = strings.ToLower

type cache map[int]*cacheModel

type Updater struct {
	MDao          dao.MatchesDAO
	PDao          dao.PredictionDAO
	TDao          dao.TeamDAO
	PlayerDao     dao.PlayerDAO
	cache         cache
	teamCache     map[string]int
	teamAbbrCache map[string]string
	playerCache   map[string]int
	once          sync.Once
}

func (u *Updater) Update(scrap map[int]*models.ScraperMatchModel) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	u.buildCaches()
	u.assignPoints(scrap)
}

func (u *Updater) assignPoints(scrap map[int]*models.ScraperMatchModel) {
	for mid, _ := range scrap {
		result := scrap[mid]
		//calculate points for prediction
		cTeam := u.teamCache[tl(result.Winner)]
		cMoM := u.playerCache[tl(result.MoM)]
		if !u.cache[mid].match.Lock {
			log.Println("updating result for match", mid)
			status := "completed"
			if result.Abandoned {
				log.Println("match abandoned points rule changed")
				status = "abandoned"
			}
			log.Println(mid, cTeam, cMoM, status)
			u.MDao.UpdateResultById(mid, cTeam, cMoM, status)
		} else {
			u.cache[mid].match.Lock = true
			continue
		}

		if len(u.cache[mid].pCache) != 0 {
			//whatever is to be written
			for pid, pInfo := range u.cache[mid].pCache {
				//some prediction found
				mType := map[bool]int{true: 0, false: 1}[result.Abandoned]
				points := u.getPoints(cTeam, cMoM, mType, pInfo)

				log.Println("adding prediction", pid, result.Winner, result.MoM, cTeam, cMoM, points, pInfo.INumber)
				if err := u.PDao.WritePredictionResult(pid, cTeam, cMoM, points); err != nil {
					log.Println("could not update", pid, err)
				}
			}
		} else {
			log.Println("no predictions for match id", mid)
		}
	}
}

func (u *Updater) getPoints(cTeam, cMoM, mType int, pInfo *models.PredictionsModel) int {
	coin := (pInfo.CoinUsed != nil) && *pInfo.CoinUsed

	if mType == 0 {
		//abandoned
		if pInfo.TeamVote != 0 {
			if coin {
				return 5
			}
			return 1
		}
		//did not vote
		return 0
	}

	tPoints := 0
	if cTeam == pInfo.TeamVote {
		log.Println("vote correct", pInfo)
		switch mType {
		case 1:
			//league
			tPoints = 2
			if coin {
				tPoints *= 5
			}
		case 2:
			//qualifier
			tPoints = 20
		case 3:
			//final
			tPoints = 30
		}
	}

	mPoints := 0
	if cMoM == pInfo.MoMVote {
		mPoints = 1
	}

	return tPoints + mPoints
}

func (u *Updater) buildCaches() {
	log.Println("building caches")
	u.once.Do(func() {
		u.buildPermCaches()
	})

	u.buildMatchCache()
}

func (u *Updater) buildPermCaches() {
	log.Println("building perm cache")
	u.teamCache = map[string]int{}
	u.playerCache = map[string]int{}

	//buildTeamCache
	if info, err := u.TDao.GetAllTeams(); err != nil {
		log.Println("error building team cache")
		panic(err)
	} else {
		for _, v := range info.Teams {
			u.teamCache[tl(v.TeamName)] = v.TeamId
		}
	}

	//buildPlayerCache
	if info, err := u.PlayerDao.GetAllPlayers(); err != nil {
		log.Println("error building player cache")
		panic(err)
	} else {
		for _, v := range info.Players {
			u.playerCache[tl(v.Name)] = v.PlayerId
		}
	}

	log.Println("teamCache", len(u.teamCache), "playerCache", len(u.playerCache))
}

func (u *Updater) buildMatchCache() {
	log.Println("building match cache")
	u.cache = cache{}
	//get all matches
	if matches, err := u.MDao.GetAllMatches(); err != nil || len(matches.Matches) == 0 {
		log.Println("Updater: error building cache", err)
		panic("error getting matches")
	} else {
		for _, m := range matches.Matches {
			u.cache[m.MatchId] = &cacheModel{match: &m}
		}
	}

	//get all predictions
	predMap := map[int]map[int]*models.PredictionsModel{}

	if preds, err := u.PDao.GetAllPredictions(); err != nil {
		log.Println("Updater: error building cache", err)
		panic("error getting predictions")
	} else {
		for _, v := range preds {
			if predMap[v.MatchId] == nil {
				predMap[v.MatchId] = map[int]*models.PredictionsModel{}
			}
			predMap[v.MatchId][v.PredictionId] = v
		}
	}

	//combine
	for k, _ := range u.cache {
		if _, ok := predMap[k]; ok {
			u.cache[k].pCache = predMap[k]
		}
	}

	log.Println("match cache", len(u.cache))
}
