package util

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/memodota/gramarr/internal/users"
	"gopkg.in/tucnak/telebot.v2"
)

func Send(bot *telebot.Bot, to telebot.Recipient, msg string) {
	bot.Send(to, msg, telebot.ModeMarkdown)
}

func SendError(bot *telebot.Bot, to telebot.Recipient, msg string) {
	bot.Send(to, msg, telebot.ModeMarkdown)
}

func SendAdmin(bot *telebot.Bot, to []users.User, msg string) {
	SendMany(bot, to, fmt.Sprintf("*[Admin]* %s", msg))
}

func SendKeyboardList(bot *telebot.Bot, to telebot.Recipient, msg string, list []string) {
	var buttons []telebot.ReplyButton
	for _, item := range list {
		buttons = append(buttons, telebot.ReplyButton{Text: item})
	}

	var replyKeys [][]telebot.ReplyButton
	for _, b := range buttons {
		replyKeys = append(replyKeys, []telebot.ReplyButton{b})
	}

	bot.Send(to, msg, &telebot.ReplyMarkup{
		ReplyKeyboard:   replyKeys,
		OneTimeKeyboard: true,
	})
}

func SendMany(bot *telebot.Bot, to []users.User, msg string) {
	for _, user := range to {
		bot.Send(user, msg, telebot.ModeMarkdown)
	}
}

func DisplayName(u *telebot.User) string {
	if u.FirstName != "" && u.LastName != "" {
		return EscapeMarkdown(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
	}

	return EscapeMarkdown(u.FirstName)
}

func EscapeMarkdown(s string) string {
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	s = strings.Replace(s, "_", "\\_", -1)
	return s
}

func BoolToYesOrNo(condition bool) string {
	if condition {
		return "Yes"
	}
	return "No"
}

func FormatDate(t time.Time) string {
	if t.IsZero() {
		return "Unknown"
	}
	return t.Format("02.01.2006")
}

func FormatDateTime(t time.Time) string {
	if t.IsZero() {
		return "Unknown"
	}
	return t.Format("02.01.2006 15:04:05")
}

func GetRootFolderFromPath(path string) string {
	return strings.Title(filepath.Base(filepath.Dir(path)))
}
