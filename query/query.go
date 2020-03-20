package query

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

// QueryUri query an url to get dom tree
func QueryUri(uri string) *goquery.Document {
	dom, err := goquery.NewDocument(uri)
	if err != nil {
		log.Println(err)
	}
	return dom
}

// QueryDom query dom to find matched nodes
func QueryDom(dom interface{}, pattern string) []string {
	var s []string
	d, ok := dom.(interface {
		Find(p string) *goquery.Selection
	})
	if ok {
		d.Find(pattern).Each(func(i int, selection *goquery.Selection) {
			s = append(s, strings.TrimSpace(selection.Text()))
		})
	}
	return s
}
