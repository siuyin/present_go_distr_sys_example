package mun

// Configuration for the manager kitMun.

import (
	"siuyin/junk/nats/exampleA/cfg"
	"time"
)

// FileDetails is the message passed by FileWatchers to Manager kitMun.
type FileDetails struct {
	WorkingDirectory string
	FileName         string
	IsDir            bool
	Size             int64
	ModTime          time.Time
	FileWatcher      cfg.NRS
}

type noMoreWorkErr struct{}

func (e noMoreWorkErr) Error() string {
	return "No More Work"
}

// NoMoreWorkError signals no more work
var NoMoreWorkError noMoreWorkErr
