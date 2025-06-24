package commands

import (
	"log"

	"alex.com/agent_application/internal/application/factories"
	"alex.com/agent_application/internal/application/helpers"
	"alex.com/agent_application/internal/infrastructure/models"
	"alex.com/agent_application/internal/infrastructure/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Resolver struct {
	MsgContainer *models.MessageContainer
	MsgFactory   *factories.MsgFactory
	ApplRepo     *repositories.ApplicationRepository
	QuestionRepo *repositories.QuestionRepository
	Bot          *tgbotapi.BotAPI
}

func (reslv Resolver) Resolve(chatId int64, username string, text string) *tgbotapi.MessageConfig {
	question, err := reslv.QuestionRepo.GetLastWithEmptyText(chatId)
	if err != nil {
		log.Printf("error while trying to get last user question, more: %s", err)
		return nil
	}
	log.Println(question)
	if question != nil {
		reslv.QuestionRepo.UpdateText(chatId, text)
		if reslv.MsgContainer != nil {
			rMsg := (*reslv.MsgContainer).Messages["question_success"]
			msg := reslv.MsgFactory.Create(rMsg, chatId)
			return msg
		} else {
			return nil
		}
	}

	if text == "/question" {
		reslv.QuestionRepo.CreateEmpty(chatId, username)
		if reslv.MsgContainer != nil {
			rMsg := (*reslv.MsgContainer).Messages["question_greeting"]
			msg := reslv.MsgFactory.Create(rMsg, chatId)
			return msg
		} else {
			return nil
		}
	}

	var msgType string
	appl, err := reslv.ApplRepo.GetLastByUserId(chatId)
	if err != nil {
		log.Printf("error while getting last user by chat id, more: %s", err)
		return nil
	}

	if appl == nil {
		err = reslv.ApplRepo.CreateEmpty(chatId, username)
		if err != nil {
			log.Printf("error while creating new appl, more: %s", err)
			msgType = "error"
		} else {
			msgType = "gretting"
		}
	}
	log.Printf("type is : %s", msgType)
	if msgType == "" {
		resolved := false
		switch text {
		case "/menu":
			msgType = "menu"
		case "/start":
			err = reslv.ApplRepo.CreateEmpty(chatId, username)
			if err != nil {
				log.Printf("error update %s", err)
				msgType = "error"
			} else {
				msgType = "gretting"
			}
			resolved = true

		case "/cooperation_advisor":
			if appl != nil {
				appl.ApplType.String = "cooperation_advisor"
				appl.ApplType.Valid = true
				appl.Step.String = "subscription"
				appl.Step.Valid = true

				err = reslv.ApplRepo.Update(appl)
				if err != nil {
					log.Printf("create error, more: %s", err)
					msgType = "error"
				} else {
					msgType = "subscribing"
				}
			}
			resolved = true
		case "/cooperation_agent":
			appl.ApplType.String = "cooperation_agent"
			appl.ApplType.Valid = true
			appl.Step.String = "subscription"
			appl.Step.Valid = true
			err = reslv.ApplRepo.Update(appl)
			if err != nil {
				log.Printf("error update, more: %s", err)
				msgType = "error"
			} else {
				msgType = "subscribing"
			}
			resolved = true
		case "/subscibed_yes":
			// @TODO
			isSubscribed, err := helpers.CheckUserSubscription(reslv.Bot, "@automorevlad", chatId)
			if err != nil {
				log.Printf("error while trying to check subscription %s", err)
				msgType = "error"
			}

			if !isSubscribed {
				msgType = "no_subscription"
			} else {

				if appl.Step.String == "subscription" {
					if appl.ApplType.String == "cooperation_agent" {
						msgType = "agent_fio_agent"
						appl.Step.String = "agent_fio_agent"
						appl.Step.Valid = true
					} else {
						msgType = "agent_fio_advisor"
						appl.Step.String = "agent_fio_advisor"
						appl.Step.Valid = true
					}

					err = reslv.ApplRepo.Update(appl)
					if err != nil {
						msgType = "error"
					}
				}
			}
			resolved = true
		case "/subscibed_no":
			msgType = "no_subscription"
			appl.Step.String = "subscription"
			appl.Step.Valid = true
			resolved = true
		}
		if !resolved {
			msgType = appl.Step.String
			switch msgType {
			// start agent cases
			case "agent_fio_agent":
				appl.AgentFio.String = text
				appl.AgentFio.Valid = true
				appl.Step.String = "agent_fio_phone"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "agent_fio_phone"
			case "agent_fio_phone":
				isValid := helpers.IsPhoneValid(text)
				if !isValid {
					msgType = "invalid_phone"
				} else {
					appl.AgentPhone.String = text
					appl.AgentPhone.Valid = true
					appl.Step.String = "agent_person_name"
					appl.Step.Valid = true
					reslv.ApplRepo.Update(appl)
					msgType = "agent_person_name"
				}
			case "agent_person_name":
				appl.PersonName.String = text
				appl.PersonName.Valid = true
				appl.Step.String = "agent_person_phone"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "agent_person_phone"
			case "agent_person_phone":
				isValid := helpers.IsPhoneValid(text)
				if !isValid {
					msgType = "invalid_phone"
				} else {
					appl.PersonPhone.String = text
					appl.PersonPhone.Valid = true
					appl.Step.String = "agent_country"
					appl.Step.Valid = true
					reslv.ApplRepo.Update(appl)
					msgType = "agent_country"
				}
			case "agent_country":
				appl.Country.String = text
				appl.Country.Valid = true
				appl.Step.String = "agent_mark_or_conditions"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "agent_mark_or_conditions"
			case "agent_mark_or_conditions":
				appl.MarkOrConditions.String = text
				appl.MarkOrConditions.Valid = true
				appl.Step.String = "agent_budget"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "agent_budget"
			case "agent_budget":
				appl.Budget.String = text
				appl.Budget.Valid = true
				appl.Step.String = "agent_appl_success"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "agent_appl_success"
			// Start advisor cases
			case "agent_fio_advisor":
				appl.AgentFio.String = text
				appl.AgentFio.Valid = true
				appl.Step.String = "agent_phone_advisor"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "agent_phone_advisor"
			case "agent_phone_advisor":
				isValid := helpers.IsPhoneValid(text)
				if !isValid {
					msgType = "invalid_phone"
				} else {
					appl.AgentPhone.String = text
					appl.AgentPhone.Valid = true
					appl.Step.String = "advisor_person_name"
					appl.Step.Valid = true
					reslv.ApplRepo.Update(appl)
					msgType = "advisor_person_name"
				}
			case "advisor_person_name":
				appl.PersonName.String = text
				appl.PersonName.Valid = true
				appl.Step.String = "advisor_person_phone"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "advisor_person_phone"
			case "advisor_person_phone":
				isValid := helpers.IsPhoneValid(text)
				if !isValid {
					msgType = "invalid_phone"
				} else {
					appl.PersonPhone.String = text
					appl.PersonPhone.Valid = true
					appl.Step.String = "advisor_person_email"
					appl.Step.Valid = true
					reslv.ApplRepo.Update(appl)
					msgType = "advisor_person_email"
				}
			case "advisor_person_email":
				appl.Email.String = text
				appl.Email.Valid = true
				appl.Step.String = "advisor_person_city"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "advisor_person_city"
			case "advisor_person_city":
				appl.City.String = text
				appl.City.Valid = true
				appl.Step.String = "advisor_person_coutry"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "advisor_person_coutry"
			case "advisor_person_coutry":
				appl.Country.String = text
				appl.Country.Valid = true
				appl.Step.String = "advisor_person_mark"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "advisor_person_mark"
			case "advisor_person_mark":
				appl.Mark.String = text
				appl.Mark.Valid = true
				appl.Step.String = "advisor_person_model"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "advisor_person_model"
			case "advisor_person_model":
				appl.Model.String = text
				appl.Model.Valid = true
				appl.Step.String = "advisor_person_car_year"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "advisor_person_car_year"
			case "advisor_person_car_year":
				appl.CarYear.String = text
				appl.CarYear.Valid = true
				appl.Step.String = "advisor_person_budget"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "advisor_person_budget"
			case "advisor_person_budget":
				appl.Budget.String = text
				appl.Budget.Valid = true
				appl.Step.String = "advisor_person_additional_info"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "advisor_person_additional_info"
			case "advisor_person_additional_info":
				appl.AdditionalInfo.String = text
				appl.AdditionalInfo.Valid = true
				appl.Step.String = "agent_appl_success"
				appl.Step.Valid = true
				reslv.ApplRepo.Update(appl)
				msgType = "agent_appl_success"
			}
		}
	}
	if reslv.MsgContainer != nil {
		rMsg := (*reslv.MsgContainer).Messages[msgType]
		msg := reslv.MsgFactory.Create(rMsg, chatId)
		return msg
	} else {
		// msg type is undefined
		// @Todo
		return nil
	}
}

func NewResolver(
	msgContainer *models.MessageContainer,
	msgFactory *factories.MsgFactory,
	applRepo *repositories.ApplicationRepository,
	questionRepo *repositories.QuestionRepository,
	bot *tgbotapi.BotAPI,
) *Resolver {
	return &Resolver{
		MsgContainer: msgContainer,
		MsgFactory:   msgFactory,
		ApplRepo:     applRepo,
		QuestionRepo: questionRepo,
		Bot:          bot,
	}
}
