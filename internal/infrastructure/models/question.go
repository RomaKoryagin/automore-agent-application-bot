package models

import "database/sql"

type Question struct {
	ID               int            `json:"id"`
	ChatId           int64          `json:"chat_id"`
	TelegramId       string         `json:"telegram_id"`
	Text             sql.NullString `json:"descript"`
	SendedToTelegram bool           `json:"sended_to_telegram"`
	SendedToBitrix   bool           `json:"sended_to_bitrix"`
}
