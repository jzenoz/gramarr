package app

import (
	"github.com/tommy647/gramarr/internal/message"
	"github.com/tommy647/gramarr/internal/router"
	"github.com/tommy647/gramarr/internal/users"

	"gopkg.in/tucnak/telebot.v2"
)

// @todo: comment...
func (s *Service) SetupHandlers(r *router.Router) {
	// Send all telegram messages to our custom router
	s.Bot.Handle(telebot.OnText, r.Route)

	// Commands
	r.HandleFunc("/auth", s.WithAccessMiddleWare(users.UANone, s.HandleAuth))
	r.HandleFunc("/start", s.WithAccessMiddleWare(users.UANone, s.HandleStart))
	r.HandleFunc("/help", s.WithAccessMiddleWare(users.UANone, s.HandleStart))
	r.HandleFunc("/cancel", s.WithAccessMiddleWare(users.UANone, s.HandleCancel))
	r.HandleFunc("/addmovie", s.WithAccessMiddleWare(users.UAMember, s.HandleAddMovie))
	r.HandleFunc("/addtv", s.WithAccessMiddleWare(users.UAMember, s.HandleAddTVShow))
	r.HandleFunc("/users", s.WithAccessMiddleWare(users.UAAdmin, s.HandleUsers))

	// Catchall Command
	r.HandleFallback(s.WithAccessMiddleWare(users.UANone, s.HandleFallback))

	// Conversation Commands
	r.HandleConvoFunc("/cancel", s.HandleConvoCancel)
}

func (s *Service) WithAccessMiddleWare(access users.UserAccess, h router.Handler) router.Handler {
	return s.WithMiddleware(s.RequireAuth(access, h))
}

type Middleware func(func(*message.Message)) func(*message.Message)

func (s *Service) WithMiddleware(h router.Handler) router.Handler {
	mw := []Middleware{
		s.WithUser,
		s.RequirePrivate,
	}

	for _, m := range mw {
		h = m(h)
	}
	return h
}
