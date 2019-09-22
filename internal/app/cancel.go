package app

import (
	"fmt"
	"strings"

	"github.com/tommy647/gramarr/internal/conversation"
	"github.com/tommy647/gramarr/internal/util"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Service) HandleCancel(m *tb.Message) {
	util.Send(s.Bot, m.Sender, "There is no active command to cancel. I wasn't doing anything anyway. Zzzzz...")
}

func (s *Service) HandleConvoCancel(c conversation.Conversation, m *tb.Message) {
	s.CM.StopConversation(c)

	var msg []string
	msg = append(msg, fmt.Sprintf("The '*%s*' command was cancelled. Anything else I can do for you?", c.Name()))
	msg = append(msg, "")
	msg = append(msg, "Send /help for a list of commands.")
	util.Send(s.Bot, m.Sender, strings.Join(msg, "\n"))
}
