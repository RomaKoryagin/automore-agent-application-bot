package loader

import (
	"encoding/json"
	"log"
	"os"

	"alex.com/agent_application/internal/infrastructure/models"
)

type MessageLoader struct {
}

func (loader *MessageLoader) Load(path string) *models.MessageContainer {
	configFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("error while trying to open file: %v\n", err)
	}
	defer configFile.Close()

	var config models.MessageContainer
	decoder := json.NewDecoder(configFile)

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("eror while tryign to decode JSON: %v\n", err)
	}
	return &config
}

func NewMessageLoader() *MessageLoader {
	return &MessageLoader{}
}
