package main

import (
	"log"
	"time"

	"siuyin/junk/nats/exampleA/cfg"

	"github.com/nats-io/go-nats"
	"github.com/siuyin/dflt"
	"github.com/siuyin/watch"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()

	myID := ""
	for myID == "" {
		c.Request(cfg.IDOffice, "I'd like an ID please.", &myID, time.Second)
	}
	log.Printf("My ID is %v", myID)

	//010_OMIT
	name := dflt.EnvString("NAME", "FileWatcher1")
	me := cfg.NRS{Name: name, Rank: cfg.FileWatcher, ID: myID}

	monPath := dflt.EnvString("MONPATH", ".")
	w := watch.NewWatcher(monPath, time.Second, 3*time.Second)
	wt := w.Watch()
	//020_OMIT

	log.Printf("%s Starting...", name)
	tkr := time.Tick(time.Second)
	for {
		//030_OMIT
		select {
		case <-tkr:
			c.Publish(cfg.HeartBeat, me)
		case f := <-wt: // f is a string // HL
			c.Publish(cfg.StableFilesA, f)
		}
		//040_OMIT
	}
}
