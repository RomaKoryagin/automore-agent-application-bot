package factories

import (
	"alex.com/agent_application/internal/infrastructure/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MsgFactory struct {
}

func (fact MsgFactory) Create(msg models.Message, chatId int64) *tgbotapi.MessageConfig {
	tMsg := tgbotapi.NewMessage(chatId, msg.Text)
	if len(msg.ButtonRows) != 0 {
		var btns [][]tgbotapi.InlineKeyboardButton
		for _, row := range msg.ButtonRows {
			var rowButtons []tgbotapi.InlineKeyboardButton
			for _, el := range row.Buttons {
				rowButtons = append(
					rowButtons,
					tgbotapi.NewInlineKeyboardButtonData(el.Text, el.Command),
				)
			}
			btns = append(btns, rowButtons)
		}
		replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(btns...)
		tMsg.ReplyMarkup = replyKeyboard
	} else {
		replyKeyboard := tgbotapi.NewRemoveKeyboard(true)
		tMsg.ReplyMarkup = replyKeyboard
	}
	return &tMsg
}

func NewMsgFactory() *MsgFactory {
	return &MsgFactory{}
}
