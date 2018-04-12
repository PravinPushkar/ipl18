package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

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

func Start() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("scraper panicked", r)
		}
	}()

	urls := []string{}
	getDocument(apiUrl).Find("div #results a").Each(func(i int, s *goquery.Selection) {
		url := s.AttrOr("href", "/")
		log.Println("scraper: found result url:", url)
		urls = append(urls, url)
	})

	for _, url := range urls {
		getDocument(baseUrl + url).Find("div .gp__cricket__gameHeader").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				no := getMatchMetaData(s)
				matches[no] = &models.ScraperMatchModel{
					no,
					"",
					"",
					"To be decided",
					getMoM(s),
					false,
				}
			}
		})
	}

	index := 0
	failed := false
	getDocument(winnerUrl).Find("b").Each(func(i int, s *goquery.Selection) {
		if (i & 1) == 1 {
			index = index%len(matches) + 1
			winner, isAbandoned := getWinnerInfo(s.Text())
			if isAbandoned {
				log.Println("scraper: match abandoned")
				matches[index].Abandoned = true
			} else if winner != "" {
				matches[index].Winner = strings.Trim(winner, " ")
			} else {
				log.Println("scraper: failed to determine result")
				failed = true
			}
		}
	})
	if failed {
		return
	}

	for k, v := range matches {
		log.Println("scraper:", k, v)
	}

	upd8 := Updater{
		PlayerDao: dao.PlayerDAO{},
		TDao:      dao.TeamDAO{},
		PDao:      dao.PredictionDAO{},
		MDao:      dao.MatchesDAO{},
		UDao:      dao.UserDAO{},
	}
	upd8.Update(matches)
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
