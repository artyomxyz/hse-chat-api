package main

import "sync"

// EventManager manages events and listeners
type EventManager struct {
	mx        sync.RWMutex
	listeners []chan Event
}

// Emit emit event to all listeners
func (evtMngr *EventManager) Emit(evt Event) {
	evtMngr.mx.RLock()
	for _, ln := range evtMngr.listeners {
		ln <- evt
	}
	evtMngr.mx.RUnlock()
}

// AddListener add channel as a listener to message event
func (evtMngr *EventManager) AddListener(ln chan Event) {
	evtMngr.mx.Lock()
	evtMngr.listeners = append(evtMngr.listeners, ln)
	evtMngr.mx.Unlock()
}

// RemoveListener removes channel from listeners
func (evtMngr *EventManager) RemoveListener(rln chan Event) {
	evtMngr.mx.Lock()
	for i, ln := range evtMngr.listeners {
		if ln == rln {
			evtMngr.listeners = append(evtMngr.listeners[:i], evtMngr.listeners[i+1:]...)
			break
		}
	}
	evtMngr.mx.Unlock()
}

// NewEventManager creates new message managers
func NewEventManager() *EventManager {
	return &EventManager{}
}
