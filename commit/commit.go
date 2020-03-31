package commit

import "time"

// Commit an universal commit struct
type Commit struct {
	Version string
	Date    time.Time
	Message string
	Sha     string
}

// UnstableVersion return the unstable version in string
func (c Commit) UnstableVersion() string {
	return c.Version + "+git" + c.Date.Format("20060102") + "." + c.Sha
}
