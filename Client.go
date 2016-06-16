package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hse-chat/hse-chat-api/HseMsg"
)

type clientError struct {
	reason string
}

func (clErr *clientError) Error() string {
	return fmt.Sprintf("Client error: %s", clErr.reason)
}

// Client object, represents conected client
type Client struct {
	user     *User
	signedIn bool
}

// NewClient Creates new unauthorized client
func NewClient() Client {
	return Client{nil, false}
}

// SignUp perform sign up action
func (cl *Client) SignUp(username string, password string) (*SignUpResult, error) {
	if cl.signedIn {
		return nil, &clientError{"Signed in already"}
	}

	if len(username) < 4 || len(password) < 4 {
		return &SignUpResult{HseMsg.Result_SignUpResult_VALIDATION_ERROR}, nil
	}

	user := User{
		Username: username,
		Password: password,
	}

	err := usrMngr.AddUser(user)

	if _, ok := err.(*AddUserUsernameIsTakenError); ok {
		return &SignUpResult{HseMsg.Result_SignUpResult_USERNAME_IS_TAKEN}, nil
	}

	if err != nil {
		return nil, err
	}

	cl.user = &User{Username: username}
	cl.signedIn = true

	return &SignUpResult{HseMsg.Result_SignUpResult_SIGNED_UP}, nil
}

// SignIn perform sign in action
func (cl *Client) SignIn(username string, password string) (*SignInResult, error) {
	if cl.signedIn {
		return nil, &clientError{"Signed in already"}
	}

	user, err := usrMngr.FindByUsernameAndPassword(username, password)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return &SignInResult{HseMsg.Result_SignInResult_USER_NOT_FOUND}, nil
	}

	cl.user = user
	cl.signedIn = true
	return &SignInResult{HseMsg.Result_SignInResult_SIGNED_IN}, nil
}

// GetUsers get users
func (cl *Client) GetUsers() (*GetUsersResult, error) {
	if !cl.signedIn {
		return nil, &clientError{"Not signed in"}
	}

	users, err := usrMngr.GetUsers()

	if err != nil {
		return nil, err
	}

	return &GetUsersResult{users}, nil
}

// GetMessagesWithUser gets messages with users
func (cl *Client) GetMessagesWithUser(peer string) (*GetMessagesWithUserResult, error) {
	if !cl.signedIn {
		return nil, &clientError{"Not signed in"}
	}

	messages, err := msgMngr.GetMessagesBetweenTwoUsers(cl.user.Username, peer)
	log.Print(messages)
	if err != nil {
		return nil, err
	}

	return &GetMessagesWithUserResult{messages}, nil
}

// SendMessageToUser send message to user
func (cl *Client) SendMessageToUser(receiver string, text string) (*SendMessageToUserResult, error) {
	if !cl.signedIn {
		return nil, &clientError{"Not signed in"}
	}

	if len(text) == 0 {
		return &SendMessageToUserResult{HseMsg.Result_SendMessageToUserResult_EMPTY_MESSAGE}, nil
	}

	exists, err := usrMngr.Exists(receiver)

	if err != nil {
		return nil, err
	}

	if !exists {
		return &SendMessageToUserResult{HseMsg.Result_SendMessageToUserResult_USER_NOT_FOUND}, nil
	}

	date := time.Now().Unix()

	err = msgMngr.AddMessage(Message{
		Author:   cl.user.Username,
		Text:     text,
		Receiver: receiver,
		Date:     date,
	})

	if err != nil {
		return nil, err
	}

	return &SendMessageToUserResult{HseMsg.Result_SendMessageToUserResult_SENT}, nil
}
