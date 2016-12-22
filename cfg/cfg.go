package cfg

// IMPORTANT -- ALL entries in this file are append only.
//
// CHANGING / EDITING entries require retiring the generation that used the old config.
// Think DNA. Our DNA is mostly the result responses to past pathogen attacks.
//
// Once a response is coded it stays in our DNA

import (
	"time"

	"github.com/nats-io/go-nats"
)

// Application Constants and Service Names
const (
	App       = "EgA"
	IDOffice  = App + ".IDOffice"
	HeartBeat = App + ".HeartBeat"
)

// Ranks or Position or Function
const (
	Unassigned  = App + ".Unassigned"
	FileWatcher = App + ".FileWatcher"
	IDOfficer   = App + ".IDOfficer"
	HBListener  = App + ".HeartBeatListener"
	ManagerA    = App + ".ManagerA"
)

// Stable Files Pool
const (
	StableFilesA = App + ".StableFilesPool.A"
)

// NRS Name, Rank and Serial Number (ID)
type NRS struct {
	Name, Rank, ID string
}

//GetID requests an ID from the ID Office.
func GetID(c *nats.EncodedConn) string {
	id := ""
	for id == "" {
		c.Request(IDOffice, "I'd like an ID please.", &id, time.Second)
	}
	return id
}

// SendHeartBeat sends a heart beat to the HeartBeat endpoint.
func SendHeartBeat(c *nats.EncodedConn, me *NRS) {
	tkr := time.Tick(time.Second)
	go func() {
		for {
			select {
			case <-tkr:
				c.Publish(HeartBeat, me)
			}
		}
	}()
}
