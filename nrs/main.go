package main

import (
	"log"

	"siuyin/junk/nats/exampleA/cfg"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()

	//010_OMIT
	myID := cfg.GetID(c) // will block until IDOffice is open
	log.Printf("My ID is %v", myID)

	me := &cfg.NRS{Name: "NameA", Rank: cfg.Unassigned, ID: myID}
	cfg.SendHeartBeat(c, me)
	log.Println("Service Starting...")
	select {}
	//020_OMIT
}
