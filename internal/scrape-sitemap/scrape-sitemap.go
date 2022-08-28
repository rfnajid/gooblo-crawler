package scrapesitemap

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"

	colors "github.com/rfnajid/gooblo-crawler/internal/const/colors"
	path "github.com/rfnajid/gooblo-crawler/internal/const/path"
)

var sitemaps = []string{}

func BulkScrape() {
	fmt.Println(colors.Green, "Starting bulk scraping sitemap...", colors.Reset)

	argsLength := len(os.Args[1:])

	if argsLength > 0 {
		sitemaps = os.Args[1:]
	} else {
		sitemaps = inputFromFile()
	}

	// clear output dir & create new
	e := os.RemoveAll(path.OutputDir)
	if e != nil {
		log.Fatal(e)
	}

	if err := os.Mkdir(path.OutputDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	for _, item := range sitemaps {
		Scrape(item)
	}

	fmt.Println(colors.Green, "Bulk Scraping sitemap done", colors.Reset)
}

func Scrape(url string) {
	fmt.Println(colors.Yellow, "Scraping XML : ", colors.Reset, url)

	// Array containing all the known URLs in a sitemap
	knownUrls := []string{}
	knownSitemaps := []string{}

	c := colly.NewCollector()

	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		knownUrls = append(knownUrls, e.Text)
	})

	c.OnXML("//sitemapindex/sitemap/loc", func(e *colly.XMLElement) {
		knownSitemaps = append(knownSitemaps, e.Text)
	})

	// Start the collector
	c.Visit(url)

	if len(knownSitemaps) > 0 {
		for _, item := range knownSitemaps {
			Scrape(item)
		}
	}

	fmt.Println("Saving all Urls...")
	urlFile, err := os.OpenFile(path.OutputUrl,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	for _, knownUrl := range knownUrls {
		knownUrl = strings.TrimSpace(knownUrl)
		if _, err := urlFile.WriteString(knownUrl + "\n"); err != nil {
			log.Println(err)
		}
	}
	fmt.Println("Collected", len(knownUrls), "URLs")

	fmt.Println("Saving all sitemaps...")
	sitemapFile, err := os.OpenFile(path.OutputSitemap,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	for _, knownSitemap := range knownSitemaps {
		knownSitemap = strings.TrimSpace(knownSitemap)
		if _, err := sitemapFile.WriteString(knownSitemap + "\n"); err != nil {
			log.Println(err)
		}
	}
	fmt.Println("Collected", len(knownSitemaps), "Sitemaps")

	urlFile.Close()
	sitemapFile.Close()
}

func inputFromFile() (res []string) {

	inputSrc := path.InputSitemap

	f, err := os.Open(inputSrc)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if len(url) > 0 {
			res = append(res, url)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}
