package app

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Service) HandleStart(m *tb.Message) {

	user, exists := s.Users.User(m.Sender.ID)

	var msg []string
	msg = append(msg, fmt.Sprintf("Hello, I'm %s! Use these commands to control me:", s.Bot.Name()))

	if !exists {
		msg = append(msg, "")
		msg = append(msg, "/auth [password] - authenticate with the bot")
	}

	if exists && user.IsAdmin() {
		msg = append(msg, "")
		msg = append(msg, "*Admin*")
		msg = append(msg, "/users - list all bot users")
	}

	if exists && (user.IsMember() || user.IsAdmin()) {
		msg = append(msg, "")
		msg = append(msg, "*Media*")
		msg = append(msg, "/addmovie - add a movie")
		msg = append(msg, "/addtv - add a tv show")
		msg = append(msg, "")
		msg = append(msg, "/cancel - cancel the current operation")
	}

	_ = s.Bot.Send(user, strings.Join(msg, "\n")) // @todo: handle error
}
