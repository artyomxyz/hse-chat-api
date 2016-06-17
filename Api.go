package main

import "gopkg.in/mgo.v2"

// API is servers instance
type API struct {
	db             *mgo.Database
	messageManager *MessageManager
	userManager    *UserManager
	eventManager   *EventManager
	listener       *ProtoListener
}

// Initialize initializes api
func (api *API) Initialize() error {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		return err
	}

	db := session.DB("chat")
	err = db.C("users").EnsureIndex(mgo.Index{
		Key:      []string{"username"},
		Unique:   true,
		DropDups: true,
	})
	if err != nil {
		return err
	}

	api.eventManager = NewEventManager()
	api.messageManager = NewMessageManager(api.eventManager.InputChannel, db)
	api.userManager = NewUserManager(api.eventManager.InputChannel, db)

	api.listener, err = ProtoListen("tcp", ":8080")
	if err != nil {
		return err
	}

	return nil
}

// Loop accept loop
func (api *API) Loop() {
	for {
		connection, err := api.listener.Accept()
		if err != nil {
			connection.Close()
		}
		NewClientConnection(connection, api)
	}
}
