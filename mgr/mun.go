package main

// Manager Mun

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"siuyin/junk/nats/exampleA/cfg"
	//030_OMIT
	"siuyin/junk/nats/exampleA/mgr/mun"
	//040_OMIT

	"github.com/boltdb/bolt"
	"github.com/nats-io/go-nats"
	"github.com/siuyin/dflt"
)

const workQSize = 1000

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json")
	defer c.Close()

	myID := ""
	for myID == "" {
		c.Request(cfg.IDOffice, "I'd like an ID please.", &myID, time.Second)
	}
	log.Printf("My ID is %v", myID)

	me := cfg.NRS{Name: "Mun", Rank: cfg.ManagerA, ID: myID}
	log.Printf("Manager %s Starting...", me.Name)

	initDB() // see file db.go
	defer db.Close()

	//010_OMIT
	workQ := make(chan workItem, workQSize)
	c.Subscribe(cfg.StableFilesA, func(subj, reply string, fd *mun.FileDetails) {
		w := workItem{Subject: subj, Data: *fd}
		workQ <- w
		// log.Printf("Received from %s f:%s d:%s D:%v s:%v t:%s",
		// 	subj, fd.FileName, fd.WorkingDirectory, fd.IsDir, fd.Size, fd.ModTime)
	})
	//020_OMIT

	tkr := time.Tick(time.Second)
	tkr2 := time.Tick(5 * time.Second)
	for {
		select {
		case <-tkr:
			c.Publish(cfg.HeartBeat, me)
		case <-tkr2:
			dumpDB()
		case w := <-workQ:
			saveToDB(w)
		}
	}
}

// =========================================================================
// Database for Manager Mun

type workItem struct {
	Subject string
	Data    mun.FileDetails
}

const timeFmt = "2006-01-02 15:04:05.000000"

var db *bolt.DB

// var (
// 	gobEnc    *gob.Encoder
// 	gobEncBuf bytes.Buffer
//
// 	gobDec    *gob.Decoder
// 	gobDecBuf bytes.Buffer
// )

// Remember to defer db.Close() in main.
func initDB() {
	var err error
	dbName := dflt.EnvString("DBNAME", "mgr/mun.db")
	db, err = bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf(`%v: Opening DB.
If you see a timeout:
  lsof mun.db
  killall mun
to kill all processes accessing the DB
I have difficulty with go present not killing the process tree.`, err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Jobs"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	log.Println("db Ready")
}

func saveToDB(w workItem) {
	db.Update(func(tx *bolt.Tx) error {
		byt, err := json.Marshal(w)
		if err != nil {
			log.Println("json encode error: ", err)
			return err
		}

		b := tx.Bucket([]byte("Jobs"))
		err = b.Put([]byte(time.Now().Format(timeFmt)), byt)

		log.Println(w.Data.FileName)

		return err
	})
}

func dumpDB() {
	fmt.Println("DB Dump")
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Jobs"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			wi := workItem{}
			err := json.Unmarshal(v, &wi)
			if err != nil {
				log.Fatal(err)
			}
			if wi.Data.IsDir {
				fmt.Printf("t:%s f:%v, DIR\n", k, wi.Data.FileName)
			} else {
				fmt.Printf("t:%s f:%v, s:%v\n", k, wi.Data.FileName, wi.Data.Size)
			}
		}

		return nil
	})
	fmt.Println("End Dump")
}
