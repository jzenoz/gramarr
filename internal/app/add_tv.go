package app

import "github.com/tommy647/gramarr/internal/message"

func (s *Service) HandleAddTVShow(m *message.Message) {
	s.CM.StartConversation(NewAddTVShowConversation(s), m)
}
