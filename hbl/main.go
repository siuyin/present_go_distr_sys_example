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
	tkr := time.Tick(time.Second)
	id := ""
	for id == "" {
		c.Request(cfg.IDOffice, "May I an ID please?", &id, time.Second)
	}
	me := cfg.NRS{Name: "Herbert", Rank: cfg.HBListener, ID: id}
	//010_OMIT
	log.Println("Listener Starting...")

	c.Subscribe(cfg.HeartBeat, func(agent *nrs) {
		log.Printf("%v", agent)
	})
	//020_OMIT

	for {
		select {
		case <-tkr:
			c.Publish(cfg.HeartBeat, me)
		}
	}
}
