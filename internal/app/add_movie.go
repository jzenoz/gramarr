package app

func (s *Service) HandleAddMovie(m interface{}) {
	s.CM.StartConversation(NewAddMovieConversation(s), m)
}
