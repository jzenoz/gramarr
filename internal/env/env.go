package env

import (
	"fmt"
	"strings"

	"github.com/jzenoz/gramarr/internal/config"
	"github.com/jzenoz/gramarr/internal/conversation"
	"github.com/jzenoz/gramarr/internal/radarr"
	"github.com/jzenoz/gramarr/internal/sonarr"
	"github.com/jzenoz/gramarr/internal/users"
	"github.com/jzenoz/gramarr/internal/util"
	"gopkg.in/tucnak/telebot.v2"
)

// Env struct for gramarr environment
type Env struct {
	Config *config.Config
	Users  *users.UserDB
	Bot    *telebot.Bot
	CM     *conversation.ConversationManager
	Radarr *radarr.Client
	Sonarr *sonarr.Client
}

// RequirePrivate ensures private message usage
func (e *Env) RequirePrivate(h func(m *telebot.Message)) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		if !m.Private() {
			return
		}
		h(m)
	}
}

// RequireAuth ensures user auth before message response
func (e *Env) RequireAuth(access users.UserAccess, h func(m *telebot.Message)) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		user, _ := e.Users.User(m.Sender.ID)
		var msg []string

		// Is Revoked?
		if user.IsRevoked() {
			// Notify User
			msg = append(msg, "Your access has been revoked and you cannot reauthorize.")
			msg = append(msg, "Please reach out to the bot owner for support.")
			util.SendError(e.Bot, m.Sender, strings.Join(msg, "\n"))

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Revoked users %s attempted the following command:", util.DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			util.SendAdmin(e.Bot, e.Users.Admins(), strings.Join(msg, "\n"))
			return
		}

		// Is Not Member?
		isAuthorized := user.IsAdmin() || user.IsMember()
		if !isAuthorized && access != users.UANone {
			// Notify User
			util.SendError(e.Bot, m.Sender, "You are not authorized to use this bot.\n`/auth [password]` to authorize.")

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Unauthorized users %s attempted the following command:", util.DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			util.SendAdmin(e.Bot, e.Users.Admins(), strings.Join(msg, "\n"))
			return
		}

		// Is Non-Admin and requires Admin?
		if !user.IsAdmin() && access == users.UAAdmin {
			// Notify User
			util.SendError(e.Bot, m.Sender, "Only admins can use this command.")

			// Notify Admins
			msg = append(msg, fmt.Sprintf("User %s attempted the following admin command:", util.DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			util.SendAdmin(e.Bot, e.Users.Admins(), strings.Join(msg, "\n"))
			return
		}

		h(m)
	}
}
