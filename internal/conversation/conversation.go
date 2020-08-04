package conversation

import (
	"fmt"

	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/tucnak/telebot.v2"
)

type Conversation interface {
	Run(m *telebot.Message)
	CurrentStep() func(*telebot.Message)
	Name() string
}

type ConversationManager struct {
	convos *cache.Cache
}

func NewConversationManager() *ConversationManager {
	convos := cache.New(30*time.Minute, 10*time.Minute)
	return &ConversationManager{convos: convos}
}

func (cm *ConversationManager) ProcessMessage(m *telebot.Message) {
	key := cm.convoKey(m)
	if convo, ok := cm.convos.Get(key); ok {
		c := convo.(Conversation)
		c.CurrentStep()(m)
	}
}

func (cm *ConversationManager) HasConversation(m *telebot.Message) bool {
	_, exists := cm.convos.Get(cm.convoKey(m))
	return exists
}

func (cm *ConversationManager) StartConversation(c Conversation, m *telebot.Message) {
	c.Run(m)
	cm.convos.SetDefault(cm.convoKey(m), c)
}

func (cm *ConversationManager) StopConversation(c Conversation) {
	for key, item := range cm.convos.Items() {
		current := item.Object.(Conversation)
		if c == current {
			cm.convos.Delete(key)
		}
	}
}

func (cm *ConversationManager) Conversation(m *telebot.Message) (Conversation, bool) {
	c, exists := cm.convos.Get(cm.convoKey(m))
	return c.(Conversation), exists
}

func (cm *ConversationManager) convoKey(m *telebot.Message) string {
	return fmt.Sprintf("%d:%d", m.Chat.ID, m.Sender.ID)
}
