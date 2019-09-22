package bot

import "github.com/tommy647/gramarr/internal/users"

type Options func(*Service)

// WithBot adds the instance bot to the service
func WithBot(b Bot) Options {
	return func(s *Service) {
		s.bot = b
	}
}

// WithAdmins adds our admins to the service
func WithAdmins(admins []users.User) Options {
	return func(s *Service) {
		s.admins = admins
	}
}
