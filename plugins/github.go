package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/openSUSE/sauria/query"
	"strings"
	"time"
)

type commit struct {
	Date    time.Time
	Message string
	Sha     string
}

func main() {
	var unstable bool
	var generateChangelog bool
	unstable = true

	dom := query.QueryUri("https://github.com/marguerite/util" + "/releases")
	versions := query.QueryDom(dom, "h4[class*=commit-title]")

	if generateChangelog || unstable || len(versions) == 0 {
		dom := query.QueryUri("https://github.com/marguerite/util/commits/master")
		var commits []commit
		timeForm := "Jan 02, 2006"
		dom.Find("div[class=commit-group-title]").Each(func(i int, selection *goquery.Selection) {
			d := strings.Replace(strings.TrimSpace(selection.Text()), "Commits on ", "", -1)
			t, _ := time.Parse(timeForm, d)
			messages := query.QueryDom(dom.Find("ol[class*=commit-group]").Eq(i), "p[class*=commit-title]>a")
			shas := query.QueryDom(dom.Find("ol[class*=commit-group]").Eq(i), "div[class*=commit-links-group]>a")
			for i, j := range messages {
				commits = append(commits, commit{t, j, shas[i]})
			}
		})
		fmt.Println(commits[0])
	}
}
