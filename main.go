package main

import (
	"log"

	"gopkg.in/mgo.v2"
)

var db *mgo.Database
var msgMngr *MessageManager

func main() {
	msgMngr = NewMessageManager()

	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	log.Print("Connected to DB")

	db = session.DB("chat")

	err = db.C("users").EnsureIndex(mgo.Index{
		Key:      []string{"username"},
		Unique:   true,
		DropDups: true,
	})

	if err != nil {
		panic(err)
	}

	log.Print("Ensured username index in db")

	listener, err := ProtoListen("tcp", ":8080")

	if err != nil {
		panic(err)
	}

	log.Print("Listening port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
		}
		NewClientConnection(conn)
	}
}
