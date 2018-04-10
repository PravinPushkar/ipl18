package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

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

	regMatchNo     *regexp.Regexp
	regMatchWinner *regexp.Regexp
	matches        map[int]*models.ScraperMatchModel
)

func getDocument(url string) *goquery.Document {
	log.Println("scraper: hitting", url)
	doc, err := goquery.NewDocument(url)
	errors.PanicOnErr(err, "scraper: error creating document for goquery")

	return doc
}

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
					"",
					getMoM(s),
				}
			}
		})
	}

	index := 0
	getDocument(winnerUrl).Find("b").Each(func(i int, s *goquery.Selection) {
		if (i & 1) == 1 {
			index = index%len(matches) + 1
			if res := regMatchWinner.FindStringSubmatch(s.Text()); len(res) != 3 {
				errors.PanicOnErr(errScrapingTeam, "scraper: could not parse team name")
			} else {
				matches[index].Winner = res[1]
			}
		}
	})

	for k, v := range matches {
		log.Println("scraper:", k, v)
	}
}

func getMoM(s *goquery.Selection) string {
	player := s.Find("div .gp__cricket__player-match__player__detail__link").Before("span").Text()
	if player != "" {
		log.Println("scraper: found mom data", player)
		return player
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

func init() {
	regMatchNo, _ = regexp.Compile(`^\d+`)
	regMatchWinner, _ = regexp.Compile(`^([a-zA-Z0-9 ]+)(won by)`)
	matches = map[int]*models.ScraperMatchModel{}
}
