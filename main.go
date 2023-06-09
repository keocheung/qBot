package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"qb-monitor/client"
	"qb-monitor/model"

	"github.com/antonmedv/expr"
)

const (
	webURLEnvKey      = "WEB_URL"
	configPathEnvKey  = "CONFIG_PATH"
	defaultConfigPath = "./config.json"
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
	webURL := os.Getenv(webURLEnvKey)
	if webURL == "" {
		return model.Config{}, fmt.Errorf("WEB_URL not set")
	}
	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		configPath = defaultConfigPath
	}
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return model.Config{}, fmt.Errorf("Read from config file error: %v", err)
	}
	var c model.Config
	if err := json.Unmarshal(configFile, &c); err != nil {
		return model.Config{}, fmt.Errorf("Unmarshal config file error: %v", err)
	}
	for _, v := range c.Rules {
		evaluator, err := expr.Compile(v.Condition)
		if err != nil {
			return model.Config{}, fmt.Errorf("Compile rule (%s) error: %v", v.Condition, err)
		}
		v.Evaluator = evaluator
	}
	return c, nil
}
