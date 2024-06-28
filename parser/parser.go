package parser

import (
	"fmt"
	"io"
	"net/url"

	"golang.org/x/net/html"
)

type UrlSetsList []string

// Extract link from body and filter those who are not from same host
func ParseLinkFromBody(body io.ReadCloser, targetUrl string) []string {
	defer body.Close()
	foundLinks := make(UrlSetsList, 0)

	parsedLink, err := url.Parse(targetUrl)
	if err != nil {
		fmt.Println("Error parsing link:", err)
		return foundLinks
	}

	z := html.NewTokenizer(body)
	for {
		next := z.Next()

		if next == html.ErrorToken {
			return foundLinks.unique()
		}

		if next == html.StartTagToken || next == html.SelfClosingTagToken {
			t := z.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						link := a.Val
						foundLinkParsed, err := url.Parse(link)
						if err != nil {
							fmt.Println("Error parsing link:", err)
							return foundLinks.unique()
						}

						if foundLinkParsed.Host == parsedLink.Host {
							foundLinks = append(foundLinks, link)
						}

						break
					}
				}
			}
		}
	}
}

// Return array of urls and remove duplicated links from html body parse
func (elements UrlSetsList) unique() []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}
