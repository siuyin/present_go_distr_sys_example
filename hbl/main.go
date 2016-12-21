package main

import (
	"log"

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
	log.Println("Listener Starting...")

	c.Subscribe(cfg.HeartBeat, func(agent *nrs) {
		log.Printf("%v", agent)
	})
	//020_OMIT

	select {}
}
