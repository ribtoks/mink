package main

import (
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type PageResponse struct {
	Url        string
	StatusCode int
	Data       []byte
}

type Scraper struct {
	MaxDepth    int
	Website     string
	Recursively bool
	PrintLogs   bool
	Async       bool
	pages       chan *PageResponse
	waitGroup   sync.WaitGroup
	stats       map[string]*PageStats
	mutex       *sync.Mutex
}

func prepareAllowedDomain(requestURL string) ([]string, error) {
	requestURL = "https://" + trimProtocol(requestURL)
	u, err := url.ParseRequestURI(requestURL)
	if err != nil {
		return nil, err
	}
	hostname := u.Hostname()
	domain := strings.TrimLeft(hostname, "wwww.")
	return []string{
		domain,
		"www." + domain,
		"http://" + domain,
		"https://" + domain,
		"http://www." + domain,
		"https://www." + domain,
	}, nil
}

func trimProtocol(requestURL string) string {
	return strings.Trim(strings.Trim(requestURL, "http://"), "https://")
}

func (s *Scraper) Log(v ...interface{}) {
	if s.PrintLogs {
		log.Println(v)
	}
}

func (s *Scraper) GetWebsite(secure bool) string {
	if secure {
		return "https://" + s.Website
	}
	return "http://" + s.Website
}

func (s *Scraper) Scrape(url string) error {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"

	c.MaxDepth = s.MaxDepth
	c.Async = s.Async
	allowedDomains, err := prepareAllowedDomain(s.Website)
	if err != nil {
		return err
	}
	c.AllowedDomains = allowedDomains
	s.Website = trimProtocol(s.Website)

	c.OnResponse(func(r *colly.Response) {
		p := &PageResponse{
			Url:        r.Request.URL.String(),
			Data:       r.Body,
			StatusCode: r.StatusCode,
		}
		go func(pp *PageResponse) {
			s.waitGroup.Add(1)
			s.pages <- pp
		}(p)
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if s.Recursively {
			link := e.Attr("href")
			s.Log("visiting: ", link)
			if err := e.Request.Visit(link); err != nil {
				if err != colly.ErrAlreadyVisited {
					s.Log("error while linking: ", err.Error())
				}
			}
		}
	})

	// Start the scrape
	if err := c.Visit(s.GetWebsite(true)); err != nil {
		s.Log("error while visiting: ", err.Error())
	}

	// Wait for concurrent scrapes to finish
	c.Wait()
	// wait for page processing to finish
	s.waitGroup.Wait()

	return nil
}

func (s *Scraper) Report() []*PageStats {
	result := make([]*PageStats, len(s.stats))
	for _, v := range s.stats {
		result = append(result, v)
	}
	return result
}
