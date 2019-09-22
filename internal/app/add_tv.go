package app

import (
	"gopkg.in/tucnak/telebot.v2"
)

func (s *Service) HandleAddTVShow(m *telebot.Message) {
	s.CM.StartConversation(NewAddTVShowConversation(s), m)
}
