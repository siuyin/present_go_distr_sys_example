package cfg

// IMPORTANT -- ALL entries in this file are append only.
//
// CHANGING / EDITING entries require retiring the generation that used the old config.
// Think DNA. Our DNA is mostly the result responses to past pathogen attacks.
//
// Once a response is coded it stays in our DNA

import (
	"strings"
	"time"

	"github.com/nats-io/go-nats"
)

// Application Constants and Service Names
const (
	App = "EgA"
)

// Rank is a rank or position within the organisation / system.
type Rank string

// Ranks or Position or Function
const (
	Unassigned  = Rank(App + ".Unassigned")
	FileWatcher = Rank(App + ".FileWatcher")
	IDOfficer   = Rank(App + ".IDOfficer")
	HBListener  = Rank(App + ".HeartBeatListener") // they listen to agent heart beats
	ManagerA    = Rank(App + ".ManagerA")          // Example manager to manage files
	MathSolver  = Rank(App + ".MathSolver")
	MathExpert  = Rank(App + ".MathExpert") // they hang out at the MathSolvers bulletin board
	FileMover   = Rank(App + ".FileMover")
)

// Board is a Bulletin Board or Radio Frequency Channel.
type Board string

// Bulletin Boards or Frequency-- agents post their messages here.
const (
	StableFilesA    = Board(App + ".StableFilesPool.A")
	MathProblemsA   = Board(App + ".MathProblems.A")
	MathSolversAOut = Board(App + ".MathSolvers.A.Outbox")
	FileMoversA     = Board(App + ".FileMovers.A")
	FileMoversAOut  = Board(App + ".FileMovers.A.Outbox")
	HeartBeat       = Board(App + ".HeartBeat")
	IDOffice        = Board(App + ".IDOffice")
)

// NRS Name, Rank and Serial Number (ID)
type NRS struct {
	Name   string
	Rank   Rank
	ID     string
	Tx, Rx []Board
}

//RxList returns a comma separated list of Rx Board names
func (n NRS) RxList() string {
	s := []string{}
	for _, v := range n.Rx {
		s = append(s, string(v))
	}
	return strings.Join(s, ",")
}

//TxList returns a comma separated list of Tx Board names
func (n NRS) TxList() string {
	s := []string{}
	for _, v := range n.Tx {
		s = append(s, string(v))
	}
	return strings.Join(s, ",")
}

//GetID requests an ID from the ID Office.
func GetID(c *nats.EncodedConn) string {
	id := ""
	for id == "" {
		c.Request(string(IDOffice), "I'd like an ID please.", &id, time.Second)
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
				c.Publish(string(HeartBeat), me)
			}
		}
	}()
}
