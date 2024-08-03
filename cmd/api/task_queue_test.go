package main

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAddTask(t *testing.T) {
	logger := logrus.New()
	mockURLManager := &MockURLManager{
		AddURLFunc: func(url string) *URLInfo {
			return &URLInfo{ID: 1, URL: url, State: Pending, UploadedAt: time.Now()}
		},
		GetURLStateFunc: func(id int) URLState {
			return Pending
		},
		UpdateURLStateFunc: func(id int, state URLState) {},
	}

	tq := NewTaskQueue(2, mockURLManager, logger)

	urlInfo := mockURLManager.AddURL("http://example.com")
	task, err := tq.AddTask(urlInfo)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, urlInfo.ID, task.ID)
}

func TestStopTask(t *testing.T) {
	logger := logrus.New()
	mockURLManager := &MockURLManager{
		AddURLFunc: func(url string) *URLInfo {
			return &URLInfo{ID: 1, URL: url, State: Pending, UploadedAt: time.Now()}
		},
		GetURLStateFunc: func(id int) URLState {
			return Pending
		},
		UpdateURLStateFunc: func(id int, state URLState) {},
	}

	tq := NewTaskQueue(2, mockURLManager, logger)

	urlInfo := mockURLManager.AddURL("http://example.com")
	task, err := tq.AddTask(urlInfo)
	assert.NoError(t, err)

	stoppedTask, err := tq.StopTask(task.ID)
	assert.NoError(t, err)
	assert.NotNil(t, stoppedTask)
	assert.True(t, stoppedTask.Stop)
}
