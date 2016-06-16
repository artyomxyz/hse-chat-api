package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserManager manages messages and listeners to new messages
type UserManager struct{}

// AddUserUsernameIsTakenError occurs when username is taken
type AddUserUsernameIsTakenError struct{}

func (err *AddUserUsernameIsTakenError) Error() string {
	return "Username is taken"
}

// AddUser add new message and emit events
func (usrMngr *UserManager) AddUser(usr User) error {
	err := db.C("users").Insert(bson.M{
		"username": usr.Username,
		"password": usr.Password,
	})

	if merr, ok := err.(*mgo.LastError); ok && merr.Code == 11000 {
		return &AddUserUsernameIsTakenError{}
	}

	if err != nil {
		return err
	}

	// TODO: do it using goroutine and channels in event manager
	go evtMngr.Emit(NewUserEvent{usr})

	return nil
}

// FindByUsernameAndPassword finds user by username and password
func (usrMngr *UserManager) FindByUsernameAndPassword(username string, password string) (*User, error) {
	var user User
	err := db.C("users").Find(bson.M{
		"username": username,
		"password": password,
	}).One(&user)

	if err == mgo.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUsers return slice of users
func (usrMngr *UserManager) GetUsers() ([]User, error) {
	var users []User

	err := db.C("users").Find(nil).Select(bson.M{
		"username": 1,
	}).All(&users)

	return users, err
}

// Exists check if user with certain username exists
func (usrMngr *UserManager) Exists(username string) (bool, error) {
	var user User
	err := db.C("users").Find(bson.M{
		"username": username,
	}).One(&user)

	if err == mgo.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// NewUserManager creates new message managers
func NewUserManager() *UserManager {
	return &UserManager{}
}
