package main

// Data Mover
import (
	"fmt"
	"log"
	"os"

	"siuyin/junk/nats/exampleA/cfg"
	"siuyin/junk/nats/exampleA/mgr/mun"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()

	myID := cfg.GetID(c) // will block until IDOffice is open
	log.Printf("My ID is %v", myID)

	me := &cfg.NRS{Name: "FileMover1", Rank: cfg.FileMover, ID: myID,
		Tx: []cfg.Board{cfg.FileMoversAOut}, Rx: []cfg.Board{cfg.FileMoversA}}
	cfg.SendHeartBeat(c, me)
	log.Printf("%s Starting...", me.Name)

	//010_OMIT
	c.Subscribe(string(cfg.FileMoversA), func(cmd *mun.FileMoveCmd) {
		r := mun.FileMoveStatus{ResponseID: cfg.GetID(c),
			CmdID: cmd.ID, Status: "Ack"}
		c.Publish(string(cfg.FileMoversAOut), r)
		fmt.Println(r)
		// perform the move and update FileMoversA bulletin board
		//020_OMIT

		if _, err := os.Stat(cmd.From); err != nil {
			r.ResponseID = cfg.GetID(c)
			r.Status = fmt.Sprintf("From Path: Error: %v", err)
			c.Publish(string(cfg.FileMoversAOut), r)
			fmt.Println(r.Status)
			return
		}

		r.ResponseID = cfg.GetID(c)
		r.Status = fmt.Sprintf("From Path: %s: found.", cmd.From)
		c.Publish(string(cfg.FileMoversAOut), r)
		fmt.Println(r.Status)

		if _, err := os.Stat(cmd.To); err != nil {
			r.ResponseID = cfg.GetID(c)
			r.Status = fmt.Sprintf("To Path: Error: %v", err)
			c.Publish(string(cfg.FileMoversAOut), r)
			fmt.Println(r.Status)
			return
		}

		r.ResponseID = cfg.GetID(c)
		r.Status = fmt.Sprintf("To Path: %s: found.", cmd.To)
		c.Publish(string(cfg.FileMoversAOut), r)
		fmt.Println(r.Status)

		// TODO
		// switch cmd.Op {
		// case mun.FileCopy:
		// 	b, err := ioutil.ReadFile(cmd.From)
		// 	if err != nil {
		// 		log.Println(err)
		// 		return
		// 	}
		// 	err = ioutil.WriteFile(cmd.To, b, os.ModePerm)
		// 	if err != nil {
		// 		log.Println(err)
		// 		return
		// 	}
		// case mun.FileMove:
		// 	os.Rename(cmd.From, cmd.To)
		// }
	})
	select {}
}
