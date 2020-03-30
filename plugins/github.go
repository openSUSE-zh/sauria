package main

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/openSUSE/sauria/common"
	"github.com/openSUSE/sauria/query"
)

// FetchNewVersion fetch a new version from github
func FetchNewVersion(uri string, unstable, genChange bool) (common.Commit, error) {
	dom := query.QueryUri(uri + "/releases")
	versions := query.QueryDom(dom, "h4[class*=commit-title]")

	if genChange || unstable || len(versions) == 0 {
		dom := query.QueryUri(uri + "/commits/master")
		var commits []common.Commit
		timeForm := "Jan 02, 2006"
		dom.Find("div[class=commit-group-title]").Each(func(i int, selection *goquery.Selection) {
			d := strings.Replace(strings.TrimSpace(selection.Text()), "Commits on ", "", -1)
			t, _ := time.Parse(timeForm, d)
			messages := query.QueryDom(dom.Find("ol[class*=commit-group]").Eq(i), "p[class*=commit-title]>a")
			shas := query.QueryDom(dom.Find("ol[class*=commit-group]").Eq(i), "div[class*=commit-links-group]>a")
			for i, j := range messages {
				commits = append(commits, common.Commit{t, j, shas[i]})
			}
		})
		return commits[0], nil
	}
	return common.Commit{}, nil
}
