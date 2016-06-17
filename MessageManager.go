package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MessageManager manages messages and listeners to new messages
type MessageManager struct {
	eventChannel chan Event
	db           *mgo.Database
}

// AddMessage add new message and emit events
func (messageManager *MessageManager) AddMessage(message Message) error {
	err := messageManager.db.C("messages").Insert(bson.M{
		"author":   message.Author,
		"text":     message.Text,
		"receiver": message.Receiver,
		"date":     message.Date,
	})

	if err != nil {
		return err
	}

	messageManager.eventChannel <- NewMessageEvent{message}

	return nil
}

// GetMessagesBetweenTwoUsers returns messages between user1 and user2
func (messageManager *MessageManager) GetMessagesBetweenTwoUsers(user1 string, user2 string) ([]Message, error) {
	var messages []Message

	err := messageManager.db.C("messages").Find(bson.M{
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
func NewMessageManager(eventChannel chan Event, db *mgo.Database) *MessageManager {
	return &MessageManager{eventChannel, db}
}
