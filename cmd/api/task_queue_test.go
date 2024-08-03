package main

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MockPageAnalyzer struct {
	AnalyzePageFunc func(url string, task *Task) (*DataInfo, error)
}

func (m *MockPageAnalyzer) AnalyzePage(url string, task *Task) (*DataInfo, error) {
	return m.AnalyzePageFunc(url, task)
}

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

	mockPageAnalyzer := &MockPageAnalyzer{
		AnalyzePageFunc: func(url string, task *Task) (*DataInfo, error) {
			return &DataInfo{
				HTMLVersion:       "HTML5",
				PageTitle:         "Mock Page",
				HeadingTagsCount:  map[string]int{"h1": 1, "h2": 2},
				InternalLinks:     1,
				ExternalLinks:     1,
				InaccessibleLinks: 0,
				HasLoginForm:      true,
			}, nil
		},
	}

	tq := NewTaskQueue(2, mockURLManager, mockPageAnalyzer, logger)

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

	mockPageAnalyzer := &MockPageAnalyzer{
		AnalyzePageFunc: func(url string, task *Task) (*DataInfo, error) {
			return &DataInfo{
				HTMLVersion:       "HTML5",
				PageTitle:         "Mock Page",
				HeadingTagsCount:  map[string]int{"h1": 1, "h2": 2},
				InternalLinks:     1,
				ExternalLinks:     1,
				InaccessibleLinks: 0,
				HasLoginForm:      true,
			}, nil
		},
	}

	tq := NewTaskQueue(2, mockURLManager, mockPageAnalyzer, logger)

	urlInfo := mockURLManager.AddURL("http://example.com")
	task, err := tq.AddTask(urlInfo)
	assert.NoError(t, err)

	stoppedTask, err := tq.StopTask(task.ID)
	assert.NoError(t, err)
	assert.NotNil(t, stoppedTask)
	assert.True(t, stoppedTask.Stop)
}
