package main

import (
	"bufio"
	"flag"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

// flags
var (
	maxDepthFlag = flag.Int("d", 1, "Maximum depth for crawling")
	verboseFlag  = flag.Bool("v", false, "Write verbose logs")
	formatFlag   = flag.String("f", "table", "Format of the output (table|csv|tsv)")
)

var (
	id int32
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
		return NewTableReporter()
	}
}

func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrencyNumber())
	reporter := NewReporter()
	reports := make(chan *PageStats)

	go func() {
		for r := range reports {
			reporter.Append(r)
			wg.Done()
		}
	}()

	for scanner.Scan() {
		wg.Add(1)
		sem <- struct{}{}
		go func(url string, rc chan *PageStats) {
			defer wg.Done()
			processUrl(url, rc, &wg)
			<-sem
		}(scanner.Text(), reports)
	}
	wg.Wait()
	close(reports)

	reporter.Render()
}

func concurrencyNumber() int {
	count := runtime.NumCPU() / *maxDepthFlag
	if count < 1 {
		count = 1
	}
	return count
}

func processUrl(url string, reports chan *PageStats, wg *sync.WaitGroup) {
	s := &Scraper{
		ID:        atomic.AddInt32(&id, 1),
		Website:   url,
		MaxDepth:  *maxDepthFlag,
		PrintLogs: *verboseFlag,
		Async:     true,
		mutex:     &sync.Mutex{},
		stats:     make(map[string]*PageStats),
	}
	err := s.Scrape()
	if err != nil {
		s.Log("Error while scraping:", err)
		return
	}
	for _, r := range s.Report() {
		wg.Add(1)
		reports <- r
	}
}
