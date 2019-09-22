package app

import (
	"fmt"
	"strings"

	"github.com/tommy647/gramarr/internal/conversation"
	"github.com/tommy647/gramarr/internal/message"
)

func (s *Service) HandleCancel(m *message.Message) {
	_ = s.Bot.Send(m.User, "There is no active command to cancel. I wasn't doing anything anyway. Zzzzz...") // @todo: handle error
}

func (s *Service) HandleConvoCancel(c conversation.Conversation, m *message.Message) {
	s.CM.StopConversation(c)

	var msg []string
	msg = append(msg, fmt.Sprintf("The '*%s*' command was cancelled. Anything else I can do for you?", c.Name()))
	msg = append(msg, "")
	msg = append(msg, "Send /help for a list of commands.")
	_ = s.Bot.Send(m.User, strings.Join(msg, "\n")) // @todo: handle error
}
