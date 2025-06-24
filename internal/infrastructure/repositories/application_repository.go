package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"alex.com/agent_application/internal/infrastructure/models"
)

type ApplicationRepository struct {
	Connection *sql.DB
}

func (repo ApplicationRepository) GetSubmittedAndTgNotSended() ([]*models.Application, error) {
	query := `
		select 
			id,
			chat_id,
			telegram_id,
			country,
			mark_or_conditions,
			budget,
			steering_wheel_type,
			city,
			person_name,
			person_phone,
			step,
			created_at, 
			updated_at,
			sended_telegram,
			sended_bitrix,
			mark,
			model,
			email,
			additional_info,
			car_year,
			appl_type,
			agent_fio,
			agent_phone
		from
			applications 
		where not sended_telegram and step = 'agent_appl_success'
		order by id desc
	`

	rows, err := repo.Connection.Query(query)

	if err != nil {
		if rows != nil {
			defer rows.Close()
		}
		return nil, err
	}
	var appls []*models.Application

	for rows.Next() {
		var appl models.Application
		err := rows.Scan(
			&appl.ID,
			&appl.ChatId,
			&appl.TelegramId,
			&appl.Country,
			&appl.MarkOrConditions,
			&appl.Budget,
			&appl.SteeringWheelType,
			&appl.City,
			&appl.PersonName,
			&appl.PersonPhone,
			&appl.Step,
			&appl.CreatedAt,
			&appl.UpdatedAt,
			&appl.SendedToTelegram,
			&appl.SendedToBitfix,
			&appl.Mark,
			&appl.Model,
			&appl.Email,
			&appl.AdditionalInfo,
			&appl.CarYear,
			&appl.ApplType,
			&appl.AgentFio,
			&appl.AgentPhone,
		)

		if err != nil {
			return nil, err
		}

		appls = append(appls, &appl)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appls, nil
}

func (repo ApplicationRepository) CreateEmpty(chatId int64, telegramId string) error {
	sql := `insert into applications(chat_id, telegram_id, created_at, updated_at) values (?, ?, datetime('now'), datetime('now'))`

	_, err := repo.Connection.Exec(sql, chatId, telegramId)

	if err != nil {
		log.Printf("error while trying to create empty application, more: %s", err)
		return err
	}

	return nil
}

func (repo ApplicationRepository) GetLastByUserId(userId int64) (*models.Application, error) {
	query := `
		select 
			id,
			chat_id,
			telegram_id,
			country,
			mark_or_conditions,
			budget,
			steering_wheel_type,
			city,
			person_name,
			person_phone,
			agent_fio,
			agent_phone,
			step,
			created_at, 
			updated_at,
			sended_telegram,
			sended_bitrix,
			mark,
			model,
			email,
			additional_info,
			car_year,
			appl_type
		from
			applications 
		where chat_id = ?
		order by id desc
		limit 1
	`

	row := repo.Connection.QueryRow(query, userId)

	var appl models.Application

	err := row.Scan(
		&appl.ID,
		&appl.ChatId,
		&appl.TelegramId,
		&appl.Country,
		&appl.MarkOrConditions,
		&appl.Budget,
		&appl.SteeringWheelType,
		&appl.City,
		&appl.PersonName,
		&appl.PersonPhone,
		&appl.AgentFio,
		&appl.AgentPhone,
		&appl.Step,
		&appl.CreatedAt,
		&appl.UpdatedAt,
		&appl.SendedToTelegram,
		&appl.SendedToBitfix,
		&appl.Mark,
		&appl.Model,
		&appl.Email,
		&appl.AdditionalInfo,
		&appl.CarYear,
		&appl.ApplType,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &appl, nil
}

func (repo ApplicationRepository) Update(appl *models.Application) error {
	sql := `
        update applications 
        set 
            country = ?,
            mark_or_conditions = ?,
            budget = ?,
            steering_wheel_type = ?,
            city = ?,
            person_name = ?,
            person_phone = ?,
			agent_fio = ?,
			agent_phone = ?,
            step = ?,
            updated_at = datetime('now'),
			sended_telegram = ?,
			sended_bitrix = ?,
			chat_id = ?,
			telegram_id = ?,
			mark = ?,
			model = ?,
			email = ?,
			additional_info = ?,
			car_year = ?,
			appl_type = ?,
			subscription_checked = ?
        where 
            id = ?
    `

	result, err := repo.Connection.Exec(sql,
		appl.Country,
		appl.MarkOrConditions,
		appl.Budget,
		appl.SteeringWheelType,
		appl.City,
		appl.PersonName,
		appl.PersonPhone,
		appl.AgentFio,
		appl.AgentPhone,
		appl.Step,
		appl.SendedToTelegram,
		appl.SendedToBitfix,
		appl.ChatId,
		appl.TelegramId,
		appl.Mark,
		appl.Model,
		appl.Email,
		appl.AdditionalInfo,
		appl.CarYear,
		appl.ApplType,
		appl.SubscriptionChecked,
		appl.ID,
	)

	if err != nil {
		log.Printf("error while executing update, more: %s", err)
		return fmt.Errorf("error while executing update, more: %s", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while tryning to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("there is no rows affected for application: %d", appl.ID)
	}

	return nil
}

func NewApplicationRepository(conn *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{Connection: conn}
}
