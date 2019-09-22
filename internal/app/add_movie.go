package app

import "github.com/tommy647/gramarr/internal/message"

func (s *Service) HandleAddMovie(m *message.Message) {
	s.CM.StartConversation(NewAddMovieConversation(s), m)
}
