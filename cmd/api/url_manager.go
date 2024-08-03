package main

import (
	"sync"
	"time"
)

type URLState string

const (
	Pending    URLState = "pending"
	Processing URLState = "processing"
	Stopped    URLState = "stopped"
	Completed  URLState = "completed"
	Failed     URLState = "failed"
)

type URLInfo struct {
	ID            int       `json:"id"`
	URL           string    `json:"url"`
	State         URLState  `json:"state"`
	ProcessedData *DataInfo `json:"processed_data,omitempty"`
	UploadedAt    time.Time `json:"uploaded_at"`
}

type DataInfo struct {
	HTMLVersion        string         `json:"html_version"`
	PageTitle          string         `json:"page_title"`
	HeadingTagsCount   map[string]int `json:"heading_tags_count"`
	InternalLinks      int            `json:"internal_links"`
	ExternalLinks      int            `json:"external_links"`
	InaccessibleLinks  int            `json:"inaccessible_links"`
	HasLoginForm       bool           `json:"has_login_form"`
	ProcessingFinished time.Time      `json:"processing_finished"`
}

type URLManagerInterface interface {
	AddURL(url string) *URLInfo
	UpdateURLState(id int, state URLState)
	UpdateProcessedData(id int, data *DataInfo)
	GetURLInfo(id int) *URLInfo
	GetAllURLs() []*URLInfo
	nextID() int
	GetURLState(id int) URLState
}

type URLManager struct {
	mu           sync.RWMutex
	urls         map[int]*URLInfo
	idCounter    int
	counterMutex sync.Mutex
}

func NewURLManager() *URLManager {
	return &URLManager{
		urls: make(map[int]*URLInfo),
	}
}

func (manager *URLManager) nextID() int {
	manager.counterMutex.Lock()
	defer manager.counterMutex.Unlock()
	manager.idCounter++
	return manager.idCounter
}

func (manager *URLManager) AddURL(url string) *URLInfo {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	id := manager.nextID()

	urlInfo := &URLInfo{
		ID:         id,
		URL:        url,
		State:      Pending,
		UploadedAt: time.Now(),
	}

	manager.urls[id] = urlInfo
	return urlInfo
}

func (manager *URLManager) UpdateURLState(id int, state URLState) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if urlInfo, exists := manager.urls[id]; exists {
		urlInfo.State = state
	}
}

func (manager *URLManager) UpdateProcessedData(id int, data *DataInfo) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if urlInfo, exists := manager.urls[id]; exists {
		urlInfo.State = Completed
		urlInfo.ProcessedData = data
		urlInfo.ProcessedData.ProcessingFinished = time.Now()
	}
}

func (manager *URLManager) GetURLInfo(id int) *URLInfo {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	return manager.urls[id]
}

func (manager *URLManager) GetAllURLs() []*URLInfo {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	allURLs := make([]*URLInfo, 0, len(manager.urls))
	for _, urlInfo := range manager.urls {
		allURLs = append(allURLs, urlInfo)
	}
	return allURLs
}

func (manager *URLManager) GetURLState(id int) URLState {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	if urlInfo, exists := manager.urls[id]; exists {
		return urlInfo.State
	}
	return ""
}
