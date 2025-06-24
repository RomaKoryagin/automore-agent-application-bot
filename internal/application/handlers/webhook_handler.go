package handlers

import (
	"log"
	"net/http"

	"alex.com/agent_application/internal/application/commands"
	"github.com/gin-gonic/gin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WebhookHandler struct {
	Resolver *commands.Resolver
	Bot      *tgbotapi.BotAPI
}

func (handler *WebhookHandler) Handle(c *gin.Context) {
	var update tgbotapi.Update

	if err := c.BindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error while binding JSON": err.Error()})
		return
	}

	var msg string
	var chatId int64
	var username string

	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.From.ID
		msg = update.CallbackQuery.Data
		username = update.CallbackQuery.From.UserName
	}

	if update.Message != nil {
		msg = update.Message.Text
		chatId = update.Message.Chat.ID
		username = update.Message.From.UserName
	}

	rMsg := handler.Resolver.Resolve(chatId, username, msg)

	if rMsg != nil {
		log.Printf("sending msg: %s, chat id : %d", rMsg.Text, rMsg.ChatID)

		handler.Bot.Send(rMsg)

	}
}

func NewWebhookHandler(resolver *commands.Resolver, bot *tgbotapi.BotAPI) *WebhookHandler {
	return &WebhookHandler{Resolver: resolver, Bot: bot}
}
