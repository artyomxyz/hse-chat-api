package main

import "github.com/hse-chat/hse-chat-api/HseMsg"

// ClientConnection represents client connection
type ClientConnection struct {
	conn            ProtoConnection
	client          *Client
	requestsChannel chan *HseMsg.Request
	eventChannel    chan Event
	doneChan        chan bool
	api             *API
}

func (clientConnection ClientConnection) process() {
	defer clientConnection.close()
	for {
		select {
		case <-clientConnection.doneChan:
			return
		case req := <-clientConnection.requestsChannel:
			err := clientConnection.handleRequest(req)
			if err != nil {
				clientConnection.doneChan <- true
			}
		case evt := <-clientConnection.eventChannel:
			err := clientConnection.handleEvent(evt)
			if err != nil {
				clientConnection.doneChan <- true
			}
		}
	}
}

func (clientConnection ClientConnection) handleEvent(evt Event) error {
	if evt.IsAccessibleBy(clientConnection.client.user) {
		return clientConnection.conn.Write(
			EventToServerMessage(evt),
		)
	}
	return nil
}

func (clientConnection ClientConnection) close() {
	clientConnection.client.Finish()
	close(clientConnection.doneChan)
	close(clientConnection.requestsChannel)
	clientConnection.api.eventManager.RemoveListener(clientConnection.eventChannel)
	close(clientConnection.eventChannel)
	clientConnection.conn.Close()
}

func (clientConnection ClientConnection) receiveRequests() {
	for {
		req := &HseMsg.Request{}
		err := clientConnection.conn.Read(req)
		if err != nil {
			break
		}

		clientConnection.requestsChannel <- req
	}

	clientConnection.doneChan <- true
}

func (clientConnection ClientConnection) handleRequest(req *HseMsg.Request) error {
	var err error
	var res Result

	if signUp := req.GetSignUp(); signUp != nil {
		res, err = clientConnection.client.SignUp(signUp.GetUsername(), signUp.GetPassword())
	}

	if signIn := req.GetSignIn(); signIn != nil {
		res, err = clientConnection.client.SignIn(signIn.GetUsername(), signIn.GetPassword())
	}

	if getUsers := req.GetGetUsers(); getUsers != nil {
		res, err = clientConnection.client.GetUsers()
	}

	if getMessagesWithUser := req.GetGetMessagesWithUser(); getMessagesWithUser != nil {
		res, err = clientConnection.client.GetMessagesWithUser(getMessagesWithUser.GetWith())
	}

	if sendMessageToUser := req.GetSendMessageToUser(); sendMessageToUser != nil {
		res, err = clientConnection.client.SendMessageToUser(sendMessageToUser.GetReceiver(), sendMessageToUser.GetText())
	}

	if err != nil {
		return err
	}

	if res != nil {
		serverMessage := res.ToServerMessage(req.GetId())
		err = clientConnection.conn.Write(serverMessage)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewClientConnection creates new client connection and start listening on it
func NewClientConnection(conn ProtoConnection, api *API) ClientConnection {
	client := NewClient(api.userManager, api.messageManager)
	clientConnection := ClientConnection{
		conn,
		&client,
		make(chan *HseMsg.Request),
		make(chan Event),
		make(chan bool),
		api,
	}

	api.eventManager.AddListener(clientConnection.eventChannel)

	go clientConnection.receiveRequests()
	go clientConnection.process()

	return clientConnection
}
