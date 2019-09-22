package app

import (
	"strings"

	"github.com/tommy647/gramarr/internal/message"
)

func (s *Service) HandleFallback(m *message.Message) {

	var msg []string
	msg = append(msg, "I'm sorry, I don't recognize that command.")
	msg = append(msg, "Type /help to see the available bot commands.")
	_ = s.Bot.Send(m.User, strings.Join(msg, "\n")) // @todo: handle error
}
