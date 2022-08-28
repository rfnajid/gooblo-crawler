package scrapehtml

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"

	colors "github.com/rfnajid/gooblo-crawler/internal/const/colors"
	path "github.com/rfnajid/gooblo-crawler/internal/const/path"
)

type ScrapeResult struct {
	url, title, keywords, description string
}

var outputResult = path.OutputScrape

func Init() {
	file, err := os.Create(outputResult)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", outputResult, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"url", "title", "keywords", "description"})
}

func appendOutput(sr ScrapeResult) {

	f, err := os.OpenFile(outputResult, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	w := csv.NewWriter(f)
	w.Write([]string{sr.url, sr.title, sr.keywords, sr.description})
	w.Flush()
}

func Scrape(url string) {

	// Instantiate default collector
	c := colly.NewCollector()
	res := ScrapeResult{url: url}

	c.OnHTML("html", func(e *colly.HTMLElement) {

		res.title = e.DOM.Find("title").Text()

		// Extract meta tags from the document
		metaTags := e.DOM.Find("meta")
		metaTags.Each(func(_ int, s *goquery.Selection) {

			if name, _ := s.Attr("name"); strings.EqualFold(name, "description") {
				res.description, _ = s.Attr("content")
				fmt.Println(colors.Purple, "desc : \n", colors.Reset, res.description)
			}

			if name, _ := s.Attr("name"); strings.EqualFold(name, "keywords") {
				res.keywords, _ = s.Attr("content")
				fmt.Println(colors.Purple, "keywords: \n", colors.Reset, res.keywords)
			}
		})

		appendOutput(res)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(colors.Yellow, "Scraping : ", colors.Reset, r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(colors.Green, "Finished : ", colors.Reset, r.Request.URL)
	})

	c.Visit(url)
}

func BulkScrape(inputSrc string) {
	fmt.Println(colors.Green, "Bulk Scraping Url...", colors.Reset)

	f, err := os.Open(inputSrc)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	Init()

	counter := 0
	for scanner.Scan() {
		// do something with a line
		url := strings.TrimSpace(scanner.Text())
		if len(url) > 0 {
			Scrape(url)
			counter++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(colors.Green, "Bulk Scraping URls completed!!! Total : ", counter, colors.Reset)
}
