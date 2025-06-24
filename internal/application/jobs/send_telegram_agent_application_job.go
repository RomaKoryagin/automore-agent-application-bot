package jobs

import (
	"log"
	"time"

	"alex.com/agent_application/internal/application/factories"
	"alex.com/agent_application/internal/infrastructure/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SendTelegramAgentApplicationJob struct {
	ApplicationRepository *repositories.ApplicationRepository
	Bot                   *tgbotapi.BotAPI
	SuccessApplMsgFactory *factories.SuccessApplMsgFactory
}

func (job SendTelegramAgentApplicationJob) Execute() {
	go func() {
		for {
			appls, err := job.ApplicationRepository.GetSubmittedAndTgNotSended()
			if err != nil {
				log.Printf("Error while trying to get subbmited applications, more: %s", err)
				continue
			}

			for _, appl := range appls {
				msg := job.SuccessApplMsgFactory.Create(appl)

				_, err := job.Bot.Send(msg)
				if err != nil {
					log.Println(err)
					continue
				}
				appl.SendedToTelegram = true
				job.ApplicationRepository.Update(appl)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func NewSendTelegramAgentApplicationJob(
	applicationRepository *repositories.ApplicationRepository,
	bot *tgbotapi.BotAPI,
	fact *factories.SuccessApplMsgFactory,
) *SendTelegramAgentApplicationJob {
	return &SendTelegramAgentApplicationJob{
		ApplicationRepository: applicationRepository,
		Bot:                   bot,
		SuccessApplMsgFactory: fact,
	}
}
