package main

import (
	"log"

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
	log.Println("Listener Starting...")

	c.Subscribe("EgA.HeartBeat", func(agent *nrs) {
		log.Printf("%v", agent)
	})

	select {}
}

//020_OMIT
