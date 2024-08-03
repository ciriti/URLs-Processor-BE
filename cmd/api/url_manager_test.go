package main

import (
	"testing"
)

func TestAddURL(t *testing.T) {
	manager := NewURLManager()
	url := "http://example.com"

	urlInfo := manager.AddURL(url)

	if urlInfo.URL != url {
		t.Errorf("expected URL %s, got %s", url, urlInfo.URL)
	}

	if urlInfo.State != Pending {
		t.Errorf("expected state %s, got %s", Pending, urlInfo.State)
	}

	if urlInfo.ID != 1 {
		t.Errorf("expected ID %d, got %d", 1, urlInfo.ID)
	}

	if urlInfo.UploadedAt.IsZero() {
		t.Error("expected UploadedAt to be set")
	}
}

func TestUpdateURLState(t *testing.T) {
	manager := NewURLManager()
	url := "http://example.com"

	urlInfo := manager.AddURL(url)
	manager.UpdateURLState(urlInfo.ID, Processing)

	if urlInfo.State != Processing {
		t.Errorf("expected state %s, got %s", Processing, urlInfo.State)
	}
}

func TestUpdateProcessedData(t *testing.T) {
	manager := NewURLManager()
	url := "http://example.com"
	data := &DataInfo{
		HTMLVersion:       "HTML5",
		PageTitle:         "Example",
		HeadingTagsCount:  map[string]int{"h1": 1, "h2": 2},
		InternalLinks:     5,
		ExternalLinks:     3,
		InaccessibleLinks: 1,
		HasLoginForm:      true,
	}

	urlInfo := manager.AddURL(url)
	manager.UpdateProcessedData(urlInfo.ID, data)

	if urlInfo.State != Completed {
		t.Errorf("expected state %s, got %s", Completed, urlInfo.State)
	}

	if urlInfo.ProcessedData == nil {
		t.Fatal("expected ProcessedData to be set")
	}
}

func TestGetURLInfo(t *testing.T) {
	manager := NewURLManager()
	url := "http://example.com"

	addedURL := manager.AddURL(url)
	retrievedURL := manager.GetURLInfo(addedURL.ID)

	if retrievedURL == nil {
		t.Fatal("expected URLInfo to be retrieved")
	}

	if retrievedURL.URL != url {
		t.Errorf("expected URL %s, got %s", url, retrievedURL.URL)
	}
}

func TestGetAllURLs(t *testing.T) {
	manager := NewURLManager()
	url1 := "http://example1.com"
	url2 := "http://example2.com"

	manager.AddURL(url1)
	manager.AddURL(url2)

	urls := manager.GetAllURLs()
	if len(urls) != 2 {
		t.Errorf("expected 2 URLs, got %d", len(urls))
	}
}

func TestGetURLState(t *testing.T) {
	manager := NewURLManager()
	url := "http://example.com"

	urlInfo := manager.AddURL(url)
	state := manager.GetURLState(urlInfo.ID)

	if state != Pending {
		t.Errorf("expected state %s, got %s", Pending, state)
	}
}
