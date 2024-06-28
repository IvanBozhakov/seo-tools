package main

import (
	"sync"

	"github.com/IvanBozhakov/seo-tools/crawler"
	"github.com/IvanBozhakov/seo-tools/parser"
)

func main() {
	cr := crawler.New()

	linksChan := make(chan []string)

	var wg sync.WaitGroup

	page := cr.Scan("https://example.bg")

	sitemap := (parser.Sitemap{})
	sitemap.Init("sitemap.xml")

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for linkList := range linksChan {
				for i := 0; i < len(linkList); i++ {
					var link string = linkList[i]
					if cr.Buffer.IsVisited(link) {
						continue
					}

					page := cr.Scan(link)

					sitemap.Add(page.URL)

					go func(p crawler.Page) {
						linksChan <- p.Links
					}(page)
				}
			}
		}()
	}

	linksChan <- page.Links

	wg.Wait()

}
