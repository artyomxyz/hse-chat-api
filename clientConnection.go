package main

import (
	"log"

	"github.com/hse-chat/hse-chat-api/HseMsg"
)

// ClientConnection represents client connection
type ClientConnection struct {
	conn    ProtoConnection
	client  *Client
	reqChan chan *HseMsg.Request
	msgChan chan Message
}

func (clConn ClientConnection) process() {
	for {
		select {
		case req := <-clConn.reqChan:
			err := clConn.handleRequest(req)
			if err != nil {
				clConn.close()
			}
		case msg := <-clConn.msgChan:
			err := clConn.handleNewMessage(msg)
			if err != nil {
				clConn.close()
			}
		}
	}
}

func (clConn ClientConnection) handleNewMessage(msg Message) error {
	if clConn.client.CanReadMessage(msg) {
		err := clConn.conn.Write(msg.ToServerMessage())
		if err != nil {
			return err
		}
	}
	return nil
}

func (clConn ClientConnection) close() {
	clConn.conn.Close()
}

func (clConn ClientConnection) receiveRequests() {
	defer clConn.close()

	for {
		req := &HseMsg.Request{}
		err := clConn.conn.Read(req)
		if err != nil {
			return
		}
		clConn.reqChan <- req
	}
}

func (clConn ClientConnection) handleRequest(req *HseMsg.Request) error {
	var err error
	log.Printf("Received %s", req)

	var res Result

	if signUp := req.GetSignUp(); signUp != nil {
		res, err = clConn.client.SignUp(signUp.GetUsername(), signUp.GetPassword())
	}

	if signIn := req.GetSignIn(); signIn != nil {
		res, err = clConn.client.SignIn(signIn.GetUsername(), signIn.GetPassword())
	}

	if getUsers := req.GetGetUsers(); getUsers != nil {
		res, err = clConn.client.GetUsers()
	}

	if getMessagesWithUser := req.GetGetMessagesWithUser(); getMessagesWithUser != nil {
		res, err = clConn.client.GetMessagesWithUser(getMessagesWithUser.GetWith())
	}

	if sendMessageToUser := req.GetSendMessageToUser(); sendMessageToUser != nil {
		res, err = clConn.client.SendMessageToUser(sendMessageToUser.GetReceiver(), sendMessageToUser.GetText())
	}

	if err != nil {
		return err
	}

	if res != nil {
		serverMessage := res.ToServerMessage(req.GetId())
		err = clConn.conn.Write(serverMessage)
		if err != nil {
			return err
		}

		log.Printf("Sent %s", serverMessage)
	}

	return nil
}

// NewClientConnection creates new client connection and start listening on it
func NewClientConnection(conn ProtoConnection) ClientConnection {
	client := NewClient()
	clConn := ClientConnection{
		conn,
		&client,
		make(chan *HseMsg.Request),
		make(chan Message),
	}

	msgMngr.AddListener(clConn.msgChan)

	go clConn.receiveRequests()
	go clConn.process()

	return clConn
}
