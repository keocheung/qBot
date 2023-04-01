package main

import (
	"fmt"
	"log"
	"os"

	"qb-monitor/client"
	"qb-monitor/model"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}
	qbClient := client.NewQbClient(config)
	taskManager := client.NewTaskManager(qbClient)
	taskManager.Start()
}

func loadConfig() (model.Config, error) {
	webURL := os.Getenv("WEB_URL")
	if webURL == "" {
		return model.Config{}, fmt.Errorf("WEB_URL not set")
	}
	apiKey := os.Getenv("API_KEY")
	return model.Config{
		WebURL: webURL,
		APIKey: apiKey,
	}, nil
}
