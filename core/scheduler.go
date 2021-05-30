package core

import (
	"fmt"
	"sync"
	"time"

	config "github.com/karthikic/techblogs/configs"
)

func worker(wg *sync.WaitGroup, name string, source config.Source, quit chan bool) {
	defer wg.Done()

	scrape_interval := source.ScrapeInterval
	if scrape_interval == 0 {
		scrape_interval = int(config.GetScrapeInterval())
	}

	ticker := time.NewTicker(time.Duration(scrape_interval) * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Woke up to scrape:", name, ", Every:", scrape_interval)
		case <-quit:
			fmt.Println("Killing worker for ", name)
			ticker.Stop()
			return
		}
	}

}

func NewScheduler(stop chan bool) {
	wg := &sync.WaitGroup{}

	for name, source := range config.GetSources() {
		wg.Add(1)
		go worker(wg, name, source, stop)
	}

	wg.Wait()

}
