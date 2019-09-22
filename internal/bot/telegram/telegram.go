package telegram

import (
	"encoding/json"
	"time"

	"github.com/tommy647/gramarr/internal/bot"
	"github.com/tommy647/gramarr/internal/users"

	"gopkg.in/tucnak/telebot.v2"
)

var _ bot.Bot = (*Service)(nil)

type Service struct {
	bot *telebot.Bot
}

func New(cfg Config) (*Service, error) {
	b, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.Token,
		Poller: &telebot.LongPoller{Timeout: cfg.Timeout},
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		bot: b,
	}, nil
}

func (s Service) Start() {
	s.bot.Start()
}

func (s *Service) Handle(endpoint interface{}, handler interface{}) {
	s.bot.Handle(endpoint, handler)
}

// Send a message to a user via telegram
// @todo: implement some type switching on the message so we can send strings and photos
func (s Service) Send(to users.User, msg interface{}) error {
	_, err := s.bot.Send(to, msg, telebot.ModeMarkdown)
	return err
}

func (s Service) SendKeyboardList(to users.User, msg string, list []string) error {
	var buttons []telebot.ReplyButton
	for _, item := range list {
		buttons = append(buttons, telebot.ReplyButton{Text: item})
	}

	var replyKeys [][]telebot.ReplyButton
	for _, b := range buttons {
		replyKeys = append(replyKeys, []telebot.ReplyButton{b})
	}

	_, err := s.bot.Send(to, msg, &telebot.ReplyMarkup{
		ReplyKeyboard:   replyKeys,
		OneTimeKeyboard: true,
	})

	return err
}

type Config struct {
	Token   string        `json:"token"`
	Timeout time.Duration `json:"timeout"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := struct {
		Timeout string `json:"timeout"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	timeout, err := time.ParseDuration(aux.Timeout)
	if err != nil {
		return err
	}
	c.Timeout = timeout
	return nil
}
