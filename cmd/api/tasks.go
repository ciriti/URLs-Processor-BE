package main

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

type TaskState string

type Task struct {
	ID     int
	URL    string
	Result *DataInfo
	Err    error
	State  URLState
	Done   chan struct{}
}

type TaskQueueInterface interface {
	AddTask(urlInfo *URLInfo) Task
	StopTask(id int)
}

type TaskQueue struct {
	tasks       chan Task
	workerCount int
	urlManager  *URLManager
	logger      *logrus.Logger
}

func NewTaskQueue(workerCount int, urlManager *URLManager, logger *logrus.Logger) *TaskQueue {
	tq := &TaskQueue{
		tasks:       make(chan Task, 100), // buffer size of 100
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
	for task := range tq.tasks {
		tq.processTask(task)
	}
}

func (tq *TaskQueue) processTask(task Task) {
	tq.logger.Infof("Processing task ID: %d, URL: %s", task.ID, task.URL)
	task.State = Processing
	tq.urlManager.UpdateURLState(task.ID, Processing)

	data, err := processURL(task.URL, &task, tq.logger)

	if task.State == Stopped {
		tq.urlManager.UpdateURLState(task.ID, Stopped)
	} else {
		task.Result = data
		task.Err = err

		if err == nil {
			task.State = Completed
			tq.urlManager.UpdateProcessedData(task.ID, data)
		} else {
			task.State = Stopped
			tq.urlManager.UpdateURLState(task.ID, Stopped)
		}
	}

	close(task.Done)
}

func (tq *TaskQueue) StopTask(id int) {
	tq.urlManager.UpdateURLState(id, Stopped)
}

func (tq *TaskQueue) AddTask(urlInfo *URLInfo) Task {
	task := Task{
		ID:    urlInfo.ID,
		URL:   urlInfo.URL,
		State: Pending,
		Done:  make(chan struct{}),
	}
	tq.tasks <- task
	return task
}

func processURL(url string, task *Task, logger *logrus.Logger) (*DataInfo, error) {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)

		if task.State == Stopped {
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
