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
	"regexp"
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
	re := regexp.MustCompile("st|nd|rd|th")

	codecademy := "www.codecademy.com"
	codecademyurl := "https://www.codecademy.com/%s"

	//datacamp := "www.datacamp.com"
	//datacampurl := "https://www.datacamp.com/profile/%s"

	//khan := "www.khanacademy.org"
	//khanurl := "https://www.khanacademy.org/profile/%s"
	//re2 := regexp.MustCompile("(.*\"points\": *)(\\d+)(,.*)")

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
		Delay: 10 * time.Second,
		// Add an additional random delay
		RandomDelay: 10 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("div[class='js-yearly-contributions'] h2[class='f4 text-normal mb-2']", func(e *colly.HTMLElement) {
		pos := strings.Index(e.Text, " contribution")
		record = append(record, strings.TrimSpace(strings.Replace(e.Text[0:pos],",", "", -1)))
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
		Delay: 25 * time.Second,
		// Add an additional random delay
		RandomDelay: 35 * time.Second,
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
		record = append(record, strings.Replace(strings.TrimPrefix(e.Text,"Honor:"),",","",-1))
	})

	co.OnXML("//div[@class='stat-box'][ancestor::div[@class='stat-container']/h2/text()='Progress']/div[@class='stat'][b/text()='Leaderboard Position:']", func(e *colly.XMLElement) {
		record = append(record, strings.Replace(strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(e.Text,"Leaderboard Position:")),"#"),",","",-1))
	})

	co.OnXML("//div[@class='stat-box'][ancestor::div[@class='stat-container']/h2/text()='Progress']/div[@class='stat'][b/text()='Honor Percentile:']", func(e *colly.XMLElement) {
		record = append(record, re.ReplaceAllString(strings.TrimPrefix(e.Text,"Honor Percentile:"), ""))
	})

	co.OnXML("//div[@class='stat-box'][ancestor::div[@class='stat-container']/h2/text()='Progress']/div[@class='stat'][b/text()='Total Completed Kata:']", func(e *colly.XMLElement) {
		record = append(record, strings.Replace(strings.TrimPrefix(e.Text,"Total Completed Kata:"),",","",-1))
	})

	c1 := colly.NewCollector(
		colly.AllowedDomains(codecademy),
		//colly.CacheDir(""),
	)

	c1.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob:  "codecademy.com/*",
		// Set a delay between requests to these domains
		Delay: 10 * time.Second,
		// Add an additional random delay
		RandomDelay: 10 * time.Second,
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

	/*
	c2 := colly.NewCollector(
		colly.AllowedDomains(datacamp),
		//colly.CacheDir(""),
	)

	c2.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob:  "datacamp.com/*",
		// Set a delay between requests to these domains
		Delay: 9 * time.Second,
		// Add an additional random delay
		RandomDelay: 9 * time.Second,
	})

	c2.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c2.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c2.OnXML("//div[contains(@class,'profile-header__stats')]/div[1]//strong", func(e *colly.XMLElement) {
		record = append(record, e.Text)
	})

	c2.OnXML("//div[contains(@class,'profile-header__stats')]/div[2]//strong", func(e *colly.XMLElement) {
		record = append(record, e.Text)
	})

	c2.OnXML("//div[contains(@class,'profile-header__stats')]/div[3]//strong", func(e *colly.XMLElement) {
		record = append(record, e.Text)
	})
*/
/*
	c3 := colly.NewCollector(
		colly.AllowedDomains(khan),
		//colly.CacheDir(""),
	)

	c3.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob:  "khanacademy.org/*",
		// Set a delay between requests to these domains
		Delay: 9 * time.Second,
		// Add an additional random delay
		RandomDelay: 9 * time.Second,
	})

	c3.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c3.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c3.OnXML("//script[contains(text(),'Profile.init')]", func(e *colly.XMLElement) {
		record = append(record, re2.ReplaceAllString(e.Text,"${2}"))
	})
*/

	for i, user := range users {
		po := float64(0)
		row := strings.Split(user,",")
		record = append(record, row[0])
		c.Visit(fmt.Sprintf(url, row[0]))

		if len(row) == 2 || len(row) == 3 || len(row) == 4 {
			co.Visit(fmt.Sprintf(codewarsurl, row[1]))
			if len(record) == 9 {
				end := record[8]
				record = append(record, "", "")
				record[8] = ""
				record[9] = ""
				record[10] = end
			}
			if len(row) == 2 {
				record = append(record, "","","")
			}
			if len(row) == 3 || len(row) == 4 {
				c1.Visit(fmt.Sprintf(codecademyurl, row[2]))
				/*
				if len(row) != 4 {
					record = append(record, "")
				}
				if len(row) == 4 {
					c3.Visit(fmt.Sprintf(khanurl, row[3]))
				}*/
			}
		}
		if len(row) == 1 {
			record = append(record, "","","","","","","","")
		}

		for j:=1; j<14; j++ {
			if j != 5 && j != 6 && j != 8 && record[j] != "" {
				points, _ := strconv.ParseFloat(record[j], 32)
				po += points
			}
		}
		record = append(record, "", "")
		copy(record[2:], record[0:])
		record[0] = strconv.FormatFloat(po, 'f', 0, 32)
		badge := "1"
		if po > 100000 {
			badge = "14"
		} else if po > 80000 {
			badge = "13"
		} else if po > 60000 {
			badge = "12"
		} else if po > 50000 {
			badge = "11"
		} else if po > 40000 {
			badge = "10"
		} else if po > 30000 {
			badge = "9"
		} else if po > 20000 {
			badge = "8"
		} else if po > 15000 {
			badge = "7"
		} else if po > 12000 {
			badge = "6"
		} else if po > 10000 {
			badge = "5"
		} else if po > 7500 {
			badge = "4"
		} else if po > 3500 {
			badge = "3"
		} else if po > 1000 {
			badge = "2"
		} else {
			badge = "1"
		}
		record[1] = badge

		ret[i] = record
		record = nil
	}
	d := ScrapedData{Data: ret}
	return d
}
