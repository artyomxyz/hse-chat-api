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
	eventManager       chan Event
	db                 *mgo.Database
}

// IncUserSessionsCount inc number of current session of user on server
func (userManager *UserManager) IncUserSessionsCount(username string) {
	userManager.mx.Lock()

	sessionsCount, ok := userManager.usersSessionsCount[username]

	if !ok {
		sessionsCount = 0
	}

	userManager.usersSessionsCount[username] = sessionsCount + 1

	if sessionsCount == 0 {
		userManager.eventManager <- UpdateUserEvent{
			user: User{
				Username: username,
				Online:   true,
			},
		}
	}

	userManager.mx.Unlock()
}

// DecUserSessionsCount dec number of current session of user on server
func (userManager *UserManager) DecUserSessionsCount(username string) {
	userManager.mx.Lock()

	sessionsCount, ok := userManager.usersSessionsCount[username]

	if !ok || sessionsCount < 1 {
		sessionsCount = 1
	}

	if sessionsCount == 1 {
		userManager.eventManager <- UpdateUserEvent{
			user: User{
				Username: username,
				Online:   true,
			},
		}
	}

	userManager.usersSessionsCount[username] = sessionsCount - 1
	delete(userManager.usersSessionsCount, username)

	userManager.mx.Unlock()
}

// IsUserOnline return is user online
func (userManager *UserManager) IsUserOnline(username string) bool {
	userManager.mx.RLock()
	sessionsCount, ok := userManager.usersSessionsCount[username]
	userManager.mx.RUnlock()

	return ok && sessionsCount != 0
}

// AddUserUsernameIsTakenError occurs when username is taken
type AddUserUsernameIsTakenError struct{}

func (err *AddUserUsernameIsTakenError) Error() string {
	return "Username is taken"
}

// AddUser add new message and emit events
func (userManager *UserManager) AddUser(user User) error {
	err := userManager.db.C("users").Insert(bson.M{
		"username": user.Username,
		"password": user.Password,
	})

	if merr, ok := err.(*mgo.LastError); ok && merr.Code == 11000 {
		return &AddUserUsernameIsTakenError{}
	}

	if err != nil {
		return err
	}

	userManager.eventManager <- NewUserEvent{user}

	return nil
}

// FindByUsernameAndPassword finds user by username and password
func (userManager *UserManager) FindByUsernameAndPassword(username string, password string) (*User, error) {
	var user User
	err := userManager.db.C("users").Find(bson.M{
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
func (userManager *UserManager) GetUsers() ([]User, error) {
	var users []User

	err := userManager.db.C("users").Find(nil).Select(bson.M{
		"username": 1,
	}).All(&users)

	for i := range users {
		users[i].Online = userManager.IsUserOnline(users[i].Username)
	}

	return users, err
}

// Exists check if user with certain username exists
func (userManager *UserManager) Exists(username string) (bool, error) {
	var user User
	err := userManager.db.C("users").Find(bson.M{
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
func NewUserManager(eventManager chan Event, db *mgo.Database) *UserManager {
	return &UserManager{
		usersSessionsCount: make(map[string]int),
		eventManager:       eventManager,
		db:                 db,
	}
}
