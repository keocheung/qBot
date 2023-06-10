package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"qb-monitor/client"
	"qb-monitor/model"
	"qb-monitor/util/logger"

	"github.com/antonmedv/expr"
	"github.com/fsnotify/fsnotify"
)

const (
	configPathEnvKey  = "CONFIG_PATH"
	defaultConfigPath = "./config.json"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	config, err := loadConfig(watcher)
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}
	logger.Infof("config loaded: %+v", config)
	qbClient := client.NewQbClient(config.WebURL, config.APIKey)
	taskManager := client.NewTaskManager(config, qbClient)
	taskManager.Start()
}

func loadConfig(watcher *fsnotify.Watcher) (*model.Config, error) {
	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		configPath = defaultConfigPath
	}

	conf := &model.Config{}
	go func(configPath string, conf *model.Config) {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logger.Infof("event:", event)
				if event.Has(fsnotify.Write) {
					c, err := loadConfigFromFile(configPath)
					if err != nil {
						logger.Errorf("error:", err)
					} else {
						*conf = c
						logger.Infof("Reload config")
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Errorf("error:", err)
			}
		}
	}(configPath, conf)

	err := watcher.Add(configPath)
	if err != nil {
		return nil, err
	}
	logger.Infof("watching config file: %s", configPath)

	c, err := loadConfigFromFile(configPath)
	if err != nil {
		return nil, err
	}
	*conf = c
	return conf, nil
}

func loadConfigFromFile(path string) (model.Config, error) {
	configFile, err := ioutil.ReadFile(path)
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
