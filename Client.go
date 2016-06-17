package main

import (
	"fmt"
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
	user           *User
	signedIn       bool
	userManager    *UserManager
	messageManager *MessageManager
}

// NewClient Creates new unauthorized client
func NewClient(userManager *UserManager, messageManager *MessageManager) Client {
	return Client{nil, false, userManager, messageManager}
}

// Finish Finishes user connection
func (client *Client) Finish() {
	if client.user != nil {
		client.userManager.DecUserSessionsCount(client.user.Username)
	}
}

// SignUp perform sign up action
func (client *Client) SignUp(username string, password string) (*SignUpResult, error) {
	if client.signedIn {
		return nil, &clientError{"Signed in already"}
	}

	if len(username) < 4 || len(password) < 4 {
		return &SignUpResult{HseMsg.Result_SignUpResult_VALIDATION_ERROR}, nil
	}

	user := User{
		Username: username,
		Password: password,
	}

	err := client.userManager.AddUser(user)

	if _, ok := err.(*AddUserUsernameIsTakenError); ok {
		return &SignUpResult{HseMsg.Result_SignUpResult_USERNAME_IS_TAKEN}, nil
	}

	if err != nil {
		return nil, err
	}

	client.user = &User{Username: username}
	client.signedIn = true

	// TODO: investigae this go
	go client.userManager.IncUserSessionsCount(username)

	return &SignUpResult{HseMsg.Result_SignUpResult_SIGNED_UP}, nil
}

// SignIn perform sign in action
func (client *Client) SignIn(username string, password string) (*SignInResult, error) {
	if client.signedIn {
		return nil, &clientError{"Signed in already"}
	}

	user, err := client.userManager.FindByUsernameAndPassword(username, password)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return &SignInResult{HseMsg.Result_SignInResult_USER_NOT_FOUND}, nil
	}

	client.user = user
	client.signedIn = true

	go client.userManager.IncUserSessionsCount(username)

	return &SignInResult{HseMsg.Result_SignInResult_SIGNED_IN}, nil
}

// GetUsers get users
func (client *Client) GetUsers() (*GetUsersResult, error) {
	if !client.signedIn {
		return nil, &clientError{"Not signed in"}
	}

	users, err := client.userManager.GetUsers()

	if err != nil {
		return nil, err
	}

	return &GetUsersResult{users}, nil
}

// GetMessagesWithUser gets messages with users
func (client *Client) GetMessagesWithUser(peer string) (*GetMessagesWithUserResult, error) {
	if !client.signedIn {
		return nil, &clientError{"Not signed in"}
	}

	messages, err := client.messageManager.GetMessagesBetweenTwoUsers(client.user.Username, peer)
	if err != nil {
		return nil, err
	}

	return &GetMessagesWithUserResult{messages}, nil
}

// SendMessageToUser send message to user
func (client *Client) SendMessageToUser(receiver string, text string) (*SendMessageToUserResult, error) {
	if !client.signedIn {
		return nil, &clientError{"Not signed in"}
	}

	if len(text) == 0 {
		return &SendMessageToUserResult{HseMsg.Result_SendMessageToUserResult_EMPTY_MESSAGE}, nil
	}

	exists, err := client.userManager.Exists(receiver)

	if err != nil {
		return nil, err
	}

	if !exists {
		return &SendMessageToUserResult{HseMsg.Result_SendMessageToUserResult_USER_NOT_FOUND}, nil
	}

	date := time.Now().Unix()

	err = client.messageManager.AddMessage(Message{
		Author:   client.user.Username,
		Text:     text,
		Receiver: receiver,
		Date:     date,
	})

	if err != nil {
		return nil, err
	}

	return &SendMessageToUserResult{HseMsg.Result_SendMessageToUserResult_SENT}, nil
}
