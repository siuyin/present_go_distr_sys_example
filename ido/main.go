package main

// Generates base32 encoded random IDs.

import (
	"crypto/rand"
	"encoding/base32"
	"log"

	"github.com/nats-io/go-nats"
)

//010_OMIT
func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()
	log.Println("ID Issuer Starting...")

	c.Subscribe("EgA.IDOffice", func(subj, reply string, req *string) {
		c.Publish(reply, randID())
	})

	select {}
}

//020_OMIT

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
