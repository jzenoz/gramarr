package env

import (
	"strings"

	"github.com/jzenoz/gramarr/internal/util"
	"gopkg.in/tucnak/telebot.v2"
)

func (e *Env) HandleFallback(m *telebot.Message) {
	var msg []string
	msg = append(msg, "I'm sorry, I don't recognize that command.")
	msg = append(msg, "Type /help to see the available bot commands.")
	util.Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
}
