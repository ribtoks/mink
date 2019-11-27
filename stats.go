package main

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"time"
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

func isRedirect(document *goquery.Document) bool {
	hasCanonicalLink := false
	document.Find("link").Each(func(i int, s *goquery.Selection) {
		if rel, _ := s.Attr("rel"); rel == "canonical" {
			hasCanonicalLink = true
		}
	})
	return hasCanonicalLink
}

func Indexibility(statusCode int, document *goquery.Document) string {
	indexable := statusCode/100 == 2
	if indexable {
		indexable = !isRedirect(document)
	}

	if indexable {
		return "Indexable"
	}
	return "Non-indexable"
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

func (ps *PageStats) countLinks(hostname string, links map[string]int) {
	for k, v := range links {
		l, err := url.ParseRequestURI(k)
		if err != nil {
			continue
		}
		if l.IsAbs() && l.Hostname() != hostname {
			ps.Outlinks += v
			ps.UniqueOutlinks += 1
		} else {
			ps.Inlinks += v
			ps.UniqueInlinks += 1
		}
	}
}

func extractMetaDescription(document *goquery.Document) string {
	description := ""
	document.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); name == "description" {
			description, _ = s.Attr("content")
		}
	})
	return description
}

func extractMetaKeywords(document *goquery.Document) (string, int) {
	keywords := ""
	document.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); name == "keywords" {
			keywords, _ = s.Attr("content")
		}
	})
	if len(keywords) == 0 {
		return "", 0
	}
	arr := strings.Split(keywords, ",")
	return keywords, len(arr)
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
	ps.Url = p.Url
	ps.Domain = u.Hostname()
	ps.StatusCode = p.StatusCode
	ps.Status = http.StatusText(p.StatusCode)
	ps.Indexibility = Indexibility(p.StatusCode, document)
	ps.Title = document.Find("title").Text()
	ps.TitleLength = len(ps.Title)
	ps.MetaDescription = extractMetaDescription(document)
	ps.MetaDescriptionLength = len(ps.MetaDescription)
	keywords, keywordsCount := extractMetaKeywords(document)
	ps.MetaKeywords = keywords
	ps.MetaKeywordsCount = keywordsCount
	ps.Size = len(p.Data)
	ps.WordCount = countWords(stripHtml(p.Data))
	ps.CrawlDepth = p.Depth
	ps.ResponseTimeMillis = int(p.Duration / time.Millisecond)

	links := extractLinks(document)
	s.Log("Page links:", links)
	ps.countLinks(strings.ToLower(u.Hostname()), links)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.stats[p.Url] = ps
}
