package main

import (
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/hse-chat/hse-chat-api/HseMsg"
)

// Message struct represents message in system
type Message struct {
	Author   string
	Receiver string
	Text     string
	Date     int64
}

// ToServerMessage converts message to ServerMessage event NewMessage
func (msg Message) ToServerMessage() *HseMsg.ServerMessage {
	return &HseMsg.ServerMessage{
		Message: &HseMsg.ServerMessage_Event{
			Event: &HseMsg.Event{
				Event: &HseMsg.Event_NewMessage_{
					NewMessage: &HseMsg.Event_NewMessage{
						Message: &HseMsg.Message{
							Author:   &msg.Author,
							Receiver: &msg.Receiver,
							Date:     &msg.Date,
							Text:     &msg.Text,
						},
					},
				},
			},
		},
	}
}

// MessageManager manages messages and listeners to new messages
type MessageManager struct {
	mx        sync.RWMutex
	listeners []chan Message
}

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

	go msgMngr.emit(msg)

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

func (msgMngr *MessageManager) emit(msg Message) {
	msgMngr.mx.RLock()
	for _, ln := range msgMngr.listeners {
		ln <- msg
	}
	msgMngr.mx.RUnlock()
}

// AddListener add channel as a listener to message event
func (msgMngr *MessageManager) AddListener(ln chan Message) {
	msgMngr.mx.Lock()
	msgMngr.listeners = append(msgMngr.listeners, ln)
	msgMngr.mx.Unlock()
}

// RemoveListener removes channel from listeners
func (msgMngr *MessageManager) RemoveListener(rln chan Message) {
	msgMngr.mx.Lock()
	for i, ln := range msgMngr.listeners {
		if ln == rln {
			msgMngr.listeners = append(msgMngr.listeners[:i], msgMngr.listeners[i+1:]...)
			break
		}
	}
	msgMngr.mx.Unlock()
}

// NewMessageManager creates new message managers
func NewMessageManager() *MessageManager {
	return &MessageManager{}
}
