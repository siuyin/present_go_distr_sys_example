package main

import (
	"log"
	"time"

	"siuyin/junk/nats/exampleA/cfg"

	"github.com/nats-io/go-nats"
)

type nrs struct {
	Name, Rank, ID string
}

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()

	//010_OMIT
	myID := ""
	for myID == "" {
		c.Request(cfg.IDOffice, "I'd like an ID please.", &myID, time.Second)
	}
	log.Printf("My ID is %v", myID)

	me := cfg.NRS{Name: "NameA", Rank: cfg.Unassigned, ID: myID}
	log.Println("Service Starting...")
	tkr := time.Tick(time.Second)
	for {
		select {
		case <-tkr:
			c.Publish(cfg.HeartBeat, me)
		}
	}
	//020_OMIT
}
