package main

import (
	"bufio"
	"os"
	"sync"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	reporter := NewTableReporter()

	for scanner.Scan() {
		s := &Scraper{
			Website:     scanner.Text(),
			MaxDepth:    1,
			Recursively: false,
			PrintLogs:   true,
			Async:       true,
			mutex:       &sync.Mutex{},
			stats:       make(map[string]*PageStats),
		}
		err := s.Scrape()
		if err != nil {
			s.Log("Error while scraping", s.Website)
			continue
		}
		for _, r := range s.Report() {
			reporter.Append(r)
		}
	}

	if scanner.Err() != nil {
		// handle error.
	}
	reporter.Render()
}
