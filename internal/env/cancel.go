package env

import (
	"fmt"
	"strings"

	"github.com/memodota/gramarr/internal/conversation"
	"github.com/memodota/gramarr/internal/util"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (e *Env) HandleCancel(m *tb.Message) {
	util.Send(e.Bot, m.Sender, "There is no active command to cancel. I wasn't doing anything anyway. Zzzzz...")
}

func (e *Env) HandleConvoCancel(c conversation.Conversation, m *tb.Message) {
	var cancelkeyboard []string
	cancelkeyboard = append(cancelkeyboard, "/help")
	util.SendKeyboardList(e.Bot, m.Sender, "", cancelkeyboard)

	var msg []string
	msg = append(msg, fmt.Sprintf("The '*%s*' command was cancelled. Anything else I can do for you?", c.Name()))
	msg = append(msg, "")
	msg = append(msg, "Send /help for a list of commands.")
	util.Send(e.Bot, m.Sender, strings.Join(msg, "\n"))

	e.CM.StopConversation(c)
}
