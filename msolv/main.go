package main

// Extensible Manager MathSolver Marsha
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"siuyin/junk/nats/exampleA/cfg"
	"siuyin/junk/nats/exampleA/msolv/msh"

	"github.com/nats-io/go-nats"
)

type work struct {
	Problem    msh.MathProblem
	ReceivedAt time.Time
	Done       bool
	Answers    []msh.MathAnswer
}

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()

	mtx := sync.Mutex{}
	jobs := map[string]*work{}
	myID := cfg.GetID(c) // will block until IDOffice is open
	log.Printf("My ID is %v", myID)

	me := &cfg.NRS{Name: "Marsha", Rank: cfg.MathSolver, ID: myID,
		Rx: []cfg.Board{cfg.MathProblemsA, cfg.MathSolversAOut},
		Tx: []cfg.Board{cfg.MathProblemsA}}
	cfg.SendHeartBeat(c, me)
	log.Printf("MathSolver %s Starting...", me.Name)

	//010_OMIT
	c.Subscribe(string(cfg.MathProblemsA), func(mp *msh.MathProblem) {
		addToJobs(&jobs, &mtx, mp)
		// No need to inform Experts (workers) as they already listen in.
	})

	c.Subscribe(string(cfg.MathSolversAOut), func(ans *msh.MathAnswer) {
		updateJobs(&jobs, &mtx, ans)
	})
	//020_OMIT
	selfTest(c)
	flagUnsolvedJobs(&jobs, &mtx)
	flagMultipleAnswerJobs(&jobs, &mtx)
	showAnswer(&jobs, &mtx)
	//listJobs(&jobs, &mtx)

	select {}
}

func selfTest(c *nats.EncodedConn) {
	myProb := struct {
		Op   string
		Data []int
	}{"Add", []int{2, 3}}
	dat, err := json.Marshal(myProb)
	if err != nil {
		log.Fatal(err)
	}
	tkr := time.Tick(1500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-tkr:
				mp := msh.MathProblem{Name: "2 + 3", ID: cfg.GetID(c), Data: dat}
				c.Publish(string(cfg.MathProblemsA), mp)
			}
		}
	}()
}

func showAnswer(j *map[string]*work, m *sync.Mutex) {
	go func() {
		tkr := time.Tick(1 * time.Second)
		for {
			select {
			case <-tkr:
				m.Lock()
				for k, v := range *j {
					now := time.Now()
					if !v.Done && len(v.Answers) == 1 {
						fmt.Printf("%s ID:%s solved by %s/%s in %.6f ans: %s\n", now.Format("05.000000"),
							k, v.Answers[0].SolverID, v.Answers[0].AnswerID, v.Answers[0].AnswerTime.Sub(v.ReceivedAt).Seconds(),
							v.Answers[0].Answer)
						delete(*j, k) // optional
					}
					//v.Done = true
				}
				m.Unlock()
			}
		}
	}()
}

func addToJobs(jobs *map[string]*work, mtx *sync.Mutex, mp *msh.MathProblem) {
	mtx.Lock()
	(*jobs)[mp.ID] = &work{ReceivedAt: time.Now(), Problem: *mp}
	mtx.Unlock()
}
func updateJobs(jobs *map[string]*work, mtx *sync.Mutex, ans *msh.MathAnswer) {
	mtx.Lock()
	w := (*jobs)[ans.ProblemID]
	w.Answers = append(w.Answers, *ans)
	mtx.Unlock()
}
func flagUnsolvedJobs(j *map[string]*work, m *sync.Mutex) {
	go func() {
		tkr := time.Tick(1 * time.Second)
		for {
			select {
			case <-tkr:
				m.Lock()
				for k, v := range *j {
					now := time.Now()
					if now.Sub(v.ReceivedAt) > 100*time.Millisecond && len(v.Answers) == 0 {
						fmt.Printf("%s ID:%s not solved for %.6f\n", now.Format("05.000000"),
							k, now.Sub(v.ReceivedAt).Seconds())
						delete(*j, k)
					}
				}
				m.Unlock()
			}
		}
	}()
}
func flagMultipleAnswerJobs(j *map[string]*work, m *sync.Mutex) {
	go func() {
		tkr := time.Tick(1 * time.Second)
		for {
			select {
			case <-tkr:
				m.Lock()
				for k, v := range *j {
					now := time.Now()
					if now.Sub(v.ReceivedAt) > 10*time.Millisecond && len(v.Answers) > 1 {
						a0 := v.Answers[0].Answer
						unanimous := true
						for _, v := range v.Answers[1:] {
							if !bytes.Equal(a0, v.Answer) {
								unanimous = false
								break
							}
						}
						if unanimous {
							fmt.Printf("ID: %s Unanimous answer from %d: %s\n", k, len(v.Answers), a0)
							delete(*j, k)
						} else {
							fmt.Printf("%s ID:%s multiple differing Ans: %d\n", now.Format("05.000000"),
								k, len(v.Answers))
							delete(*j, k)
						}
					}
				}
				m.Unlock()
			}
		}
	}()
}
func listJobs(j *map[string]*work, m *sync.Mutex) {
	go func() {
		tkr := time.Tick(3 * time.Second)
		for {
			select {
			case <-tkr:
				m.Lock()
				for k, v := range *j {
					fmt.Printf("ID:%s, r:%s s:%t\n", k,
						v.ReceivedAt.Format("15:04:05"), v.Done)
				}
				m.Unlock()
			}
		}
	}()
}
