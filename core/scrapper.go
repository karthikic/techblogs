package core

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	config "github.com/karthikic/techblogs/configs"
	"mvdan.cc/xurls/v2"
)

func GetLatestBlog(company string, source config.Source, alreadyScrapped []Blogs) ([]Blogs, error) {
	titles := []Blogs{}

	page := 1
	for true {

		new_titles, next := getBlogsPerPage(company, page, source, alreadyScrapped)
		if new_titles != nil {

			titles = append(titles, new_titles...)
		}

		if !next {
			break
		}

		page = page + 1
	}
	return titles, nil
}

func getBlogsPerPage(company string, page int, source config.Source, alreadyScrapped []Blogs) ([]Blogs, bool) {
	// Get the HTML
	xurlsStrict := xurls.Strict()
	url := fmt.Sprintf("%s%s%s", source.Url, source.PagePath, fmt.Sprint(page))

	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, false
	}

	// Convert HTML into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, false
	}

	// Save each .post-title as a list
	var blogs []Blogs
	doc.Find(source.TitleKey).Each(func(i int, s *goquery.Selection) {

		found := false
		title := s.Text()
		raw_link, _ := s.Children().Html()
		links := xurlsStrict.FindAllString(raw_link, -1)
		link := ""
		if len(links) > 0 {
			link = links[len(links)-1]
		}
		for _, v := range alreadyScrapped {
			if v.Title == title && v.Link == link {
				found = true
				fmt.Print("Already Saved")
				break
			}
		}

		if !found {

			blog := Blogs{
				Title:   title,
				Company: company,
				Link:    link,
			}
			blogs = append(blogs, blog)
		}

		found = false
	})

	return blogs, true
}
