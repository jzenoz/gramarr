package bot

import "github.com/tommy647/gramarr/internal/users"

var _ Bot = (*Service)(nil)

type Bot interface {
	Start()
	Send(users.User, string) error
}

// Service our bot service struct
type Service struct {
	name string
	bot  Bot
}

// New instantiates a new Service
func New(cfg Config, opts ...Options) *Service {
	s := &Service{
		name: cfg.Name,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Service) Start() {
	s.bot.Start()
}

// Send sends a message to recipient
func (s Service) Send(to users.User, msg string) error {
	return s.bot.Send(to, msg)
}

// SendError sends an error message to recipient
func (s Service) SendError(to users.User, msg string) error {
	return s.Send(to, msg) // @todo: why is this different?
}

// Send sends a message to recipient
func (s Service) SendAdmin(admins []users.User, msg string) error {
	for _, to := range admins {
		if err := s.bot.Send(to, msg); err != nil {
			// log error
		}
	}
	return nil
}
