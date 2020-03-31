package query

import (
	"log"

	"github.com/PuerkitoBio/goquery"
)

// URI query an url to get dom tree
func URI(uri string) *goquery.Document {
	dom, err := goquery.NewDocument(uri)
	if err != nil {
		log.Println(err)
	}
	return dom
}

// Dom query dom to find matched nodes
func Dom(dom interface{}, pattern string) []*goquery.Selection {
	var s []*goquery.Selection
	d, ok := dom.(interface {
		Find(p string) *goquery.Selection
	})
	if ok {
		d.Find(pattern).Each(func(i int, selection *goquery.Selection) {
			s = append(s, selection)
		})
	}
	return s
}
