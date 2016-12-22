package main

// Generates base32 encoded random IDs.

import (
	"crypto/rand"
	"encoding/base32"
	"log"

	"siuyin/junk/nats/exampleA/cfg"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()

	me := &cfg.NRS{Name: "IDOfc1", Rank: cfg.IDOfficer, ID: "001"}
	cfg.SendHeartBeat(c, me)
	//010_OMIT
	log.Println("ID Issuer Starting...")

	c.Subscribe(cfg.IDOffice, func(subj, reply string, req *string) {
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
