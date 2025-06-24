package factories

import (
	"fmt"

	"alex.com/agent_application/internal/infrastructure/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type QuestionMsgFactory struct {
	ChannelIdentifier string
}

func (fact QuestionMsgFactory) Create(question *models.Question) *tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	msg.ChannelUsername = fact.ChannelIdentifier

	text := "<b>Получен вопрос от пользователя</b>"
	text += fmt.Sprintf("\nTelegram ID пользовтеля: @%s ", question.TelegramId)
	text += fmt.Sprintf("\nВопрос пользователя: %s ", question.Text.String)

	button := tgbotapi.NewInlineKeyboardButtonURL("Перейти к диалогу в TG", fmt.Sprintf("https://t.me/%s", question.TelegramId))

	replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(button),
	)
	msg.ParseMode = "HTML"
	msg.Text = text
	msg.ReplyMarkup = replyKeyboard

	return &msg
}

func NewQuestionMsgFactory(channelIdentifier string) *QuestionMsgFactory {
	return &QuestionMsgFactory{ChannelIdentifier: channelIdentifier}
}
