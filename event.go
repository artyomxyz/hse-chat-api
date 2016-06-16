package main

import "github.com/hse-chat/hse-chat-api/HseMsg"

// Event represents event of any kind
type Event interface {
	ToProtoEvent() *HseMsg.Event
}

// EventToServerMessage converts event to serverMessage
func EventToServerMessage(evt Event) *HseMsg.ServerMessage {
	return &HseMsg.ServerMessage{
		Message: &HseMsg.ServerMessage_Event{
			Event: evt.ToProtoEvent(),
		},
	}
}

// NewMessageEvent represents new message event
type NewMessageEvent struct {
	message Message
}

// ToProtoEvent converts message to protobuf event
func (evt NewMessageEvent) ToProtoEvent() *HseMsg.Event {
	return &HseMsg.Event{
		Event: &HseMsg.Event_NewMessage_{
			NewMessage: &HseMsg.Event_NewMessage{
				Message: evt.message.ToProtoMessage(),
			},
		},
	}
}
