package main

import "github.com/hse-chat/hse-chat-api/HseMsg"

// Message struct represents message in system
type Message struct {
	Author   string
	Receiver string
	Text     string
	Date     int64
}

// ToProtoMessage convers struct to *HseMsg.Message
func (msg Message) ToProtoMessage() *HseMsg.Message {
	return &HseMsg.Message{
		Author:   &msg.Author,
		Receiver: &msg.Receiver,
		Date:     &msg.Date,
		Text:     &msg.Text,
	}
}
