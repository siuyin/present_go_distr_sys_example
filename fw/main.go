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

	myID := cfg.GetID(c)
	log.Printf("My ID is %v", myID)

	//010_OMIT
	name := dflt.EnvString("NAME", "FileWatcher1")
	me := &cfg.NRS{Name: name, Rank: cfg.FileWatcher, ID: myID,
		Tx: []cfg.Board{cfg.StableFilesA}}
	cfg.SendHeartBeat(c, me)

	monPath := dflt.EnvString("MONPATH", "./junk")
	w := watch.NewWatcher(monPath, time.Second, 3*time.Second)
	wt := w.Watch()
	//020_OMIT

	wd := getWorkingDirectory()

	log.Printf("%s watching %s.", name, monPath)
MAINLOOP:
	for {
		select {
		//030_OMIT
		case f := <-wt: // f is a string // HL
			fi, err := os.Stat(f)
			if err != nil {
				log.Println(err)
				continue MAINLOOP
			}
			sendFileDetails(c, wd, f, fi, me)
			//040_OMIT
		}
	}
}

func getWorkingDirectory() string {
	var (
		wd  string
		err error
	)
	for wd, err = os.Getwd(); err != nil; wd, err = os.Getwd() {
		log.Println(err)
		time.Sleep(time.Second)
	}
	return wd
}
func sendFileDetails(c *nats.EncodedConn, wd, fn string, fi os.FileInfo, me *cfg.NRS) {
	fd := mun.FileDetails{
		WorkingDirectory: wd, FileName: fn, FileWatcher: *me}
	fd.IsDir = fi.IsDir()
	fd.Size = fi.Size()
	fd.ModTime = fi.ModTime()
	c.Publish(string(cfg.StableFilesA), &fd)
}
