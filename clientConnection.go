package main

import "github.com/hse-chat/hse-chat-api/HseMsg"

// ClientConnection represents client connection
type ClientConnection struct {
	conn     ProtoConnection
	client   *Client
	reqChan  chan *HseMsg.Request
	evtChan  chan Event
	doneChan chan bool
}

func (clConn ClientConnection) process() {
	defer clConn.close()
	for {
		select {
		case <-clConn.doneChan:
			return
		case req := <-clConn.reqChan:
			err := clConn.handleRequest(req)
			if err != nil {
				clConn.doneChan <- true
			}
		case evt := <-clConn.evtChan:
			err := clConn.handleEvent(evt)
			if err != nil {
				clConn.doneChan <- true
			}
		}
	}
}

func (clConn ClientConnection) handleEvent(evt Event) error {
	if evt.IsAccessibleBy(clConn.client.user) {
		return clConn.conn.Write(
			EventToServerMessage(evt),
		)
	}
	return nil
}

func (clConn ClientConnection) close() {
	close(clConn.doneChan)
	close(clConn.reqChan)
	evtMngr.RemoveListener(clConn.evtChan)
	close(clConn.evtChan)
	clConn.conn.Close()
}

func (clConn ClientConnection) receiveRequests() {
	for {
		req := &HseMsg.Request{}
		err := clConn.conn.Read(req)
		if err != nil {
			break
		}

		clConn.reqChan <- req
	}

	clConn.doneChan <- true
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
		make(chan Event),
		make(chan bool),
	}

	evtMngr.AddListener(clConn.evtChan)

	go clConn.receiveRequests()
	go clConn.process()

	return clConn
}
