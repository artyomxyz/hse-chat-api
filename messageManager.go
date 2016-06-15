package main

import (
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/hse-chat/hse-chat-api/HseMsg"
)

type Message struct {
	author   string
	receiver string
	text     string
	date     int64
}

func (msg Message) ToServerMessage() *HseMsg.ServerMessage {
	return &HseMsg.ServerMessage{
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
}

type MessageManager struct {
	mx        sync.RWMutex // TODO: change to rwmutex
	listeners []chan Message
}

func (me *MessageManager) AddMessage(msg Message) error {
	err := db.C("messages").Insert(bson.M{
		"author":   msg.author,
		"text":     msg.text,
		"receiver": msg.receiver,
		"date":     msg.date,
	})

	if err != nil {
		return err
	}

	go me.Emit(msg)

	return nil
}

func (me *MessageManager) GetMessagesBetweenTwoUsers(user1 string, user2 string) ([]Message, error) {
	var messages []Message

	err := db.C("messages").Find(nil).Select(bson.M{
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
	}).Select(bson.M{
		"author":   1,
		"receiver": 1,
		"text":     1,
		"date":     1,
	}).All(&messages)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (me *MessageManager) Emit(msg Message) {
	me.mx.RLock()
	for _, ln := range me.listeners {
		ln <- msg
	}
	me.mx.RUnlock()
}

func (me *MessageManager) AddListener(ln chan Message) {
	me.mx.Lock()
	me.listeners = append(me.listeners, ln)
	me.mx.Unlock()
}

func (me *MessageManager) RemoveListener(chan Message) {

}

func NewMessageManager() *MessageManager {
	return &MessageManager{}
}
