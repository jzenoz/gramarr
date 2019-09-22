package auth

import (
	"github.com/tommy647/gramarr/internal/users"
)

type Options func(s *Service)

type Bot interface {
	Send(to users.User, msg interface{}) error
	SendError(to users.User, msg interface{}) error
	SendToAdmins(msg interface{}) error
}

type User interface {
	User(int) (users.User, bool)
	Update(users.User) error
	Create(users.User) error
	Admins() []users.User
}

// WithBot adds the telegram bot to the service
func WithBot(b Bot) Options {
	return func(s *Service) {
		s.bot = b
	}
}

func WithUsers(users User) Options {
	return func(s *Service) {
		s.users = users
	}
}
