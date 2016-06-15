package main

import "sync"

type Message struct {
	author   string
	receiver string
	text     string
	date     int64
}

type MessageEmitter struct {
	mx        sync.Mutex // TODO: change to rwmutex
	listeners []chan Message
}

func (me *MessageEmitter) Emit(msg Message) {
	me.mx.Lock()
	for _, ln := range me.listeners {
		ln <- msg
	}
	me.mx.Unlock()
}

func (me *MessageEmitter) AddListener(ln chan Message) {
	me.mx.Lock()
	me.listeners = append(me.listeners, ln)
	me.mx.Unlock()
}

func (me *MessageEmitter) RemoveListener(chan Message) {

}

func NewMessagEmitter() *MessageEmitter {
	return &MessageEmitter{}
}
