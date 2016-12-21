package main

import (
	"log"
	"os"
	"time"

	"siuyin/junk/nats/exampleA/cfg"
	"siuyin/junk/nats/exampleA/mgr/mun"

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

	var (
		wd  string
		err error
	)
	for wd, err = os.Getwd(); err != nil; wd, err = os.Getwd() {
		log.Println(err)
		time.Sleep(time.Second)
	}

	log.Printf("%s Starting...", name)
	tkr := time.Tick(time.Second)
MAINLOOP:
	for {
		select {
		case <-tkr:
			c.Publish(cfg.HeartBeat, me)
		case f := <-wt: // f is a string // HL
			fi, err := os.Stat(f)
			if err != nil {
				log.Println(err)
				continue MAINLOOP
			}
			//030_OMIT
			fd := mun.FileDetails{
				WorkingDirectory: wd, FileName: f, FileWatcher: me}
			fd.IsDir = fi.IsDir()
			fd.Size = fi.Size()
			fd.ModTime = fi.ModTime()
			c.Publish(cfg.StableFilesA, &fd)
			//040_OMIT
		}
	}
}
