package bot

type Options func(*Service)

// WithBot adds the telegram bot to the service
func WithBot(b Bot) Options {
	return func(s *Service) {
		s.bot = b
	}
}
