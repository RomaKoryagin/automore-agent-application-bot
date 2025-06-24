package repositories

import (
	"database/sql"
	"log"

	"alex.com/agent_application/internal/infrastructure/models"
)

type QuestionRepository struct {
	Connection *sql.DB
}

func (repo QuestionRepository) GetLastWithEmptyText(chatId int64) (*models.Question, error) {
	query := `
		select 
			id, 
			chat_id, 
			telegram_id, 
			descript, 
			sended_to_telegram, 
			sended_to_bitfix 
		from 
			questions
		where descript is null 
	`

	row := repo.Connection.QueryRow(query, chatId)

	var question models.Question

	err := row.Scan(
		&question.ID,
		&question.ChatId,
		&question.TelegramId,
		&question.Text,
		&question.SendedToTelegram,
		&question.SendedToBitrix,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &question, nil
}

func (repo QuestionRepository) CreateEmpty(chatId int64, telegramId string) error {
	sql := `insert into questions(chat_id, telegram_id) values (?, ?)`
	_, err := repo.Connection.Exec(sql, chatId, telegramId)
	if err != nil {
		log.Printf("error while trying to create empty question, more: %s", err)
		return err
	}

	return nil
}

func (repo QuestionRepository) UpdateText(chatId int64, text string) error {
	sql := `update questions set descript = ? where chat_id = ?`
	_, err := repo.Connection.Exec(sql, text, chatId)
	if err != nil {
		log.Printf("error while trying to update question text, more: %s", err)
		return err
	}

	return nil
}

func (repo QuestionRepository) SetSendedToTg(id int) error {
	sql := `update questions set sended_to_telegram = true where id = ?`
	_, err := repo.Connection.Exec(sql, id)
	if err != nil {
		log.Printf("error while trying to set sended to tg, more: %s", err)
		return err
	}

	return nil
}

func (repo QuestionRepository) GetNotSendedToTg() ([]*models.Question, error) {
	sql := `
		select 
			id, 
			chat_id, 
			telegram_id, 
			descript, 
			sended_to_telegram, 
			sended_to_bitfix 
		from 
			questions
		where not sended_to_telegram and descript is not null
	`
	rows, err := repo.Connection.Query(sql)

	if err != nil {
		if rows != nil {
			defer rows.Close()
		}
		return nil, err
	}

	var questions []*models.Question

	for rows.Next() {
		var appl models.Question
		err := rows.Scan(
			&appl.ID,
			&appl.ChatId,
			&appl.TelegramId,
			&appl.Text,
			&appl.SendedToTelegram,
			&appl.SendedToBitrix,
		)

		if err != nil {
			return nil, err
		}

		questions = append(questions, &appl)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

func NewQuestionRepository(conn *sql.DB) *QuestionRepository {
	return &QuestionRepository{Connection: conn}
}
