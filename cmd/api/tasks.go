package main

import (
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Task struct {
	ID     int
	URL    string
	Result *DataInfo
	Err    error
	State  URLState
	Done   bool
	Stop   bool
}

type TaskQueueInterface interface {
	AddTask(urlInfo *URLInfo) (*Task, error)
	StopTask(id int)
}

type TaskQueue struct {
	tasks       map[int]*Task
	workerCount int
	urlManager  *URLManager
	logger      *logrus.Logger
	mu          sync.Mutex
}

func NewTaskQueue(workerCount int, urlManager *URLManager, logger *logrus.Logger) *TaskQueue {
	tq := &TaskQueue{
		tasks:       make(map[int]*Task),
		workerCount: workerCount,
		urlManager:  urlManager,
		logger:      logger,
	}

	for i := 0; i < workerCount; i++ {
		go tq.worker()
	}

	return tq
}

func (tq *TaskQueue) worker() {
	for {
		tq.mu.Lock()
		for _, task := range tq.tasks {
			if task.State == Pending {
				task.State = Processing
				go tq.processTask(task)
			}
		}
		tq.mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func (tq *TaskQueue) processTask(task *Task) {
	tq.logger.Infof("Processing task ID: %d, URL: %s", task.ID, task.URL)
	tq.urlManager.UpdateURLState(task.ID, Processing)

	data, err := processURL(task.URL, task, tq.logger)

	if task.Stop {
		tq.logger.Infof("Task ID: %d stopped", task.ID)
		task.State = Stopped
		tq.urlManager.UpdateURLState(task.ID, Stopped)
	} else {
		task.Result = data
		task.Err = err

		if err == nil {
			task.State = Completed
			tq.urlManager.UpdateProcessedData(task.ID, data)
		} else {
			task.State = Failed
			tq.urlManager.UpdateURLState(task.ID, Failed)
		}
	}

	task.Done = true
}

func (tq *TaskQueue) StopTask(id int) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	if task, exists := tq.tasks[id]; exists {
		if !task.Stop && !task.Done {
			task.Stop = true
			task.State = Stopped
			tq.logger.Infof("Task ID: %d stop signal sent", task.ID)
		} else {
			tq.logger.Warnf("Task ID: %d already stopped or completed", task.ID)
		}
	}
}

func (tq *TaskQueue) AddTask(urlInfo *URLInfo) (*Task, error) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if task, exists := tq.tasks[urlInfo.ID]; exists && task.State != Completed && task.State != Stopped {
		return task, errors.New("task already in progress")
	}

	task := &Task{
		ID:    urlInfo.ID,
		URL:   urlInfo.URL,
		State: Pending,
		Done:  false,
		Stop:  false,
	}

	tq.tasks[task.ID] = task
	return task, nil
}

func processURL(url string, task *Task, logger *logrus.Logger) (*DataInfo, error) {
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)

		if task.Stop {
			logger.Infof("Task ID: %d processing stopped", task.ID)
			return nil, errors.New("task stopped")
		}
	}

	data := &DataInfo{
		HTMLVersion:       "HTML5",
		PageTitle:         "Example Page",
		HeadingTagsCount:  map[string]int{"h1": 1, "h2": 2, "h3": 3},
		InternalLinks:     5,
		ExternalLinks:     3,
		InaccessibleLinks: 1,
		HasLoginForm:      true,
	}

	return data, nil
}
