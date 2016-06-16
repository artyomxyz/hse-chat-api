package main

import "sync"

// User struct represents user in system
type User struct {
	Username string
}

// UserManager manages users and send events to listeners
type UserManager struct {
	mx        sync.RWMutex
	listeners []chan User
}

func (usrManager *UserManager) emit(usr User) {
	usrManager.mx.RLock()
	for _, ln := range usrManager.listeners {
		ln <- usr
	}
	usrManager.mx.RUnlock()
}

// AddListener add channel as a listener to new user event
func (usrManager *UserManager) AddListener(ln chan User) {
	usrManager.mx.Lock()
	usrManager.listeners = append(usrManager.listeners, ln)
	usrManager.mx.Unlock()
}

// RemoveListener removes channel from listeners
func (usrManager *UserManager) RemoveListener(rln chan User) {
	usrManager.mx.Lock()
	for i, ln := range usrManager.listeners {
		if ln == rln {
			usrMng.listeners = append(usrManager.listeners[:i], usrManager.listeners[i+1:]...)
			break
		}
	}
	usrManager.mx.Unlock()
}

// NewUserManager creates new user manager
func NewUserManager() *UserManager {
	return &UserManager{}
}
