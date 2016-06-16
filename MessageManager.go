package main

import "gopkg.in/mgo.v2/bson"

// MessageManager manages messages and listeners to new messages
type MessageManager struct{}

// AddMessage add new message and emit events
func (msgMngr *MessageManager) AddMessage(msg Message) error {
	err := db.C("messages").Insert(bson.M{
		"author":   msg.Author,
		"text":     msg.Text,
		"receiver": msg.Receiver,
		"date":     msg.Date,
	})

	if err != nil {
		return err
	}

	go evtMngr.Emit(NewMessageEvent{msg})

	return nil
}

// GetMessagesBetweenTwoUsers returns messages between user1 and user2
func (msgMngr *MessageManager) GetMessagesBetweenTwoUsers(user1 string, user2 string) ([]Message, error) {
	var messages []Message

	err := db.C("messages").Find(bson.M{
		"$or": []interface{}{
			bson.M{
				"author":   user1,
				"receiver": user2,
			},
			bson.M{
				"author":   user2,
				"receiver": user1,
			},
		},
	}).All(&messages)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

// NewMessageManager creates new message managers
func NewMessageManager() *MessageManager {
	return &MessageManager{}
}
