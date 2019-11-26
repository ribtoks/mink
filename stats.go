package main

import (
	"bytes"
	"net/url"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
)

type PageStats struct {
	Url                   string
	Domain                string
	StatusCode            int
	Status                string
	Indexibility          string
	Title                 string
	TitleLength           int
	MetaDescription       string
	MetaDescriptionLength int
	MetaKeywords          string
	MetaKeywordsCount     int
	Size                  int
	WordCount             int
	CrawlDepth            int
	Inlinks               int
	UniqueInlinks         int
	Outlinks              int
	UniqueOutlinks        int
	ResponseTimeMillis    int
}

func stripHtml(data []byte) string {
	stripped := bluemonday.StrictPolicy().SanitizeBytes(data)
	return string(stripped)
}

func countWords(s string) int {
	nWords := 0
	inWord := false
	for _, r := range s {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			inWord = false
		} else if inWord == false {
			inWord = true
			nWords++
		}
	}
	return nWords
}

func extractLinks(document *goquery.Document) map[string]int {
	links := make(map[string]int)
	document.Find("a").Each(func(index int, element *goquery.Selection) {
		// See if the href attribute exists on the element
		href, exists := element.Attr("href")
		if !exists {
			return
		}
		href = strings.ToLower(href)
		_, ok := links[href]
		if ok {
			links[href] += 1
		} else {
			links[href] = 1
		}
	})
	return links
}

func (s *Scraper) processPage(p *PageResponse) {
	s.Log("Processing page", p.Url)
	defer s.waitGroup.Done()
	u, err := url.ParseRequestURI(p.Url)
	if err != nil {
		s.Log("Error parsing request url.", err)
		return
	}
	buf := bytes.NewBuffer(p.Data)
	document, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		s.Log("Error loading HTTP response body.", err)
		return
	}

	ps := &PageStats{}
	ps.WordCount = countWords(stripHtml(p.Data))

	links := extractLinks(document)
	for k, v := range links {
		l, err := url.ParseRequestURI(k)
		if err != nil {
			s.Log("Error parsing on-page url.", err)
			continue
		}
		if l.Hostname() == u.Hostname() {
			ps.Inlinks += v
			ps.UniqueInlinks += 1
		} else {
			ps.Outlinks += v
			ps.UniqueOutlinks += 1
		}
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.stats[p.Url] = ps
}
