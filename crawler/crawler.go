package crawler

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/IvanBozhakov/seo-tools/parser"
)

type HTTPClient interface {
	Get(url string) (res *http.Response, err error)
}

type Buffer struct {
	Mutex   sync.Mutex
	Visited map[string]bool
}

type Crawler struct {
	Client HTTPClient
	Buffer Buffer
}

type Page struct {
	URL   string
	Body  io.ReadCloser
	Links []string
}

// Create new Crawler
func New() *Crawler {
	return &Crawler{
		Client: &http.Client{},
		Buffer: Buffer{Visited: make(map[string]bool)},
	}
}

// Look into buffer if page is already visited
func (b *Buffer) IsVisited(url string) bool {
	b.Mutex.Lock()
	_, isVisited := b.Visited[url]
	defer b.Mutex.Unlock()
	return isVisited
}

// Scan page and scedule links for crawling
func (c *Crawler) Scan(url string) Page {
	c.Buffer.Mutex.Lock()
	fmt.Printf("Crawl Page %v\n", url)

	page, err := DoRequest(c, Page{url, nil, nil})

	if err != nil {
		fmt.Println(err)
	}

	page.Links = parser.ParseLinkFromBody(page.Body, page.URL)

	c.Buffer.Visited[url] = true
	defer c.Buffer.Mutex.Unlock()
	return page

}
