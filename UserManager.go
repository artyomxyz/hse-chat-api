package main

import (
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserManager manages messages and listeners to new messages
type UserManager struct {
	usersSessionsCount map[string]int
	mx                 sync.RWMutex
}

// IncUserSessionsCount inc number of current session of user on server
func (usrMngr *UserManager) IncUserSessionsCount(username string) {
	usrMngr.mx.Lock()

	sessionsCount, ok := usrMngr.usersSessionsCount[username]

	if !ok {
		sessionsCount = 0
	}

	usrMngr.usersSessionsCount[username] = sessionsCount + 1

	if sessionsCount == 0 {
		go evtMngr.Emit(UpdateUserEvent{
			user: User{
				Username: username,
				Online:   true,
			},
		})
	}

	usrMngr.mx.Unlock()
}

// DecUserSessionsCount dec number of current session of user on server
func (usrMngr *UserManager) DecUserSessionsCount(username string) {
	usrMngr.mx.Lock()

	sessionsCount, ok := usrMngr.usersSessionsCount[username]

	if !ok || sessionsCount < 1 {
		sessionsCount = 1
	}

	if sessionsCount == 1 {
		go evtMngr.Emit(UpdateUserEvent{
			user: User{
				Username: username,
				Online:   true,
			},
		})
	}

	usrMngr.usersSessionsCount[username] = sessionsCount - 1
	delete(usrMngr.usersSessionsCount, username)

	usrMngr.mx.Unlock()
}

// IsUserOnline return is user online
func (usrMngr *UserManager) IsUserOnline(username string) bool {
	usrMngr.mx.RLock()
	sessionsCount, ok := usrMngr.usersSessionsCount[username]
	usrMngr.mx.RUnlock()

	return ok && sessionsCount != 0
}

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

	for i := range users {
		users[i].Online = usrMngr.IsUserOnline(users[i].Username)
	}

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
	return &UserManager{usersSessionsCount: make(map[string]int)}
}
