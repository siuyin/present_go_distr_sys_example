package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"siuyin/junk/nats/exampleA/cfg"

	"github.com/nats-io/go-nats"
)

type dat struct {
	cfg.NRS
	T time.Time
}

var mtx sync.Mutex

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
		mtx.Lock()
		seen[key] = *agent
		mtx.Unlock()
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

	mtx.Lock()
	for k := range *d {
		keys[i] = k
		i++
	}
	mtx.Unlock()

	sort.Strings(keys)
	fmt.Println("=================================")
	for _, k := range keys {
		mtx.Lock()
		v := (*d)[k]
		mtx.Unlock()
		s := ""
		if time.Now().Sub(v.T).Seconds() > 2 {
			s = "F: "
		}
		fmt.Printf("%s%s %s %s\n", s, v.Name, v.Rank, v.T.Format("15:04:05 MST"))
	}
}
