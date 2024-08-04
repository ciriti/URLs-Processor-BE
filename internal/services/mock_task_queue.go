package services

import "errors"

type MockTaskQueue struct {
	AddTaskFunc  func(urlInfo *URLInfo) (*Task, error)
	StopTaskFunc func(id int) (*Task, error)
}

func (m *MockTaskQueue) AddTask(urlInfo *URLInfo) (*Task, error) {
	if m.AddTaskFunc != nil {
		return m.AddTaskFunc(urlInfo)
	}
	return nil, errors.New("AddTask function not implemented")
}

func (m *MockTaskQueue) StopTask(id int) (*Task, error) {
	if m.StopTaskFunc != nil {
		return m.StopTaskFunc(id)
	}
	return nil, errors.New("StopTask function not implemented")
}
