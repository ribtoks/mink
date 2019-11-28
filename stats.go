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

const (
	INDEXABLE     = "Indexable"
	NON_INDEXABLE = "Non-Indexable"
	NOINDEX       = "noindex"
)

type PageStats struct {
	Url                   string
	Domain                string
	StatusCode            int
	Status                string
	Indexibility          string
	ContentType           string
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
	Emails                string
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

func isRedirect(hostname string, document *goquery.Document) bool {
	hasCanonicalLink := false
	document.Find("link").Each(func(i int, s *goquery.Selection) {
		if rel, _ := s.Attr("rel"); rel == "canonical" {
			if href, exists := s.Attr("href"); exists {
				if u, err := url.ParseRequestURI(href); err == nil {
					hasCanonicalLink = hostname != u.Hostname()
				}
			}
		}
	})
	return hasCanonicalLink
}

func isNoIndex(document *goquery.Document) bool {
	hasNoIndex := false
	document.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); name == "robots" {
			if content, exists := s.Attr("content"); exists {
				hasNoIndex = strings.Contains(content, NOINDEX)
			}
		}
	})
	return hasNoIndex
}

func Indexibility(p *PageResponse, hostname string, document *goquery.Document) string {
	if p.StatusCode/100 != 2 {
		return NON_INDEXABLE
	}

	robots := p.Headers.Get("X-Robots-Tag")
	if strings.Contains(robots, NOINDEX) {
		return NON_INDEXABLE
	}

	if isNoIndex(document) {
		return NON_INDEXABLE
	}

	if isRedirect(hostname, document) {
		return NON_INDEXABLE
	}

	return INDEXABLE
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
		if name, _ := s.Attr("name"); strings.ToLower(name) == "description" {
			description, _ = s.Attr("content")
		}
	})
	return description
}

func extractMetaKeywords(document *goquery.Document) (string, int) {
	keywords := ""
	document.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); strings.ToLower(name) == "keywords" {
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
	s.Log("Processing page", p.Url, len(p.Data), "bytes")
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
	ps.Indexibility = Indexibility(p, u.Hostname(), document)
	ps.ContentType = p.Headers.Get("Content-Type")
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
	ps.Emails = strings.Join(parseEmails(p.Data), ";")

	links := extractLinks(document)
	s.Log("Page links:", links)
	ps.countLinks(strings.ToLower(u.Hostname()), links)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.stats[p.Url] = ps
}
