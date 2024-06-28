package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// URLSet represents the <urlset> XML structure
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// URL represents an individual <url> entry
type URL struct {
	Loc string `xml:"loc"`
}

type Sitemap struct {
	filename string
}

// Init sitemap
// Set filename
func (sm *Sitemap) Init(filename string) {
	sm.filename = filename
}

// Add new url in sitemap
func (sm Sitemap) Add(url string) {
	var urlSet URLSet

	// Open the existing file
	file, err := os.OpenFile(sm.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening sitemap.xml: %v\n", err)
		return
	}
	defer file.Close()

	// Read the file content
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading sitemap.xml: %v\n", err)
		return
	}

	// Unmarshal the XML data into URLSet struct
	if len(data) > 0 {
		err = xml.Unmarshal(data, &urlSet)
		if err != nil {
			fmt.Printf("Error unmarshalling XML: %v\n", err)
			return
		}
	} else {
		// If the file is empty, initialize the URLSet with the XML namespace
		urlSet = URLSet{
			XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		}
	}

	// Append the new URL to the URLSet
	urlSet.URLs = append(urlSet.URLs, URL{Loc: url})

	// Marshal the URLSet back to XML
	output, err := xml.MarshalIndent(urlSet, "", "    ")
	if err != nil {
		fmt.Printf("Error marshalling XML: %v\n", err)
		return
	}

	// Truncate the file and write the updated XML data
	file.Truncate(0)
	file.Seek(0, 0)

	_, err = file.Write([]byte(xml.Header))
	if err != nil {
		fmt.Printf("Error writing XML header: %v\n", err)
		return
	}
	_, err = file.Write(output)
	if err != nil {
		fmt.Printf("Error writing XML content: %v\n", err)
		return
	}

	fmt.Printf("%v updated successfully with %v.\n", sm.filename, url)
}
