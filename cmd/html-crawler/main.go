package main

import (
	"fmt"

	path "github.com/rfnajid/gooblo-crawler/internal/const/path"
	scraper "github.com/rfnajid/gooblo-crawler/internal/scrape-html"
)

func main() {
	fmt.Println("HTML Crawler....")

	scraper.BulkScrape(path.InputUrl)
}
