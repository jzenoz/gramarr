package app

import (
	"github.com/tommy647/gramarr/internal/router"
	"github.com/tommy647/gramarr/internal/users"

	"gopkg.in/tucnak/telebot.v2"
)

// @todo: comment...
func SetupHandlers(r *router.Router, a *Service) {
	// Send all telegram messages to our custom router
	a.Bot.Handle(telebot.OnText, r.Route)

	// Commands
	r.HandleFunc("/auth", a.RequirePrivate(a.RequireAuth(users.UANone, a.HandleAuth)))
	r.HandleFunc("/start", a.WithUser(a.RequirePrivate(a.RequireAuth(users.UANone, a.HandleStart))))
	r.HandleFunc("/help", a.RequirePrivate(a.RequireAuth(users.UANone, a.HandleStart)))
	r.HandleFunc("/cancel", a.RequirePrivate(a.RequireAuth(users.UANone, a.HandleCancel)))
	r.HandleFunc("/addmovie", a.RequirePrivate(a.RequireAuth(users.UAMember, a.HandleAddMovie)))
	r.HandleFunc("/addtv", a.RequirePrivate(a.RequireAuth(users.UAMember, a.HandleAddTVShow)))
	r.HandleFunc("/users", a.RequirePrivate(a.RequireAuth(users.UAAdmin, a.HandleUsers)))

	// Catchall Command
	r.HandleFallback(a.RequirePrivate(a.RequireAuth(users.UANone, a.HandleFallback)))

	// Conversation Commands
	r.HandleConvoFunc("/cancel", a.HandleConvoCancel)
}
