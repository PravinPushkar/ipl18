package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.wdf.sap.corp/I334816/ipl18/backend/dao"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseUrl   = "http://www.espncricinfo.com"
	apiUrl    = baseUrl + "/series/_/id/8048/season/2018/indian-premier-league/"
	winnerUrl = baseUrl + "/ci/engine/series/1131611.html"
)

var (
	errScrapingPlayer = fmt.Errorf("error getting player info")
	errScrapingTeam   = fmt.Errorf("error getting team info")
	errScrapingMeta   = fmt.Errorf("error getting match metadata info")

	regMatchNo *regexp.Regexp
	matches    map[int]*models.ScraperMatchModel
	teamCache  map[string]int
)

func Start(kill chan bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("scraper panicked", r)
		}
	}()

	timer := time.NewTicker(time.Hour)
	upd8 := Updater{
		PlayerDao: dao.PlayerDAO{},
		TDao:      dao.TeamDAO{},
		PDao:      dao.PredictionDAO{},
		MDao:      dao.MatchesDAO{},
		UDao:      dao.UserDAO{},
	}
	loc, _ := time.LoadLocation("Asia/Kolkata")

	for {
		select {
		case <-timer.C:
			timeNow := time.Now().In(loc)
			log.Println("scraper: time", timeNow)
			if timeNow.Hour() == 3 {
				log.Println("scraper: running")
				getDocument(winnerUrl).Find("div .large-5 b").Each(func(i int, s *goquery.Selection) {
					winner, isAbandoned := getWinnerInfo(s.Text())
					if isAbandoned {
						matches[i+1] = &models.ScraperMatchModel{
							MatchNo:   i + 1,
							Abandoned: true,
						}
					} else {
						if winner != "" {
							matches[i+1] = &models.ScraperMatchModel{
								MatchNo:   i + 1,
								Winner:    winner,
								Abandoned: false,
							}
						}
					}
				})

				matchCount := len(matches)
				getDocument(winnerUrl).Find("span .potMatchLink").Each(func(i int, s *goquery.Selection) {
					if i+1 <= matchCount && !matches[i+1].Abandoned {
						for _, el := range s.Get(0).Attr {
							if el.Key == "href" {
								getDocument(el.Val).Find("div .gp__cricket__gameHeader").Each(func(im int, s *goquery.Selection) {
									if im == 0 {
										no := getMatchMetaData(s)
										if _, ok := matches[no]; ok {
											matches[no].MoM = getMoM(s)
										}
									}
								})
							}
						}
					}
				})

				for k, v := range matches {
					log.Println("scraper:", k, v)
				}

				upd8.Update(matches)
			}

		case <-kill:
			return
		}
	}
}

func getDocument(url string) *goquery.Document {
	log.Println("scraper: hitting", url)
	doc, err := goquery.NewDocument(url)
	errors.PanicOnErr(err, "scraper: error creating document for goquery")

	return doc
}

func getMoM(s *goquery.Selection) string {
	player := s.Find("div .gp__cricket__player-match__player__detail__link").Before("span").Text()
	if player != "" {
		log.Println("scraper: found mom data", player)
		return strings.Trim(player, " ")
	} else {
		errors.PanicOnErr(errScrapingPlayer, "scraper:")
	}

	return ""
}

func getMatchMetaData(s *goquery.Selection) int {
	if meta := s.Find("div .cscore_info-overview").Text(); meta != "" {
		if tokens := strings.Split(meta, ", "); len(tokens) != 3 {
			errors.PanicOnErr(errScrapingMeta, "scraper: meta data info invalid")
		} else {
			if matchNo := regMatchNo.FindString(tokens[0]); matchNo == "" {
				errors.PanicOnErr(errScrapingMeta, "scraper: meta data info invalid (match number)")
			} else {
				log.Println("scraper: found match number", matchNo)
				matchNoNum, _ := strconv.Atoi(matchNo)
				return matchNoNum
			}
		}
	} else {
		errors.PanicOnErr(errScrapingMeta, "scraper: could not find div")
	}
	return -1
}

func getWinnerInfo(data string) (string, bool) {
	log.Println("scraper: getting winner info")

	if strings.Contains(data, " abandoned ") {
		log.Println("scraper: match abandoned")
		return "", true
	} else {
		//search for teams
		for k, _ := range teamCache {
			if strings.Contains(data, k) {
				log.Println("scraper: found team", k)
				return k, false
			}
		}
		log.Println("scraper: winner team not found in tied match")
		return "", false
	}
	return "", false
}

func init() {
	regMatchNo, _ = regexp.Compile(`^\d+`)
	matches = map[int]*models.ScraperMatchModel{}
	teamCache = map[string]int{}
	tdao := dao.TeamDAO{}
	if info, err := tdao.GetAllTeams(); err != nil {
		log.Println("error building team cache")
	} else {
		for _, v := range info.Teams {
			teamCache[v.TeamName] = v.TeamId
		}
	}
	log.Println("teamCache", teamCache)
}
