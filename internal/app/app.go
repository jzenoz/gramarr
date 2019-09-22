package app

import (
	"fmt"
	"strings"

	"github.com/tommy647/gramarr/internal/conversation"
	"github.com/tommy647/gramarr/internal/radarr"
	"github.com/tommy647/gramarr/internal/sonarr"
	"github.com/tommy647/gramarr/internal/users"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Authoriser interface to our auth service
type Authoriser interface {
	Auth(message *tb.Message)
}

type Bot interface {
	Start()
	Send(users.User, interface{}) error
	SendKeyboardList(users.User, string, []string) error
	SendToAdmins(interface{}) error
	Name() string
	Handle(interface{}, interface{})
}

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

func (s *Service) HandleAuth(m *tb.Message) {
	s.Auth.Auth(m)
}

func (s *Service) RequirePrivate(h func(m *tb.Message)) func(m *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}
		h(m)
	}
}

func (s *Service) RequireAuth(access users.UserAccess, h func(m *tb.Message)) func(m *tb.Message) {
	return func(m *tb.Message) {
		user, _ := s.Users.User(m.Sender.ID)
		var msg []string

		// Is Revoked?
		if user.IsRevoked() {
			// Notify User
			msg = append(msg, "Your access has been revoked and you cannot reauthorize.")
			msg = append(msg, "Please reach out to the bot owner for support.")
			_ = s.Bot.Send(user, strings.Join(msg, "\n")) // @todo: handle error

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Revoked users %s attempted the following command:", user.DisplayName()))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
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
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			_ = s.Bot.SendToAdmins(strings.Join(msg, "\n")) // @todo: handle error
			return
		}

		// Is Non-Admin and requires Admin?
		if !user.IsAdmin() && access == users.UAAdmin {
			// Notify User
			_ = s.Bot.Send(user, "Only admins can use this command.") // @todo: handle error

			// Notify Admins
			msg = append(msg, fmt.Sprintf("User %s attempted the following admin command:", user.DisplayName()))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			_ = s.Bot.SendToAdmins(strings.Join(msg, "\n")) // @todo: handle error
			return
		}

		h(m)
	}
}
