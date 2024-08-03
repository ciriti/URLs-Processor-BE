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
	Done   bool
	Stop   bool
}

type TaskQueueInterface interface {
	AddTask(urlInfo *URLInfo) (*Task, error)
	StopTask(id int) (*Task, error)
}

type TaskQueue struct {
	tasks       map[int]*Task
	workerCount int
	urlManager  *URLManager
	logger      *logrus.Logger
	mu          sync.Mutex
	sem         chan struct{}
}

func NewTaskQueue(workerCount int, urlManager *URLManager, logger *logrus.Logger) *TaskQueue {
	tq := &TaskQueue{
		tasks:       make(map[int]*Task),
		workerCount: workerCount,
		urlManager:  urlManager,
		logger:      logger,
		sem:         make(chan struct{}, workerCount),
	}

	for i := 0; i < workerCount; i++ {
		go tq.worker()
	}

	return tq
}

func (tq *TaskQueue) worker() {
	for {
		tq.sem <- struct{}{}
		tq.mu.Lock()
		var taskToProcess *Task
		for _, task := range tq.tasks {
			if tq.urlManager.GetURLState(task.ID) == Pending {
				tq.urlManager.UpdateURLState(task.ID, Processing)
				taskToProcess = task
				break
			}
		}
		tq.mu.Unlock()

		if taskToProcess != nil {
			go func(task *Task) {
				tq.processTask(task)
				<-tq.sem
			}(taskToProcess)
		} else {
			<-tq.sem
			time.Sleep(1 * time.Second)
		}
	}
}

func (tq *TaskQueue) processTask(task *Task) {
	tq.logger.Infof("Processing task ID: %d, URL: %s", task.ID, task.URL)
	tq.urlManager.UpdateURLState(task.ID, Processing)

	data, err := processURL(task.URL, task, tq.logger)

	tq.mu.Lock()
	defer tq.mu.Unlock()

	if task.Stop {
		tq.logger.Infof("Task ID: %d processing stopped", task.ID)
		tq.urlManager.UpdateURLState(task.ID, Stopped)
	} else {
		task.Result = data
		task.Err = err

		if err == nil {
			tq.urlManager.UpdateProcessedData(task.ID, data)
		} else {
			tq.urlManager.UpdateURLState(task.ID, Failed)
		}
	}

	task.Done = true
}

func (tq *TaskQueue) AddTask(urlInfo *URLInfo) (*Task, error) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	task, exists := tq.tasks[urlInfo.ID]
	if exists {

		state := tq.urlManager.GetURLState(urlInfo.ID)
		tq.logger.Infof(" GetURLState: %s", state)
		if state == Stopped || state == Completed {
			task.Stop = false
			task.Done = false
			task.Result = nil
			task.Err = nil
			tq.urlManager.UpdateURLState(urlInfo.ID, Pending)
			tq.logger.Infof("Resetting task ID: %d", urlInfo.ID)
		}
	} else {

		task = &Task{
			ID:   urlInfo.ID,
			URL:  urlInfo.URL,
			Done: false,
			Stop: false,
		}
		tq.tasks[task.ID] = task
	}

	return task, nil
}

func (tq *TaskQueue) StopTask(id int) (*Task, error) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	if task, exists := tq.tasks[id]; exists {
		tq.logger.Infof("StopTask - task.Done: %v - task.Stop: %v", task.Done, task.Stop)
		if !task.Stop {
			task.Stop = true
			tq.urlManager.UpdateURLState(id, Stopped)
			tq.logger.Infof("StopTask - Task ID: %d stop signal sent", task.ID)
		} else {
			tq.logger.Infof("StopTask - Task ID: %d already stopped or completed", task.ID)
		}
		return task, nil
	} else {
		tq.logger.Warnf("Task ID: %d not found", id)
		return nil, errors.New("task not found")
	}
}

func processURL(url string, task *Task, logger *logrus.Logger) (*DataInfo, error) {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)

		if task.Stop {
			logger.Infof(" processURL - Task ID: %d processing stopped", task.ID)
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
