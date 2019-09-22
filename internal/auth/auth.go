package auth

import (
	"fmt"
	"strings"

	"github.com/tommy647/gramarr/internal/users"
	"github.com/tommy647/gramarr/internal/util"

	"gopkg.in/tucnak/telebot.v2"
)

// Service our auth handling service
type Service struct {
	bot           Bot
	users         User
	adminPassword string
	password      string
}

// New returns a new initialised instance of our Auth service
func New(cfg Config, opts ...Options) *Service {
	s := &Service{
		password:      cfg.Password,
		adminPassword: cfg.AdminPassword,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Service) Auth(m *telebot.Message) {
	var msg []string
	pass := m.Payload
	user, exists := s.users.User(m.Sender.ID)

	// Empty Password?
	if pass == "" {
		if err := s.bot.Send(user, "Usage: `/auth [password]`"); err != nil {
			// log error
		}
		return
	}

	// Is User Already Admin?
	if exists && user.IsAdmin() {
		// Notify User
		msg = append(msg, "You're already authorized.")
		msg = append(msg, "Type /start to begin.")
		if err := s.bot.Send(user, strings.Join(msg, "\n")); err != nil {
			// log error
		}
		return
	}

	// Check if pass is Admin Password
	if pass == s.adminPassword {
		if exists {
			user.Access = users.UAAdmin
			_ = s.users.Update(user) // @todo: handle error
		} else {
			newUser := users.User{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Username:  user.Username,
				Access:    users.UAAdmin,
			}
			_ = s.users.Create(newUser) // handle error
		}

		// Notify User
		msg = append(msg, "You have been authorized as an *admin*.")
		msg = append(msg, "Type /start to begin.")
		_ = s.bot.Send(user, strings.Join(msg, "\n")) // handle error

		// Notify Admin
		adminMsg := fmt.Sprintf("%s has been granted admin access.", user.DisplayName())
		_ = s.bot.SendAdmin(s.users.Admins(), adminMsg) // handle error

		return
	}

	// Check if pass is User Password
	if pass == s.password {
		if exists {
			// Notify User
			msg = append(msg, "You're already authorized.")
			msg = append(msg, "Type /start to begin.")
			_ = s.bot.Send(user, strings.Join(msg, "\n")) // @todo: handle error
			return
		}
		newUser := users.User{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Access:    users.UAMember,
		}
		_ = s.users.Create(newUser) // todo: handle error

		// Notify User
		msg = append(msg, "You have been authorized.")
		msg = append(msg, "Type /start to begin.")
		_ = s.bot.Send(user, strings.Join(msg, "\n")) // @todo: handle error

		// Notify Admin
		adminMsg := fmt.Sprintf("%s has been granted acccess.", user.DisplayName())
		_ = s.bot.SendAdmin(s.users.Admins(), adminMsg) // @todo: handle errors
		return
	}
	// Notify User
	_ = s.bot.SendError(user, "Your password is invalid.") // @todo: handle error

	// Notify Admin
	adminMsg := "%s made an invalid auth request with password: %s"
	adminMsg = fmt.Sprintf(adminMsg, user.DisplayName(), util.EscapeMarkdown(m.Payload))
	_ = s.bot.SendAdmin(s.users.Admins(), adminMsg) // @todo: handle error
}
