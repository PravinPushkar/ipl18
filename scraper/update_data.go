package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.wdf.sap.corp/I334816/ipl18/backend/cache"
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

type scraperCache map[int]*cacheModel

type Updater struct {
	MDao      dao.MatchesDAO
	PDao      dao.PredictionDAO
	TDao      dao.TeamDAO
	PlayerDao dao.PlayerDAO
	UDao      dao.UserDAO
	cache     scraperCache
}

func (u *Updater) Update(scrap map[int]*models.ScraperMatchModel) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Updater:panicked", r)
		}
	}()

	u.buildCaches()
	u.assignPoints(scrap)
}

func (u *Updater) assignPoints(scrap map[int]*models.ScraperMatchModel) {
	//over all matches whose result is declared
	for mid, _ := range scrap {
		log.Println("Updater: scrap", mid)
		result := scrap[mid]
		cTeam := cache.TeamNameCache[result.Winner].TeamId
		cMoM := cache.PlayerNameCache[result.MoM].PlayerId

		mType := map[bool]int{true: 0, false: 1}[result.Abandoned]

		log.Printf("Updater: match %d -  %v", mid, u.cache[mid].match)
		//if a match is not locked (points already allocated)
		if u.cache[mid].match.Lock == false {
			log.Println("Updater:updating result for match", mid)
			//for updating table
			status := "completed"
			if result.Abandoned {
				log.Println(mid, "match abandoned, points rule changed")
				status = "abandoned"
			}
			log.Println("Updater:match result for insertion", mid, cTeam, cMoM, status)

			if err := u.MDao.UpdateResultById(mid, cTeam, cMoM, status); err != nil {
				log.Println("Updater:unable to update match result", mid, err)
				continue
			}

			u.cache[mid].match.Lock = true
		} else {
			//match locked check next one
			continue
		}

		//if some predictions are there for match
		if len(u.cache[mid].pCache) != 0 {
			//analyze them and allocate points
			for pid, pInfo := range u.cache[mid].pCache {
				//some prediction found
				//todo:will need to change after new column added to table
				points := u.getPoints(cTeam, cMoM, mType, pInfo)

				log.Println("Updater:adding prediction", pid, result.Winner, result.MoM, cTeam, cMoM, points, pInfo.INumber)
				if err := u.PDao.WritePredictionResult(pid, cTeam, cMoM, points); err != nil {
					log.Println("Updater:could not update", pid, err)
					//ignore error, go on to update other predictions
				}

				//update user table also
				log.Println("Updater:updating points for user", pInfo.INumber, "by", points)
				if err := u.UDao.UpdateUserPointsByINumber(points, pInfo.INumber); err != nil {
					log.Println("Updater:error updating points for ", pInfo.INumber, points)
					//can continue on error
				}
			}
		} else {
			log.Println("Updater:no predictions for match id", mid)
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
		log.Println("Updater:vote correct", pInfo)
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
	log.Println("Updater:building caches")
	u.buildMatchCache()
}

func (u *Updater) buildMatchCache() {
	log.Println("Updater:building match cache")
	u.cache = scraperCache{}
	//get all matches
	if matches, err := u.MDao.GetAllMatches(); err != nil || len(matches.Matches) == 0 {
		log.Println("Updater: error building cache", err)
		panic("error getting matches")
	} else {
		for _, m := range matches.Matches {
			//necessary otherwise overwrites happen
			p := m
			u.cache[m.MatchId] = &cacheModel{match: p}
		}
	}

	//get all predictions
	predMap := map[int]map[int]*models.PredictionsModel{}

	if preds, err := u.PDao.GetAllPredictions(); err != nil {
		log.Println("Updater:Updater: error building cache", err)
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
	for mid, _ := range u.cache {
		if _, ok := predMap[mid]; ok {
			u.cache[mid].pCache = predMap[mid]
		}
	}

	log.Println("Updater:match cache", len(u.cache))
}
