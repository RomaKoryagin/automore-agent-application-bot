package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"alex.com/agent_application/internal/application/commands"
	"alex.com/agent_application/internal/application/factories"
	"alex.com/agent_application/internal/application/handlers"
	"alex.com/agent_application/internal/application/jobs"
	"alex.com/agent_application/internal/infrastructure/loader"
	"alex.com/agent_application/internal/infrastructure/repositories"
	"alex.com/agent_application/internal/infrastructure/sqlite"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// setting up environment vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// end setting up envronment vars

	// intitializing SQLite database
	executableFolderPath, _ := os.Getwd()
	db := sqlite.Database{MainDirPath: executableFolderPath}
	db.Init()
	// end initializing SQLite database

	// loading msgs from JSON
	absPath, _ := filepath.Abs(".")
	configPath := filepath.Join(absPath, "config", "message.config.json")
	msgLoader := loader.NewMessageLoader()
	msgs := msgLoader.Load(configPath)
	// end loading msgs from JSON

	// creating instance of tgBot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	// end creating isntance of tgBot

	// creating beans
	applRepo := repositories.NewApplicationRepository(db.Connection)
	questionRepo := repositories.NewQuestionRepository(db.Connection)

	msgFactory := factories.NewMsgFactory()
	successMsgFactory := factories.NewSuccessApplMsgFactory(os.Getenv("GROUP_TELEGRAM_CHAT_NAME"))
	questionMsgFactory := factories.NewQuestionMsgFactory(os.Getenv("GROUP_TELEGRAM_CHAT_NAME"))

	resolver := commands.NewResolver(msgs, msgFactory, applRepo, questionRepo, bot)
	webhookHandler := handlers.NewWebhookHandler(resolver, bot)
	// end creating beans

	// jobs
	jobs.NewSendTelegramAgentApplicationJob(applRepo, bot, successMsgFactory).Execute()
	jobs.NewSendTelegramQuestionJob(bot, questionRepo, questionMsgFactory).Execute()
	// end jobs

	r := gin.Default()
	r.POST("/v1/application-agent-bot/handle", webhookHandler.Handle)

	r.Run(fmt.Sprintf(":%s", os.Getenv("REST_API_PORT")))
}
