package common

import "time"

// Commit an universal commit struct
type Commit struct {
	Date    time.Time
	Message string
	Sha     string
}
