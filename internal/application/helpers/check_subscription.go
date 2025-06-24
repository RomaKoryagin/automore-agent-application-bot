package helpers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CheckUserSubscription(bot *tgbotapi.BotAPI, group string, userId int64) (bool, error) {
	var config tgbotapi.GetChatMemberConfig
	config.SuperGroupUsername = group
	config.UserID = userId
	member, err := bot.GetChatMember(config)
	if err != nil {
		log.Printf("error while getting chat member, more: %s", err)
		return false, err
	}
	return member.Status == "member" || member.Status == "administrator" || member.Status == "creator", nil
}
