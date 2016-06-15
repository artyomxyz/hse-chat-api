package main

import "github.com/hse-chat/hse-chat-api/HseMsg"

// Result represents result of request sent by user
type Result interface {
	ToServerMessage(id uint32) *HseMsg.ServerMessage
}

// SignUpResult result of sign up
type SignUpResult struct {
	status HseMsg.Result_SignUpResult_SignUpResultStatus
}

// ToServerMessage converts to server message
func (res SignUpResult) ToServerMessage(id uint32) *HseMsg.ServerMessage {
	return &HseMsg.ServerMessage{
		Message: &HseMsg.ServerMessage_Result{
			Result: &HseMsg.Result{
				Id: &id,
				Value: &HseMsg.Result_SignUp{
					SignUp: &HseMsg.Result_SignUpResult{
						Status: &res.status,
					},
				},
			},
		},
	}
}

// SignInResult result of sign in
type SignInResult struct {
	status HseMsg.Result_SignInResult_SignInResultStatus
}

// ToServerMessage converts to server message
func (res SignInResult) ToServerMessage(id uint32) *HseMsg.ServerMessage {
	return &HseMsg.ServerMessage{
		Message: &HseMsg.ServerMessage_Result{
			Result: &HseMsg.Result{
				Id: &id,
				Value: &HseMsg.Result_SignIn{
					SignIn: &HseMsg.Result_SignInResult{
						Status: &res.status,
					},
				},
			},
		},
	}
}

// GetUsersResult result getting users
type GetUsersResult struct {
	users []struct {
		Username string
	}
}

// ToServerMessage converts to server message
func (res GetUsersResult) ToServerMessage(id uint32) *HseMsg.ServerMessage {
	users := make([]*HseMsg.User, len(res.users))

	for i := range res.users {
		users[i] = &HseMsg.User{
			Username: &res.users[i].Username,
		}
	}

	return &HseMsg.ServerMessage{
		Message: &HseMsg.ServerMessage_Result{
			Result: &HseMsg.Result{
				Id: &id,
				Value: &HseMsg.Result_GetUsers{
					GetUsers: &HseMsg.Result_GetUsersResult{
						Users: users,
					},
				},
			},
		},
	}
}

// GetMessagesWithUserResult result of getting messages with user
type GetMessagesWithUserResult struct {
	messages []struct {
		Author string
		Text   string
		Date   int64
	}
}

// ToServerMessage converts to server message
func (res GetMessagesWithUserResult) ToServerMessage(id uint32) *HseMsg.ServerMessage {
	messages := make([]*HseMsg.Message, len(res.messages))

	for i := range res.messages {
		messages[i] = &HseMsg.Message{
			Author: &res.messages[i].Author,
			Text:   &res.messages[i].Text,
			Date:   &res.messages[i].Date,
		}
	}

	return &HseMsg.ServerMessage{
		Message: &HseMsg.ServerMessage_Result{
			Result: &HseMsg.Result{
				Id: &id,
				Value: &HseMsg.Result_GetMessagesWithUser{
					GetMessagesWithUser: &HseMsg.Result_GetMessagesWithUserResult{
						Messages: messages,
					},
				},
			},
		},
	}
}

// SendMessageToUserResult result of sending messages
type SendMessageToUserResult struct {
	status HseMsg.Result_SendMessageToUserResult_SendMessageToUserResultStatus
}

// ToServerMessage converts to server message
func (res SendMessageToUserResult) ToServerMessage(id uint32) *HseMsg.ServerMessage {
	return &HseMsg.ServerMessage{
		Message: &HseMsg.ServerMessage_Result{
			Result: &HseMsg.Result{
				Id: &id,
				Value: &HseMsg.Result_SendMessageToUser{
					SendMessageToUser: &HseMsg.Result_SendMessageToUserResult{
						Status: &res.status,
					},
				},
			},
		},
	}
}
