package mun

// Configuration for the manager Mun.

import (
	"siuyin/junk/nats/exampleA/cfg"
	"time"
)

// Constants used my manager Mun to direct FileMovers.
const (
	FileMove = iota + 1
	FileCopy
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

// FileMoveCmd is the command type sent to FileMover workers.
type FileMoveCmd struct {
	ID       string
	Op       int
	From, To string
}

// FileMoveStatus is the type used by FileMover responses.
type FileMoveStatus struct {
	ResponseID,
	CmdID,
	Status string
}

type noMoreWorkErr struct{}

func (e noMoreWorkErr) Error() string {
	return "No More Work"
}

// NoMoreWorkError signals no more work
var NoMoreWorkError noMoreWorkErr
