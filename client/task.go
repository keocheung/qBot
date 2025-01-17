package client

import (
	"fmt"
	"sync"
	"time"

	"qbot/model"
	"qbot/util/logger"

	"github.com/antonmedv/expr"
)

// TaskManager is the task manager of qbot
type TaskManager interface {
	Start() // Start the task manager
}

// NewTaskManager creates a new task manager
func NewTaskManager(config *model.Config, qbClient QbClient) TaskManager {
	taskManager := &taskManager{
		config:   config,
		qbClient: qbClient,
		tasks:    []func(config *model.Config, qbClient QbClient) error{},
	}
	taskManager.tasks = []func(config *model.Config, qbClient QbClient) error{
		checkTorrent,
	}
	return taskManager
}

var checkTorrent = func(conf *model.Config, qbClient QbClient) error {
	logger.Infof("task checkTorrent started")
	if conf == nil {
		return fmt.Errorf("config is nil")
	}
	torrents, err := qbClient.GetTorrents(model.Options{
		// TODO: set default values
		Limit:   conf.GetTorrents.Limit,
		Sort:    conf.GetTorrents.Sort,
		Reverse: conf.GetTorrents.Reverse,
	})
	if err != nil {
		logger.Errorf("get torrents error: %v", err)
		return err
	}
	for _, torrent := range torrents {
		for _, rule := range conf.Rules {
			res, err := expr.Run(rule.Evaluator, torrent)
			if err != nil {
				logger.Errorf("eval rule error: %v", err)
				continue
			}
			if res == false {
				continue
			}
			if err := execTorrentAction(qbClient, torrent, rule.Action); err != nil {
				logger.Errorf("set share limit for %s error: %v", torrent.Hash, err)
			}
		}
	}
	logger.Infof("task checkTorrent finished")
	return nil
}

func execTorrentAction(qbClient QbClient, torrent model.Torrent, action model.TorrentAction) error {
	if action.State != "" {
		switch action.State {
		case model.TorrentStateStop:
			if err := qbClient.StopTorrents([]string{torrent.Hash}); err != nil {
				return err
			}
			logger.Infof("stopped torrent %s", torrent.Hash)
		}
	}
	if action.MaxRatio != nil {
		if err := qbClient.SetShareLimits([]string{torrent.Hash}, *action.MaxRatio, torrent.MaxSeedingTime); err != nil {
			return err
		}
		logger.Infof("set share limit to %f for %s (%s)", torrent.Hash, torrent.Name, *action.MaxRatio)
	}
	return nil
}

type taskManager struct {
	config   *model.Config
	qbClient QbClient
	tasks    []func(config *model.Config, qbClient QbClient) error
}

// Start starts the task manager
func (t *taskManager) Start() {
	logger.Infof("task manager started")
	var wg sync.WaitGroup
	for _, task := range t.tasks {
		wg.Add(1)
		go func(f func(config *model.Config, qbClient QbClient) error) {
			defer wg.Done()
			for {
				logger.Debugf("running task, config: %+v", t.config)
				err := f(t.config, t.qbClient)
				if err != nil {
					logger.Errorf("Error running task: %v", err)
				}
				time.Sleep(10 * time.Minute)
			}
		}(task)
	}
	wg.Wait()
}
