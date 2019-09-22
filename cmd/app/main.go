package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/tommy647/gramarr/internal/bot/telegram"

	"github.com/tommy647/gramarr/internal/app"
	"github.com/tommy647/gramarr/internal/auth"
	"github.com/tommy647/gramarr/internal/bot"
	"github.com/tommy647/gramarr/internal/config"
	"github.com/tommy647/gramarr/internal/conversation"
	"github.com/tommy647/gramarr/internal/radarr"
	"github.com/tommy647/gramarr/internal/router"
	"github.com/tommy647/gramarr/internal/sonarr"
	"github.com/tommy647/gramarr/internal/users"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Flags
var configDir = flag.String("configDir", ".", "config dir for settings and logs")

func main() {
	flag.Parse()

	conf, err := config.LoadConfig(*configDir)
	if err != nil {
		log.Fatalf("failed to load config file: %s", err.Error())
	}

	//err = config.ValidateConfig(conf) // @todo: doesn't do anything
	//if err != nil {
	//	log.Fatal("config error: %s", err.Error())
	//}

	userPath := filepath.Join(*configDir, "users.json")
	usrs, err := users.NewUserDB(userPath)
	if err != nil {
		log.Fatalf("failed to load the usrs db %v", err)
	}

	var rc *radarr.Client
	if conf.Radarr != nil {
		rc, err = radarr.New(*conf.Radarr)
		if err != nil {
			log.Fatalf("failed to create radarr client: %v", err)
		}
	}

	var sn *sonarr.Client
	if conf.Sonarr != nil {
		sn, err = sonarr.New(*conf.Sonarr)
		if err != nil {
			log.Fatalf("failed to create sonarr client: %v", err)
		}
	}

	cm := conversation.NewConversationManager()
	r := router.NewRouter(cm)

	// @todo : move this into our bot service
	tbot, err := telegram.New(conf.Telegram)
	if err != nil {
		log.Fatalf("failed to create telegram bot client: %v", err)
	}

	boter := bot.New(conf.Bot, bot.WithBot(tbot))

	authoriser := auth.New(conf.Auth, auth.WithBot(boter), auth.WithUsers(usrs))

	a := &app.Service{
		Auth:   authoriser,
		Bot:    tbot,
		Users:  usrs,
		CM:     cm,
		Radarr: rc,
		Sonarr: sn,
	}

	setupHandlers(r, a)
	log.Print("Gramarr is up and running. Go call your bot!")
	boter.Start()
}

func setupHandlers(r *router.Router, a *app.Service) {
	// Send all telegram messages to our custom router
	a.Bot.Handle(tb.OnText, r.Route)

	// Commands
	r.HandleFunc("/auth", a.RequirePrivate(a.RequireAuth(users.UANone, a.HandleAuth)))
	r.HandleFunc("/start", a.RequirePrivate(a.RequireAuth(users.UANone, a.HandleStart)))
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
