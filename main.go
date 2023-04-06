package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"qb-monitor/client"
	"qb-monitor/model"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}
	log.Printf("config loaded: %+v", config)
	qbClient := client.NewQbClient(config.WebURL, config.APIKey)
	taskManager := client.NewTaskManager(config, qbClient)
	taskManager.Start()
}

func loadConfig() (model.Config, error) {
	webURL := os.Getenv("WEB_URL")
	if webURL == "" {
		return model.Config{}, fmt.Errorf("WEB_URL not set")
	}
	return model.Config{
		WebURL:               webURL,
		APIKey:               os.Getenv("API_KEY"),
		RatioLimitTags:       strings.Split(os.Getenv("RATIO_LIMIT_TAGS"), ","),
		RatioLimitCatogories: strings.Split(os.Getenv("RATIO_LIMIT_CATOGORIES"), ","),
	}, nil
}
