package main

type MockURLManager struct {
	AddURLFunc              func(url string) *URLInfo
	UpdateURLStateFunc      func(id int, state URLState)
	UpdateProcessedDataFunc func(id int, data *DataInfo)
	GetURLInfoFunc          func(id int) *URLInfo
	GetAllURLsFunc          func() []*URLInfo
	NextIDFunc              func() int
	GetURLStateFunc         func(id int) URLState
}

func (m *MockURLManager) AddURL(url string) *URLInfo {
	if m.AddURLFunc != nil {
		return m.AddURLFunc(url)
	}
	return nil
}

func (m *MockURLManager) UpdateURLState(id int, state URLState) {
	if m.UpdateURLStateFunc != nil {
		m.UpdateURLStateFunc(id, state)
	}
}

func (m *MockURLManager) UpdateProcessedData(id int, data *DataInfo) {
	if m.UpdateProcessedDataFunc != nil {
		m.UpdateProcessedDataFunc(id, data)
	}
}

func (m *MockURLManager) GetURLInfo(id int) *URLInfo {
	if m.GetURLInfoFunc != nil {
		return m.GetURLInfoFunc(id)
	}
	return nil
}

func (m *MockURLManager) GetAllURLs() []*URLInfo {
	if m.GetAllURLsFunc != nil {
		return m.GetAllURLsFunc()
	}
	return nil
}

func (m *MockURLManager) nextID() int {
	if m.NextIDFunc != nil {
		return m.NextIDFunc()
	}
	return 0
}

func (m *MockURLManager) GetURLState(id int) URLState {
	if m.GetURLStateFunc != nil {
		return m.GetURLStateFunc(id)
	}
	return ""
}
