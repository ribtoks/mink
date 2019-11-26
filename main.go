package main

import (
	"bufio"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	reporter := NewTableReporter()

	for scanner.Scan() {
		s := &Scraper{
			MaxDepth:    1,
			Recursively: false,
			PrintLogs:   true,
			Async:       true,
			pages:       make(chan *PageResponse),
		}
		err := s.Scrape(scanner.Text())
		if err != nil {
			for _, r := range s.Report() {
				reporter.Append(r)
			}
		}
	}

	if scanner.Err() != nil {
		// handle error.
	}
	reporter.Render()
}
