package client

import (
	"log"
	"strings"
	"sync"
	"time"

	"qb-monitor/model"
	"qb-monitor/util"
)

// TaskManager is the task manager of qb-monitor
type TaskManager interface {
	Start() // Start the task manager
}

// NewTaskManager creates a new task manager
func NewTaskManager(config model.Config, qbClient QbClient) TaskManager {
	taskManager := &taskManager{
		config:   config,
		qbClient: qbClient,
		tasks:    []func(config model.Config, qbClient QbClient) error{},
	}
	taskManager.tasks = []func(config model.Config, qbClient QbClient) error{
		limitShareRatio,
	}
	return taskManager
}

// Limit the share ratio to a certain value for certain torrents
var limitShareRatio = func(config model.Config, qbClient QbClient) error {
	log.Printf("task limitShareRatio started")
	torrents, err := qbClient.GetTorrents(model.Options{
		Limit:   10,
		Sort:    "added_on",
		Reverse: true,
	})
	if err != nil {
		log.Printf("get torrents error: %v", err)
		return err
	}
	for _, torrent := range torrents {
		if torrent.MaxRatio == -1 && (util.StringArraysHasCommon(config.RatioLimitTags, strings.Split(torrent.Tags, ",")) ||
			util.StringsContain(config.RatioLimitCatogories, torrent.Category)) {
			err := qbClient.SetShareLimits([]string{torrent.Hash}, 2.0, torrent.MaxSeedingTime)
			if err != nil {
				log.Printf("set share limit for %s error: %v", torrent.Hash, err)
			} else {
				log.Printf("set share limit to 2.0 for %s (%s)", torrent.Hash, torrent.Name)
			}
		}
	}
	log.Printf("task limitShareRatio finished")
	return nil
}

type taskManager struct {
	config   model.Config
	qbClient QbClient
	tasks    []func(config model.Config, qbClient QbClient) error
}

// Start starts the task manager
func (t *taskManager) Start() {
	log.Println("task manager started")
	var wg sync.WaitGroup
	for _, task := range t.tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				err := task(t.config, t.qbClient)
				if err != nil {
					log.Printf("task error: %v", err)
				}
				time.Sleep(10 * time.Minute)
			}
		}()
	}
	wg.Wait()
}
