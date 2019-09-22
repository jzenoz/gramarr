package app

import (
	"gopkg.in/tucnak/telebot.v2"
)

func (s *Service) HandleAddMovie(m *telebot.Message) {
	s.CM.StartConversation(NewAddMovieConversation(s), m)
}
