package router

import (
	"fmt"
	"regexp"

	"github.com/jzenoz/gramarr/internal/conversation"
	"gopkg.in/tucnak/telebot.v2"
)

var (
	cmdRx = regexp.MustCompile(`^(/\w+)(@(\w+))?(\s|$)(.+)?`)
)

type Handler func(*telebot.Message)
type ConvoHandler func(conversation.Conversation, *telebot.Message)

func NewRouter(cm *conversation.ConversationManager) *Router {
	return &Router{cm: cm, routes: map[string]Handler{}, convoRoutes: map[string]ConvoHandler{}}
}

type Router struct {
	cm          *conversation.ConversationManager
	routes      map[string]Handler
	convoRoutes map[string]ConvoHandler
	fallback    Handler
}

func (r *Router) HandleFunc(cmd string, h Handler) {
	r.routes[cmd] = h
}

func (r *Router) HandleFallback(h Handler) {
	r.fallback = h
}

func (r *Router) HandleConvoFunc(cmd string, h ConvoHandler) {
	r.convoRoutes[cmd] = h
}

func (r *Router) Route(m *telebot.Message) {
	if !r.routeConvo(m) && !r.routeCommand(m) {
		r.routeFallback(m)
	}
}

func (r *Router) routeConvo(m *telebot.Message) bool {
	fmt.Printf("%+v Command: ", m.Sender)
	fmt.Printf("%+v\n", m.Text)
	if !r.cm.HasConversation(m) {
		return false
	}

	// Global Conversation Cmd?
	if cmd, match := r.parseCommand(m); match {
		if route, exists := r.convoRoutes[cmd]; exists {
			convo, _ := r.cm.Conversation(m)
			route(convo, m)
			return true
		}
	}

	r.cm.ProcessMessage(m)
	return true
}

func (r *Router) routeCommand(m *telebot.Message) bool {
	if cmd, match := r.parseCommand(m); match {
		if route, exists := r.routes[cmd]; exists {
			route(m)
			return true
		}
	}
	return false
}

func (r Router) routeFallback(m *telebot.Message) {
	if r.fallback != nil {
		r.fallback(m)
	}
}

func (r *Router) parseCommand(m *telebot.Message) (string, bool) {
	match := cmdRx.FindAllStringSubmatch(m.Text, -1)

	if match != nil {
		return match[0][1], true
	} else {
		return "", false
	}
}
