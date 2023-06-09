package client

import (
	"log"
	"sync"
	"time"

	"qb-monitor/model"

	"github.com/antonmedv/expr"
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
		checkTorrent,
	}
	return taskManager
}

var checkTorrent = func(conf model.Config, qbClient QbClient) error {
	log.Printf("task checkTorrent started")
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
		for _, rule := range conf.Rules {
			res, err := expr.Run(rule.Evaluator, torrent)
			if err != nil {
				log.Printf("eval rule error: %v", err)
				continue
			}
			if res == false {
				continue
			}
			if err := execTorrentAction(qbClient, torrent, rule.Action); err != nil {
				log.Printf("set share limit for %s error: %v", torrent.Hash, err)
			}
		}
	}
	log.Printf("task checkTorrent finished")
	return nil
}

func execTorrentAction(qbClient QbClient, torrent model.Torrent, action model.TorrentAction) error {
	if action.MaxRatio != nil {
		if err := qbClient.SetShareLimits([]string{torrent.Hash}, *action.MaxRatio, torrent.MaxSeedingTime); err != nil {
			return err
		}
		log.Printf("set share limit for %s (%s) to %f", torrent.Hash, torrent.Name, *action.MaxRatio)
	}
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
		go func(f func(config model.Config, qbClient QbClient) error) {
			defer wg.Done()
			for {
				err := f(t.config, t.qbClient)
				if err != nil {
					log.Printf("Error running task: %v", err)
				}
				time.Sleep(10 * time.Minute)
			}
		}(task)
	}
	wg.Wait()
}
