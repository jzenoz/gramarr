package app

import (
	"github.com/tommy647/gramarr/internal/router"
	"github.com/tommy647/gramarr/internal/users"

	"gopkg.in/tucnak/telebot.v2"
)

// @todo: comment...
func (s *Service) SetupHandlers(r *router.Router) {
	// Send all telegram messages to our custom router
	s.Bot.Handle(telebot.OnText, r.Route)

	// Commands
	r.HandleFunc("/auth", s.WithMiddleWare(users.UANone, s.HandleAuth))
	r.HandleFunc("/start", s.WithMiddleWare(users.UANone, s.HandleStart))
	r.HandleFunc("/help", s.WithMiddleWare(users.UANone, s.HandleStart))
	r.HandleFunc("/cancel", s.WithMiddleWare(users.UANone, s.HandleCancel))
	r.HandleFunc("/addmovie", s.WithMiddleWare(users.UAMember, s.HandleAddMovie))
	r.HandleFunc("/addtv", s.WithMiddleWare(users.UAMember, s.HandleAddTVShow))
	r.HandleFunc("/users", s.WithMiddleWare(users.UAAdmin, s.HandleUsers))

	// Catchall Command
	r.HandleFallback(s.WithMiddleWare(users.UANone, s.HandleFallback))

	// Conversation Commands
	r.HandleConvoFunc("/cancel", s.HandleConvoCancel)
}

func (s *Service) WithMiddleWare(access users.UserAccess, h router.Handler) router.Handler {
	return s.RequirePrivate(s.RequireAuth(access, h))
}
