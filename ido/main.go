package main

// Generates base32 encoded random IDs.

import (
	"crypto/rand"
	"encoding/base32"
	"log"

	"siuyin/junk/nats/exampleA/cfg"

	"github.com/nats-io/go-nats"
	"github.com/siuyin/dflt"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()

	name := dflt.EnvString("NAME", "IDOfc1")
	id := dflt.EnvString("ID", "001")
	me := &cfg.NRS{Name: name, Rank: cfg.IDOfficer, ID: id,
		Rx: []cfg.Board{cfg.IDOffice}}
	cfg.SendHeartBeat(c, me)
	//010_OMIT
	log.Printf("%s %s Starting...\n", me.Name, me.ID)

	c.Subscribe(string(cfg.IDOffice), func(subj, reply string, req *string) {
		c.Publish(reply, randID())
	})
	//020_OMIT
	select {}
}

func randID() string {
	c := 5
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("error:", err)
		return ""
	}
	// The slice should now contain random bytes instead of only zeroes.
	str := base32.StdEncoding.EncodeToString(b)
	return str
}
