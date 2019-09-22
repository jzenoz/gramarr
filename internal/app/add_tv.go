package app

func (s *Service) HandleAddTVShow(m interface{}) {
	s.CM.StartConversation(NewAddTVShowConversation(s), m)
}
