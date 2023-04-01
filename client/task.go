package client

import (
	"log"
	"strings"
	"sync"
	"time"

	"qb-monitor/model"
)

// TaskManager is the task manager of qb-monitor
type TaskManager interface {
	Start() // Start the task manager
}

// NewTaskManager creates a new task manager
func NewTaskManager(qbClient QbClient) TaskManager {
	taskManager := &taskManager{
		qbClient: qbClient,
		tasks:    []func(qbClient QbClient){},
	}
	taskManager.tasks = []func(qbClient QbClient){
		limitShareRatio,
	}
	return taskManager
}

// Limit the share ratio to a certain value for certain torrents
var limitShareRatio = func(qbClient QbClient) {
	for {
		torrents, err := qbClient.GetTorrents(model.Options{
			Limit:   10,
			Sort:    "added_on",
			Reverse: true,
		})
		if err != nil {
			log.Printf("get torrents error: %v", err)
			continue
		}
		for _, torrent := range torrents {
			if torrent.MaxRatio == -1 {
				tags := strings.Split(torrent.Tags, ",")
				for _, tag := range tags {
					if tag == "VCB" {
						err := qbClient.SetShareLimits([]string{torrent.Hash}, 2.0, torrent.MaxSeedingTime)
						if err != nil {
							log.Printf("set share limit for %s error: %v", torrent.Hash, err)
							continue
						}
						log.Printf("set share limit to 2.0 for %s", torrent.Hash)
					}
				}
				continue
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

type taskManager struct {
	qbClient QbClient
	tasks    []func(qbClient QbClient)
}

// Start starts the task manager
func (t *taskManager) Start() {
	log.Println("task manager started")
	var wg sync.WaitGroup
	for _, task := range t.tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task(t.qbClient)
		}()
	}
	wg.Wait()
}
