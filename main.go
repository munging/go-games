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

func scrapeGitHub() ScrapedData {

	github := "github.com"
	url := "https://github.com/%s"
	users := []string{"jbampton", "ugifractal", "giacomosorbi", "tsara27", "scottyrs", "udha", "prestonhunter",
		"petraruttiger", "grfxwzdesigner", "summerhill5"}
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
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("div[class='js-yearly-contributions'] h2[class='f4 text-normal mb-2']", func(e *colly.HTMLElement) {
		pos := strings.Index(e.Text, " contributions")
		record = append(record, e.Text[0:pos])
	})

	c.OnHTML("nav > a[aria-selected='false'] > span", func(e *colly.HTMLElement) {
		record = append(record, e.Text)
	})

	for i, user := range users {
		record = append(record, user)
		c.Visit(fmt.Sprintf(url, user))
		ret[i] = record
		record = nil

	}
	d := ScrapedData{Data: ret}
	return d
}
