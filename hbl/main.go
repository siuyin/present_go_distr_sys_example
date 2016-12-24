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
	id := cfg.GetID(c)
	me := &cfg.NRS{Name: "Herbert", Rank: cfg.HBListener, ID: id,
		Rx: []cfg.Board{cfg.HeartBeat}, Tx: []cfg.Board{cfg.HeartBeat}}
	cfg.SendHeartBeat(c, me)
	seen := map[string]dat{}
	//010_OMIT
	log.Println("Listener Starting...")

	c.Subscribe(string(cfg.HeartBeat), func(agent *dat) {
		agent.T = time.Now()
		key := agent.Name + string(agent.Rank) + agent.ID
		mtx.Lock()
		seen[key] = *agent
		mtx.Unlock()
	})
	//020_OMIT

	for {
		select {
		case <-tkr:
			c.Publish(string(cfg.HeartBeat), me)
			displayDat(&seen)
		}
	}
}

func displayDat(d *map[string]dat) {
	mtx.Lock()
	keys := make([]string, len(*d))
	mtx.Unlock()
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
		fmt.Printf("%s%s %s %s %s\n  Rx:%s\n  Tx:%s\n", s, v.Name, v.Rank, v.ID,
			v.T.Format("15:04:05 MST"), v.RxList(), v.TxList())
	}
}
