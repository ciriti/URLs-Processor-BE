package main

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

type Task struct {
	ID      int
	URL     string
	Result  *DataInfo
	Err     error
	Done    chan bool
	Stopped bool
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
	// Update URL state to Processing
	tq.urlManager.UpdateURLState(task.ID, Processing)

	// Simulate processing (replace with actual processing logic)
	data, err := processURL(task.URL, &task, tq.logger)

	if task.Stopped {
		tq.urlManager.UpdateURLState(task.ID, Stopped)
	} else {
		task.Result = data
		task.Err = err

		// Update the URLManager with the result
		if err == nil {
			tq.urlManager.UpdateProcessedData(task.ID, data)
		} else {
			tq.urlManager.UpdateURLState(task.ID, Stopped)
		}
	}

	// Signal that the task is done
	task.Done <- true
}

func (tq *TaskQueue) StopTask(id int) {
	tq.urlManager.UpdateURLState(id, Stopped)
	for task := range tq.tasks {
		if task.ID == id {
			task.Stopped = true
			break
		}
	}
}

func (tq *TaskQueue) AddTask(urlInfo *URLInfo) Task {
	task := Task{
		ID:   urlInfo.ID,
		URL:  urlInfo.URL,
		Done: make(chan bool),
	}
	tq.tasks <- task
	return task
}

func processURL(url string, task *Task, logger *logrus.Logger) (*DataInfo, error) {

	// Simulate fetching and processing the URL with periodic checks for stop signal
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second) // Simulate delay

		if task.Stopped {
			return nil, errors.New("task stopped")
		}
	}

	// Replace this with actual URL processing logic
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
