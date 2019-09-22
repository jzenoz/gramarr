package bot

import "github.com/tommy647/gramarr/internal/users"

var _ Bot = (*Service)(nil)

type Bot interface {
	Start()
	Send(users.User, interface{}) error
	SendKeyboardList(users.User, string, []string) error
	Handle(interface{}, interface{})
	GetUserID(interface{}) interface{}
	IsPrivate(interface{}) bool
	GetText(interface{}) string
	GetPayload(interface{}) interface{}
}

// Service our bot service struct
type Service struct {
	name   string
	bot    Bot
	admins []users.User
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
func (s Service) GetText(msg interface{}) string                    { return s.bot.GetText(msg) }
func (s *Service) Start()                                           { s.bot.Start() }
func (s Service) Name() string                                      { return s.name }
func (s Service) IsPrivate(message interface{}) bool                { return s.bot.IsPrivate(message) }
func (s Service) GetUserID(message interface{}) interface{}         { return s.bot.GetUserID(message) }
func (s *Service) Handle(endpoint interface{}, handler interface{}) { s.bot.Handle(endpoint, handler) }
func (s *Service) GetPayload(msg interface{}) interface{}           { return s.bot.GetPayload(msg) }

// Send sends a message to recipient
func (s Service) Send(to users.User, msg interface{}) error { return s.bot.Send(to, msg) }

// SendError sends an error message to recipient @todo: why is this different?
func (s Service) SendError(to users.User, msg interface{}) error { return s.Send(to, msg) }

// Send sends a message to recipient
func (s Service) SendToAdmins(msg interface{}) error {
	for _, to := range s.admins {
		if err := s.bot.Send(to, msg); err != nil {
			// log error
		}
	}
	return nil
}

func (s Service) SendKeyboardList(to users.User, msg string, list []string) error {
	return s.bot.SendKeyboardList(to, msg, list)
}
