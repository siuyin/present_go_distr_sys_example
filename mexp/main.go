package main

// math expert who knows how to add

import (
	"fmt"
	"log"
	"time"

	"siuyin/junk/nats/exampleA/cfg"
	"siuyin/junk/nats/exampleA/msolv/msh"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()

	myID := cfg.GetID(c) // will block until IDOffice is open
	log.Printf("My ID is %v", myID)

	me := &cfg.NRS{Name: "MathExpert1", Rank: cfg.MathExpert, ID: myID,
		Rx: []cfg.Board{cfg.MathProblemsA}, Tx: []cfg.Board{cfg.MathSolversAOut}}
	cfg.SendHeartBeat(c, me)
	log.Printf("MathExpert %s Starting...", me.Name)
	//010_OMIT
	c.Subscribe(string(cfg.MathProblemsA), func(mp *msh.MathProblem) {
		myAns := msh.MathAnswer{
			SolverID:   me.ID,
			ProblemID:  mp.ID,
			AnswerID:   cfg.GetID(c),
			Answer:     []byte("42"),
			AnswerTime: time.Now(),
		}
		c.Publish(string(cfg.MathSolversAOut), myAns)
		fmt.Printf("Sent answer: %s for Problem: %s\n", myAns.Answer, myAns.ProblemID)
	})
	//020_OMIT
	select {}
}
