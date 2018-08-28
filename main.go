package main

import (
	"os"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/gocolly/colly"
	"log"
	"fmt"
	"strings"
	"time"
	"io/ioutil"
	"strconv"
)

type ScrapedData struct {
	Data [][]string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
		//log.Fatal("$PORT must be set")
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	d := scrapeGitHub()
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", d)
	})
	router.Run(":" + port)
}

func asNumber(input string) float64 {
	var str string = strings.TrimSpace(input)
	if strings.HasSuffix(str, "k") {
		str = strings.Replace(str, "k", "",1)
		val, _ := strconv.ParseFloat(str, 32)
		return val * 1000
	}else if strings.HasSuffix(str, "K") {
		str = strings.Replace(str, "K", "",1)
		val, _ := strconv.ParseFloat(str, 32)
		return val * 1000
	}
	val, _ := strconv.ParseFloat(str, 32)
	return  val
}

func scrapeGitHub() ScrapedData {

	github := "github.com"
	url := "https://github.com/%s"

	codewars := "www.codewars.com"
	codewarsurl := "https://www.codewars.com/users/%s"

	codecademy := "www.codecademy.com"
	codecademyurl := "https://www.codecademy.com/%s"

	b, err := ioutil.ReadFile("data/users.csv") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	users := strings.Split(string(b),"\n")

	var ret = make([][]string, len(users))
	var record []string

	c := colly.NewCollector(
		colly.AllowedDomains(github),
		//colly.CacheDir(""),
	)

	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob:  "github.com/*",
		// Set a delay between requests to these domains
		Delay: 9 * time.Second,
		// Add an additional random delay
		RandomDelay: 9 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("div[class='js-yearly-contributions'] h2[class='f4 text-normal mb-2']", func(e *colly.HTMLElement) {
		pos := strings.Index(e.Text, " contribution")
		record = append(record, e.Text[0:pos])
	})

	c.OnHTML("nav > a[aria-selected='false'] > span", func(e *colly.HTMLElement) {
		number := asNumber(e.Text)
		record = append(record, strconv.FormatFloat(number, 'f', 0, 32))
	})

	co := colly.NewCollector(
		colly.AllowedDomains(codewars),
		//colly.CacheDir(""),
	)

	co.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob:  "codewars.com/*",
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 18 * time.Second,
	})

	co.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	co.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	co.OnXML("//div[@class='stat-box'][ancestor::div[@class='stat-container']/h2/text()='Progress']/div[@class='stat'][b/text()='Rank:']", func(e *colly.XMLElement) {
		record = append(record, strings.TrimPrefix(e.Text,"Rank:"))
	})
	co.OnXML("//div[@class='stat-box'][ancestor::div[@class='stat-container']/h2/text()='Progress']/div[@class='stat'][b/text()='Honor:']", func(e *colly.XMLElement) {
		record = append(record, strings.TrimPrefix(e.Text,"Honor:"))
	})

	co.OnXML("//div[@class='stat-box'][ancestor::div[@class='stat-container']/h2/text()='Progress']/div[@class='stat'][b/text()='Leaderboard Position:']", func(e *colly.XMLElement) {
		record = append(record, strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(e.Text,"Leaderboard Position:")),"#"))
	})

	co.OnXML("//div[@class='stat-box'][ancestor::div[@class='stat-container']/h2/text()='Progress']/div[@class='stat'][b/text()='Honor Percentile:']", func(e *colly.XMLElement) {
		record = append(record, strings.TrimPrefix(e.Text,"Honor Percentile:"))
	})

	co.OnXML("//div[@class='stat-box'][ancestor::div[@class='stat-container']/h2/text()='Progress']/div[@class='stat'][b/text()='Total Completed Kata:']", func(e *colly.XMLElement) {
		record = append(record, strings.TrimPrefix(e.Text,"Total Completed Kata:"))
	})

	c1 := colly.NewCollector(
		colly.AllowedDomains(codecademy),
		//colly.CacheDir(""),
	)

	c1.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob:  "codecademy.com/*",
		// Set a delay between requests to these domains
		Delay: 9 * time.Second,
		// Add an additional random delay
		RandomDelay: 9 * time.Second,
	})

	c1.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c1.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c1.OnXML("//main[starts-with(@class,'profiles')]/article[3]//h3[following-sibling::small/text()='total points']", func(e *colly.XMLElement) {
		record = append(record, e.Text)
	})

	c1.OnXML("//main[starts-with(@class,'profiles')]/article[2]//div/article[2]//h3", func(e *colly.XMLElement) {
		record = append(record, e.Text)
	})

	c1.OnXML("//main[starts-with(@class,'profiles')]/article[3]//h3[following-sibling::small/text()='day streak']", func(e *colly.XMLElement) {
		record = append(record, e.Text)
	})


	for i, user := range users {
		row := strings.Split(user,",")
		record = append(record, row[0])
		c.Visit(fmt.Sprintf(url, row[0]))

		if len(row) == 2 || len(row) == 3 {
			co.Visit(fmt.Sprintf(codewarsurl, row[1]))
			if len(record) == 9 {
				end := record[8]
				record = append(record, "", "")
				record[8] = ""
				record[9] = ""
				record[10] = end
			}
			if len(row) == 3 {
				c1.Visit(fmt.Sprintf(codecademyurl, row[2]))
			} else {
				record = append(record, "", "", "")
			}
		}
		if len(row) == 1 {
			record = append(record, "","","","","","","","")
		}

		ret[i] = record
		record = nil
	}
	d := ScrapedData{Data: ret}
	return d
}
