package app

import (
	"strings"

	"github.com/tommy647/gramarr/internal/users"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Service) HandleFallback(m *tb.Message) {
	user := users.User{} // @todo: get from context

	var msg []string
	msg = append(msg, "I'm sorry, I don't recognize that command.")
	msg = append(msg, "Type /help to see the available bot commands.")
	_ = s.Bot.Send(user, strings.Join(msg, "\n")) // @todo: handle error
}
