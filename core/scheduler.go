package core

import (
	"fmt"
	"sync"
	"syscall"
	"time"

	config "github.com/karthikic/techblogs/configs"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func saveBlogs(db *gorm.DB, blogs []Blogs) {
	db.CreateInBatches(blogs, 100)
}

func alreadyScrapped(db *gorm.DB, company string) []Blogs {
	var blogs []Blogs
	db.Find(&blogs, "company = ?", company)
	return blogs
}

func worker(wg *sync.WaitGroup, db *gorm.DB, name string, source config.Source, quit chan bool) {
	defer wg.Done()

	scrape_interval := source.ScrapeInterval
	if scrape_interval == 0 {
		scrape_interval = int(config.GetScrapeInterval())
	}

	lock := sync.RWMutex{}
	ticker := time.NewTicker(time.Duration(scrape_interval) * time.Second)

	for {
		select {
		case <-ticker.C:

			lock.Lock()

			fmt.Println("Woke up to scrape:", name, ", Every:", scrape_interval)
			blogs, err := GetLatestBlog(name, source, alreadyScrapped(db, name))
			if err == nil {
				fmt.Println("Blogs for ", name, ":", blogs)
				saveBlogs(db, blogs)
			} else {
				fmt.Println("error fetchting blogs for ", name, ":", err)
			}

			if viper.GetBool("test_mode") {
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}

			lock.Unlock()
		case <-quit:
			fmt.Println("Killing worker for ", name)
			ticker.Stop()
			return
		}
	}

}

func NewScheduler(db *gorm.DB, stop chan bool) {
	wg := &sync.WaitGroup{}

	for name, source := range config.GetSources() {
		wg.Add(1)
		go worker(wg, db, name, source, stop)
	}

	wg.Wait()

}
