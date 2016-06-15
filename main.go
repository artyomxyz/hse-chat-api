package main

import (
	"log"

	"github.com/hse-chat/hse-chat-api/HseMsg"

	"gopkg.in/mgo.v2"
)

var db *mgo.Database

func handleConnection(conn ProtoConnection) {
	log.Printf("New connection")
	defer conn.Close()
	defer log.Print("Disconnected")

	cl := NewClient()

	for {
		req := &HseMsg.Request{}
		err := conn.Read(req)
		if err != nil {
			return
		}

		log.Printf("Received %s", req)

		var res Result

		if signUp := req.GetSignUp(); signUp != nil {
			res, err = cl.SignUp(signUp.GetUsername(), signUp.GetPassword())
		}

		if signIn := req.GetSignIn(); signIn != nil {
			res, err = cl.SignIn(signIn.GetUsername(), signIn.GetPassword())
		}

		if getUsers := req.GetGetUsers(); getUsers != nil {
			res, err = cl.GetUsers()
		}

		if getMessagesWithUser := req.GetGetMessagesWithUser(); getMessagesWithUser != nil {
			res, err = cl.GetMessagesWithUser(getMessagesWithUser.GetWith())
		}

		if sendMessageToUser := req.GetSendMessageToUser(); sendMessageToUser != nil {
			res, err = cl.SendMessageToUser(sendMessageToUser.GetReceiver(), sendMessageToUser.GetText())
		}

		if err != nil {
			return
		}

		if res != nil {
			serverMessage := res.ToServerMessage(req.GetId())
			err = conn.Write(serverMessage)
			if err != nil {
				return
			}

			log.Printf("Sent %s", serverMessage)
		}
	}
}

func main() {
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
			// close
		}
		go handleConnection(conn)
	}
}
