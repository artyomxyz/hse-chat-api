package main

import "github.com/hse-chat/hse-chat-api/HseMsg"

// User struct represents user in system
type User struct {
	Username string
	Password string
	Online   bool
}

// ToProtoUser convers struct to *HseMsg.User
func (usr User) ToProtoUser() *HseMsg.User {
	return &HseMsg.User{
		Username: &usr.Username,
		Online:   &usr.Online,
	}
}
