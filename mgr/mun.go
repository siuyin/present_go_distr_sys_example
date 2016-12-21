package main

// Manager Mun

import (
	"log"
	"time"

	"siuyin/junk/nats/exampleA/cfg"
	"siuyin/junk/nats/exampleA/mgr/mun"

	"github.com/nats-io/go-nats"
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

	me := cfg.NRS{Name: "Mun", Rank: cfg.ManagerA, ID: myID}
	log.Printf("Manager %s Starting...", me.Name)

	//010_OMIT
	c.Subscribe(cfg.StableFilesA, func(subj, reply string, fd *mun.FileDetails) {
		log.Printf("Received from %s f:%s d:%s D:%v s:%v t:%s",
			subj, fd.FileName, fd.WorkingDirectory, fd.IsDir, fd.Size, fd.ModTime)
	})
	//020_OMIT

	tkr := time.Tick(time.Second)
	for {
		select {
		case <-tkr:
			c.Publish(cfg.HeartBeat, me)
		}
	}
}
