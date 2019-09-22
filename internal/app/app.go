package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/tommy647/gramarr/internal/conversation"
	"github.com/tommy647/gramarr/internal/message"
	"github.com/tommy647/gramarr/internal/radarr"
	"github.com/tommy647/gramarr/internal/sonarr"
	"github.com/tommy647/gramarr/internal/users"
)

// Authoriser interface to our auth service
type Authoriser interface {
	Auth(message *message.Message)
}

type Bot interface {
	Start()
	Send(users.User, interface{}) error
	SendKeyboardList(users.User, string, []string) error
	SendToAdmins(interface{}) error
	Name() string
	Handle(interface{}, interface{})
	GetUserID(interface{}) interface{}
	IsPrivate(interface{}) bool
	GetText(interface{}) string
	WithMessage(interface{}) (*message.Message, error)
}

type ContextUserKey string

const cxtUserKey ContextUserKey = `user`

// Service our main service
// @todo: why are these exposed?
type Service struct {
	Auth   Authoriser
	Users  *users.UserDB
	Bot    Bot
	CM     *conversation.ConversationManager
	Radarr *radarr.Client
	Sonarr *sonarr.Client
}

func (s *Service) HandleAuth(m *message.Message) { s.Auth.Auth(m) }

// WithUser infers our user from the message and add into the context
func (s *Service) WithUser(h func(m *message.Message)) func(m *message.Message) {
	log.Print("WithUser")
	return func(m *message.Message) {
		user, exists := s.Users.User(m.User.ID)
		if !exists {
			if err := s.Users.Create(m.User); err != nil {
				log.Println(err.Error())
			}
		}
		m.User = user
		h(m)
	}
}

func (s *Service) RequirePrivate(h func(m *message.Message)) func(m *message.Message) {
	log.Print("RequirePrivate")
	return func(m *message.Message) {
		if !m.Private {
			log.Println("not a private message")
			return
		}
		h(m)
	}
}

func (s *Service) RequireAuth(access users.UserAccess, h func(m *message.Message)) func(m *message.Message) {
	return func(m *message.Message) {
		user := m.User
		var msg []string

		// Is Revoked?
		if user.IsRevoked() {
			// Notify User
			msg = append(msg, "Your access has been revoked and you cannot reauthorize.")
			msg = append(msg, "Please reach out to the bot owner for support.")
			_ = s.Bot.Send(user, strings.Join(msg, "\n")) // @todo: handle error

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Revoked users %s attempted the following command:", user.DisplayName()))
			msg = append(msg, fmt.Sprintf("`%s`", s.Bot.GetText(m)))
			_ = s.Bot.SendToAdmins(strings.Join(msg, "\n")) // @todo: handle error
			return
		}

		// Is Not Member?
		isAuthorized := user.IsAdmin() || user.IsMember()
		if !isAuthorized && access != users.UANone {
			// Notify User
			_ = s.Bot.Send(user, "You are not authorized to use this bot.\n`/auth [password]` to authorize.") // @todo: handle error

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Unauthorized users %s attempted the following command:", user.DisplayName()))
			msg = append(msg, fmt.Sprintf("`%s`", s.Bot.GetText(m)))
			_ = s.Bot.SendToAdmins(strings.Join(msg, "\n")) // @todo: handle error
			return
		}

		// Is Non-Admin and requires Admin?
		if !user.IsAdmin() && access == users.UAAdmin {
			// Notify User
			_ = s.Bot.Send(user, "Only admins can use this command.") // @todo: handle error

			// Notify Admins
			msg = append(msg, fmt.Sprintf("User %s attempted the following admin command:", user.DisplayName()))
			msg = append(msg, fmt.Sprintf("`%s`", s.Bot.GetText(m)))
			_ = s.Bot.SendToAdmins(strings.Join(msg, "\n")) // @todo: handle error
			return
		}

		h(m)
	}
}
