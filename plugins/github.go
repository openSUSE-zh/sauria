package main

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/openSUSE/sauria/commit"
	"github.com/openSUSE/sauria/configuration"
	"github.com/openSUSE/sauria/query"
)

// FetchNewVersion fetch a new version from github
func FetchNewVersion(c *configuration.Configuration) (commit.Commit, error) {
	var revision commit.Commit
	timeForm := "Jan 02, 2006"
	uri := c.StringAttr("url")
	nondigit := c.BoolAttr("nondigit")
	unstable := c.BoolAttr("unstable")
	genChange := c.BoolAttr("genchange")

	dom := query.Dom(query.URI(uri+"/releases"), "div[class=release-entry]")
	// len(dom) == 0 no release
	if len(dom) > 0 {
		dom1 := query.Dom(dom[0], "div[class*=label-latest]")
		next := false
		// normal release
		if len(dom1) > 0 {
			release := strings.TrimSpace(query.Dom(dom1[0], "div ul li a span")[0].Text())
			dom2 := query.Dom(dom1[0], "div[class*=release-main-section]")
			date := strings.TrimSpace(query.Dom(dom2[0], "div[class=release-header] p relative-time")[0].Text())
			messages := query.Dom(dom2[0], "div[class=markdown-body] p")
			re := regexp.MustCompile(`^\d+\.`)
			if !nondigit && !re.MatchString(release) {
				next = true
			} else {
				revision.Version = release
				t, _ := time.Parse(timeForm, date)
				revision.Date = t
				if len(messages) > 0 {
					revision.Message = strings.TrimSpace(messages[0].Text())
				}
			}
		}
		// commit release
		if len(dom1) == 0 || next {
			dom2 := query.Dom(query.URI(uri+"/releases"), "div[class*=release-timeline-tags] div[class=release-entry] div[class*=d-flex]")
			release := strings.TrimSpace(query.Dom(dom2[0], "h4[class*=commit-title] a")[0].Text())
			date, _ := query.Dom(dom2[0], "span[class*=tag-timeline-date] relative-time")[0].Attr("datetime")
			messages := query.Dom(dom2[0], "div[class=commit-desc] pre")
			revision.Version = release
			t, _ := time.Parse("2006-01-02T15:04:05Z", date)
			revision.Date = t
			if len(messages) > 0 {
				revision.Message = strings.TrimSpace(messages[0].Text())
			}
		}
	} else {
		unstable = true
	}

	if unstable || genChange && len(revision.Message) == 0 {
		//FIXME: 1. can't navigate to next commit page
		dom1 := query.URI(uri + "/commits/master")

		lastModification := c.ModificationTime()
		releaseTime := revision.Date

		dom1.Find("div[class=commit-group-title]").Each(func(i int, selection *goquery.Selection) {
			d := strings.Replace(strings.TrimSpace(selection.Text()), "Commits on ", "", -1)
			t, _ := time.Parse(timeForm, d)

			if t.After(lastModification) {
				if i == 0 {
					if len(revision.Version) == 0 {
						revision.Version = "0.0.0"
					}
					if unstable {
						revision.Date = t
						shas := query.Dom(dom1.Find("ol[class*=commit-group]").Eq(i), "div[class*=commit-links-group]>a")
						revision.Sha = strings.TrimSpace(shas[0].Text())
					}
				}
				if genChange {
					// for release without changelog, only allow commit messages before the release time.
					if releaseTime.IsZero() || !releaseTime.IsZero() && t.Before(releaseTime) {
						messages := query.Dom(dom1.Find("ol[class*=commit-group]").Eq(i), "p[class*=commit-title]>a")
						revision.Message += "* " + strings.TrimSpace(messages[0].Text()) + "\n"
					}
				}
			}
		})
		return revision, nil
	}
	return commit.Commit{}, nil
}
