package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"siuyin/junk/nats/exampleA/cfg"

	"github.com/nats-io/go-nats"
)

type dat struct {
	cfg.NRS
	T time.Time
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
	seen := map[string]dat{}
	//010_OMIT
	log.Println("Listener Starting...")

	c.Subscribe(cfg.HeartBeat, func(agent *dat) {
		agent.T = time.Now()
		key := agent.Name + agent.Rank
		seen[key] = *agent
	})
	//020_OMIT

	for {
		select {
		case <-tkr:
			c.Publish(cfg.HeartBeat, me)
			displayDat(&seen)
		}
	}
}

func displayDat(d *map[string]dat) {
	keys := make([]string, len(*d))
	i := 0
	for k := range *d {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	fmt.Println("=================================")
	for _, k := range keys {
		v := (*d)[k]
		fmt.Printf("%s %s %s\n", v.Name, v.Rank, v.T.Format("15:04:05 MST"))
	}
}
