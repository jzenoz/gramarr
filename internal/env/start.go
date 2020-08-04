package env

import (
	"fmt"
	"strings"

	"github.com/memodota/gramarr/internal/util"
	"gopkg.in/tucnak/telebot.v2"
)

func (e *Env) HandleStart(m *telebot.Message) {

	user, exists := e.Users.User(m.Sender.ID)

	var msg []string
	msg = append(msg, fmt.Sprintf("Hello, I'm %s! Use these commands to control me:", e.Bot.Me.FirstName))

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

	util.Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
	var startkeyboard []string
	if !exists {
		startkeyboard = append(startkeyboard, "/auth")
	}

	if exists && user.IsAdmin() {
		startkeyboard = append(startkeyboard, "/users")
	}

	if exists && (user.IsMember() || user.IsAdmin()) {
		startkeyboard = append(startkeyboard, "/addmovie")
		startkeyboard = append(startkeyboard, "/addtv")
		startkeyboard = append(startkeyboard, "/cancel")
	}
	util.SendKeyboardList(e.Bot, m.Sender, "Select command", startkeyboard)
}
