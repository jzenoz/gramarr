package telegram

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/tommy647/gramarr/internal/bot"
	"github.com/tommy647/gramarr/internal/message"
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

func getTelegramMessage(message interface{}) (*telebot.Message, error) {
	switch message.(type) {
	case *telebot.Message:
		return message.(*telebot.Message), nil
	}
	return nil, errors.New("not a telegram message")
}

// WithMessage converts a telegram message to a standard message
func (s Service) WithMessage(msg interface{}) (*message.Message, error) {
	var m *telebot.Message
	switch msg.(type) {
	case *telebot.Message:
		m = msg.(*telebot.Message)
	default:
		return nil, errors.New("not a telegram message")
	}

	user := users.User{
		ID:        m.Sender.ID,
		Username:  m.Sender.Username,
		FirstName: m.Sender.FirstName,
		LastName:  m.Sender.LastName,
		Access:    users.UANone,
	}

	return &message.Message{
		Payload: m.Payload,
		Text:    m.Text,
		User:    user,
		Private: m.Private(),
		Chat: message.Chat{
			ID: m.Chat.ID,
		},
	}, nil
}

// @todo: we need to flesh out the type checking here
func (s Service) IsPrivate(message interface{}) bool {
	msg, err := getTelegramMessage(message)
	if err != nil {
		log.Println("IsPrivate", err.Error())
		return false // @todo: LOG
	}
	return msg.Private()
}

// @todo: we need to flesh out the type checking here
func (s Service) GetUserID(message interface{}) interface{} {
	msg, err := getTelegramMessage(message)
	if err != nil {
		log.Println("GetUserID", err.Error())
		return nil // @todo: LOG
	}
	return msg.Sender.ID
}

func (s Service) GetPayload(message interface{}) interface{} {
	msg, err := getTelegramMessage(message)
	if err != nil {
		log.Println("GetPayload", err.Error())
		return nil // @todo: LOG
	}
	return msg.Payload
}

func (s Service) GetText(message interface{}) string {
	msg, err := getTelegramMessage(message)
	if err != nil {
		log.Println("GetText", err.Error())
		return ""
	}
	return msg.Text
}

func (s Service) Start()                                            { s.bot.Start() }
func (s *Service) Handle(endpoint interface{}, handler interface{}) { s.bot.Handle(endpoint, handler) }

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
