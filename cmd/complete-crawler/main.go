package main

import (
	"log"
	"os"

	"github.com/rfnajid/gooblo-crawler/internal/const/path"
	htmlScraper "github.com/rfnajid/gooblo-crawler/internal/scrape-html"
	sitemapScraper "github.com/rfnajid/gooblo-crawler/internal/scrape-sitemap"
)

func main() {
	// clear output dir & create new
	e := os.RemoveAll(path.OutputDir)
	if e != nil {
		log.Fatal(e)
	}

	if err := os.Mkdir(path.OutputDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	sitemapScraper.BulkScrape()
	htmlScraper.BulkScrape(path.OutputUrl)

}
