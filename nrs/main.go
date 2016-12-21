package main

import (
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

type nrs struct {
	Name, Rank, ID string
}

//010_OMIT
func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()
	me := nrs{Name: "NameA", Rank: "EgA.Unassigned", ID: "CXHSD"}
	log.Println("Service Starting...")
	tkr := time.Tick(time.Second)
	for {
		select {
		case <-tkr:
			c.Publish("EgA.HeartBeat", me)
		}
	}
}

//020_OMIT
