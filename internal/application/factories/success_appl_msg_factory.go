package factories

import (
	"fmt"

	"alex.com/agent_application/internal/infrastructure/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SuccessApplMsgFactory struct {
	ChannelIdentifier string
}

func (fact SuccessApplMsgFactory) Create(appl *models.Application) *tgbotapi.MessageConfig {
	msg := tgbotapi.MessageConfig{}
	msg.ChannelUsername = fact.ChannelIdentifier

	var text string
	if appl.ApplType.String == "cooperation_advisor" {
		text = "<b>Новая заявка агента</b>"
		text += fmt.Sprintf("\nИмя агента: %s ", appl.AgentFio.String)
		text += fmt.Sprintf("\nНомер телефона агента : %s ", appl.AgentPhone.String)
		text += fmt.Sprintf("\nТелеграм ID агента : @%s", appl.TelegramId)
		text += fmt.Sprintf("\nФио клиента : %s", appl.PersonName.String)
		text += fmt.Sprintf("\nНомер телефона клиента : %s", appl.PersonPhone.String)
		text += fmt.Sprintf("\nEmail клиента : %s", appl.Email.String)
		text += fmt.Sprintf("\nГород клиента : %s", appl.City.String)
		text += fmt.Sprintf("\nCтрана автомобиля : %s", appl.Country.String)
		text += fmt.Sprintf("\nМарка автомобиля : %s", appl.Mark.String)
		text += fmt.Sprintf("\nМодель автомобиля : %s", appl.Model.String)
		text += fmt.Sprintf("\nГод выпуска автомобиля : %s", appl.CarYear.String)
		text += fmt.Sprintf("\nБюджет : %s", appl.Budget.String)
		text += fmt.Sprintf("\nДополнительная информация : %s", appl.AdditionalInfo.String)
	} else {
		text = "<b>Новая заявка советника</b>"
		text += fmt.Sprintf("\nИмя советника: %s ", appl.AgentFio.String)
		text += fmt.Sprintf("\nНомер телефона советника : %s ", appl.AgentPhone.String)
		text += fmt.Sprintf("\nТелеграм ID советника : @%s", appl.TelegramId)
		text += fmt.Sprintf("\nФио клиента : %s", appl.PersonName.String)
		text += fmt.Sprintf("\nНомер телефона клиента : %s", appl.PersonPhone.String)
		text += fmt.Sprintf("\nCтрана автомобиля : %s", appl.Country.String)
		text += fmt.Sprintf("\nХарактеристики автомобиля : %s", appl.Mark.String)
		text += fmt.Sprintf("\nБюджет : %s", appl.Budget.String)
	}
	button := tgbotapi.NewInlineKeyboardButtonURL("Перейти к диалогу в TG", fmt.Sprintf("https://t.me/%s", appl.TelegramId))

	replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(button),
	)
	msg.ParseMode = "HTML"
	msg.Text = text
	msg.ReplyMarkup = replyKeyboard

	return &msg
}

func NewSuccessApplMsgFactory(channelIdentifier string) *SuccessApplMsgFactory {
	return &SuccessApplMsgFactory{ChannelIdentifier: channelIdentifier}
}
