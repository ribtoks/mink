package main

import (
	"bufio"
	"flag"
	"os"
	"sync"
)

// flags
var (
	maxDepthFlag = flag.Int("d", 1, "Maximum depth for crawling")
	verboseFlag  = flag.Bool("v", false, "Write verbose logs")
	formatFlag   = flag.String("f", "table", "Format of the output table|csv|tsv")
)

func NewReporter() Reporter {
	switch *formatFlag {
	case "table":
		return NewTableReporter()
	case "csv":
		return NewCSVReporter()
	case "tsv":
		return NewTSVReporter()
	default:
		return nil
	}
}

func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)

	reporter := NewReporter()

	for scanner.Scan() {
		s := &Scraper{
			Website:   scanner.Text(),
			MaxDepth:  *maxDepthFlag,
			PrintLogs: *verboseFlag,
			Async:     true,
			mutex:     &sync.Mutex{},
			stats:     make(map[string]*PageStats),
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
