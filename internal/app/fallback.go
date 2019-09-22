package app

import (
	"strings"

	"github.com/tommy647/gramarr/internal/util"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Service) HandleFallback(m *tb.Message) {
	var msg []string
	msg = append(msg, "I'm sorry, I don't recognize that command.")
	msg = append(msg, "Type /help to see the available bot commands.")
	util.Send(s.Bot, m.Sender, strings.Join(msg, "\n"))
}
