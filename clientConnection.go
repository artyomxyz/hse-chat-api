package main

import "github.com/hse-chat/hse-chat-api/HseMsg"

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
	serverMessage := &HseMsg.ServerMessage{
		Message: &HseMsg.ServerMessage_Event{
			Event: &HseMsg.Event{
				Event: &HseMsg.Event_NewMessage_{
					NewMessage: &HseMsg.Event_NewMessage{
						Message: &HseMsg.Message{
							Author:   &msg.author,
							Receiver: &msg.receiver,
							Date:     &msg.date,
							Text:     &msg.text,
						},
					},
				},
			},
		},
	}
	err := clConn.conn.Write(serverMessage)
	if err != nil {
		return err
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

	msgEmitter.AddListener(clConn.msgChan)

	go clConn.receiveRequests()
	go clConn.process()

	return clConn
}
