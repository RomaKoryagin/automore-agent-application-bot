package models

import "database/sql"

type Application struct {
	ID                  int            `json:"id"`
	ChatId              int64          `json:"chat_id"`
	TelegramId          string         `json:"telegram_id"`
	Country             sql.NullString `json:"country"`
	MarkOrConditions    sql.NullString `json:"mark_or_conditions"`
	Budget              sql.NullString `json:"budget"`
	SteeringWheelType   sql.NullString `json:"steering_wheel_type"`
	City                sql.NullString `json:"city"`
	AgentFio            sql.NullString `json:"agent_fio"`
	AgentPhone          sql.NullString `json:"agent_phone"`
	PersonName          sql.NullString `json:"person_name"`
	PersonPhone         sql.NullString `json:"person_phone"`
	Step                sql.NullString `json:"step"`
	ApplType            sql.NullString `json:"appl_type"`
	Mark                sql.NullString `json:"mark"`
	Model               sql.NullString `json:"model"`
	Email               sql.NullString `json:"email"`
	CarYear             sql.NullString `json:"car_year"`
	AdditionalInfo      sql.NullString `json:"additional_info"`
	SubscriptionChecked bool           `json:"subscription_checked"`
	UpdatedAt           string         `json:"updated_at"`
	CreatedAt           string         `json:"created_at"`
	SendedToTelegram    bool           `json:"sended_telegram"`
	SendedToBitfix      bool           `json:"sended_bitrix"`
}
