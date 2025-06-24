package jobs

import (
	"log"
	"time"

	"alex.com/agent_application/internal/application/factories"
	"alex.com/agent_application/internal/infrastructure/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SendTelegramQuestionJob struct {
	Bot                *tgbotapi.BotAPI
	QuestionRepository *repositories.QuestionRepository
	QuestionMsgFactory *factories.QuestionMsgFactory
}

func (job SendTelegramQuestionJob) Execute() {
	go func() {
		for {
			questions, err := job.QuestionRepository.GetNotSendedToTg()
			if err != nil {
				log.Printf("Error while trying to get questions, more: %s", err)
				continue
			}

			for _, question := range questions {
				msg := job.QuestionMsgFactory.Create(question)

				_, err := job.Bot.Send(msg)
				if err != nil {
					log.Println(err)
					continue
				}

				job.QuestionRepository.SetSendedToTg(question.ID)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func NewSendTelegramQuestionJob(
	bot *tgbotapi.BotAPI,
	repo *repositories.QuestionRepository,
	fact *factories.QuestionMsgFactory,
) *SendTelegramQuestionJob {
	return &SendTelegramQuestionJob{Bot: bot, QuestionRepository: repo, QuestionMsgFactory: fact}
}
