package message

import "github.com/tommy647/gramarr/internal/users"

type Message struct {
	Text    string
	Payload string
	User    users.User
	Private bool
	Chat    Chat
}

type Chat struct {
	ID interface{}
}
