package models

import "time"

type MergeRequests struct {
	Length        int
	On            time.Time
	MergeRequests []MergeRequestFileChanges
}
