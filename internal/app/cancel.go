package app

import (
	"fmt"
	"strings"

	"github.com/tommy647/gramarr/internal/users"

	"github.com/tommy647/gramarr/internal/conversation"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Service) HandleCancel(m *tb.Message) {
	user := users.User{}                                                                                   // @todo: get from context
	_ = s.Bot.Send(user, "There is no active command to cancel. I wasn't doing anything anyway. Zzzzz...") // @todo: handle error
}

func (s *Service) HandleConvoCancel(c conversation.Conversation, m *tb.Message) {
	user := users.User{} // @todo: get from context
	s.CM.StopConversation(c)

	var msg []string
	msg = append(msg, fmt.Sprintf("The '*%s*' command was cancelled. Anything else I can do for you?", c.Name()))
	msg = append(msg, "")
	msg = append(msg, "Send /help for a list of commands.")
	_ = s.Bot.Send(user, strings.Join(msg, "\n")) // @todo: handle error
}
