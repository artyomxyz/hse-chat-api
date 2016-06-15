package main

import (
	"time"

	"github.com/hse-chat/hse-chat-api/HseMsg"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type clientError struct{}

func (clErr *clientError) Error() string {
	return "Client error"
}

// Client object, represents conected client
type Client struct {
	username string
	signedIn bool
}

// NewClient Creates new unauthorized client
func NewClient() Client {
	return Client{"", false}
}

// SignUp perform sign up action
func (cl *Client) SignUp(username string, password string) (*SignUpResult, error) {
	if cl.signedIn {
		return nil, &clientError{}
	}

	if len(username) < 4 || len(password) < 4 {
		return &SignUpResult{HseMsg.Result_SignUpResult_VALIDATION_ERROR}, nil
	}

	err := db.C("users").Insert(bson.M{
		"username": username,
		"password": password,
	})

	if merr, ok := err.(*mgo.LastError); ok && merr.Code == 11000 {
		return &SignUpResult{HseMsg.Result_SignUpResult_USERNAME_IS_TAKEN}, nil
	}

	if err != nil {
		return nil, err
	}

	return &SignUpResult{HseMsg.Result_SignUpResult_SIGNED_UP}, nil
}

// SignIn perform sign in action
func (cl *Client) SignIn(username string, password string) (*SignInResult, error) {
	if cl.signedIn {
		return nil, &clientError{}
	}

	var user struct {
		Username string
		Password string
	}
	err := db.C("users").Find(bson.M{
		"username": username,
		"password": password,
	}).One(&user)

	if err == mgo.ErrNotFound {
		return &SignInResult{HseMsg.Result_SignInResult_USER_NOT_FOUND}, nil
	}

	if err != nil {
		return nil, err
	}

	cl.username = username
	cl.signedIn = true
	return &SignInResult{HseMsg.Result_SignInResult_SIGNED_IN}, nil
}

// GetUsers get users
func (cl *Client) GetUsers() (*GetUsersResult, error) {
	if !cl.signedIn {
		return nil, &clientError{}
	}

	var users []struct {
		Username string
	}

	err := db.C("users").Find(nil).Select(bson.M{
		"username": 1,
	}).All(&users)

	if err != nil {
		return nil, err
	}

	return &GetUsersResult{users}, nil
}

// GetMessagesWithUser gets messages with users
func (cl *Client) GetMessagesWithUser(peer string) (*GetMessagesWithUserResult, error) {
	if !cl.signedIn {
		return nil, &clientError{}
	}

	var messages []struct {
		Author string
		Text   string
		Date   int64
	}

	err := db.C("messages").Find(nil).Select(bson.M{
		"$or": []interface{}{
			bson.M{
				"author":   cl.username,
				"receiver": peer,
			},
			bson.M{
				"author":   peer,
				"receiver": cl.username,
			},
		},
	}).Select(bson.M{
		"author": 1,
		"text":   1,
		"date":   1,
	}).All(&messages)

	if err != nil {
		return nil, err
	}

	return &GetMessagesWithUserResult{messages}, nil
}

// SendMessageToUser send message to user
func (cl *Client) SendMessageToUser(receiver string, text string) (*SendMessageToUserResult, error) {
	if !cl.signedIn {
		return nil, &clientError{}
	}

	if len(text) == 0 {
		return &SendMessageToUserResult{HseMsg.Result_SendMessageToUserResult_EMPTY_MESSAGE}, nil
	}

	var user struct {
		Username string
	}
	err := db.C("users").Find(bson.M{
		"username": receiver,
	}).One(&user)

	if err == mgo.ErrNotFound {
		return &SendMessageToUserResult{HseMsg.Result_SendMessageToUserResult_USER_NOT_FOUND}, nil
	}

	if err != nil {
		return nil, err
	}

	err = db.C("messages").Insert(bson.M{
		"author":   cl.username,
		"text":     text,
		"receiver": receiver,
		"date":     time.Now().Unix(),
	})

	if err != nil {
		return nil, err
	}

	return &SendMessageToUserResult{HseMsg.Result_SendMessageToUserResult_SENT}, nil
}
