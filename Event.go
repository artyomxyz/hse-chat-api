package main

import "github.com/hse-chat/hse-chat-api/HseMsg"

// Event represents event of any kind
type Event interface {
	ToProtoEvent() *HseMsg.Event
	IsAccessibleBy(usr *User) bool
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
		Value: &HseMsg.Event_NewMessage_{
			NewMessage: &HseMsg.Event_NewMessage{
				Message: evt.message.ToProtoMessage(),
			},
		},
	}
}

// IsAccessibleBy checks if this event accesible by certain user
func (evt NewMessageEvent) IsAccessibleBy(usr *User) bool {
	return evt.message.Author == usr.Username || evt.message.Receiver == usr.Username
}

// NewUserEvent represents new message event
type NewUserEvent struct {
	user User
}

// ToProtoEvent converts message to protobuf event
func (evt NewUserEvent) ToProtoEvent() *HseMsg.Event {
	return &HseMsg.Event{
		Value: &HseMsg.Event_NewUser_{
			NewUser: &HseMsg.Event_NewUser{
				User: evt.user.ToProtoUser(),
			},
		},
	}
}

// IsAccessibleBy checks if this event accesible by certain user
func (evt NewUserEvent) IsAccessibleBy(usr *User) bool {
	return true
}

// UpdateUserEvent represents new message event
type UpdateUserEvent struct {
	user User
}

// ToProtoEvent converts message to protobuf event
func (evt UpdateUserEvent) ToProtoEvent() *HseMsg.Event {
	return &HseMsg.Event{
		Value: &HseMsg.Event_UpdateUser_{
			UpdateUser: &HseMsg.Event_UpdateUser{
				User: evt.user.ToProtoUser(),
			},
		},
	}
}

// IsAccessibleBy checks if this event accesible by certain user
func (evt UpdateUserEvent) IsAccessibleBy(usr *User) bool {
	return true
}
