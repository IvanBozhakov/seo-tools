package crawler

// Create GET request and fetch push body into chanel
func DoRequest(crawler *Crawler, p Page) (Page, error) {
	result, err := crawler.Client.Get(p.URL)
	if err != nil {
		return p, err
	}
	p.Body = result.Body

	return p, nil
}
