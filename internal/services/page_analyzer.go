package services

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type PageAnalyzerInterface interface {
	AnalyzePage(url string, task *Task) (*DataInfo, error)
}

type PageAnalyzer struct {
	client *http.Client
	logger *logrus.Logger
}

func NewPageAnalyzer(client *http.Client, logger *logrus.Logger) *PageAnalyzer {
	return &PageAnalyzer{client: client, logger: logger}
}

func (pa *PageAnalyzer) AnalyzePage(url string, task *Task) (*DataInfo, error) {
	pa.logger.Infof("Starting analysis for URL: %s", url)

	resp, err := pa.client.Get(url)
	if err != nil {
		pa.logger.Errorf("Failed to fetch URL: %s, error: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("failed to fetch URL: " + resp.Status)
		pa.logger.Errorf("URL returned non-OK status: %s, status code: %d", url, resp.StatusCode)
		return nil, err
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		pa.logger.Errorf("Failed to parse HTML for URL: %s, error: %v", url, err)
		return nil, err
	}

	data := &DataInfo{
		HeadingTagsCount: make(map[string]int),
	}

	// Detect HTML version
	if doc.FirstChild != nil && doc.FirstChild.Type == html.DoctypeNode {
		if strings.Contains(doc.FirstChild.Data, "html") {
			data.HTMLVersion = "HTML5"
		} else {
			data.HTMLVersion = "HTML 4.01"
		}
	} else {
		data.HTMLVersion = "HTML 4.01"
	}

	// Traverse the document
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil {
					data.PageTitle = n.FirstChild.Data
				}
			case "h1", "h2", "h3", "h4", "h5", "h6":
				data.HeadingTagsCount[n.Data]++
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						if strings.HasPrefix(attr.Val, "http") {
							data.ExternalLinks++
							if pa.isInaccessible(attr.Val) {
								data.InaccessibleLinks++
							}
						} else {
							data.InternalLinks++
						}
					}
				}
			case "form":
				for _, attr := range n.Attr {
					if attr.Key == "action" && strings.Contains(attr.Val, "login") {
						data.HasLoginForm = true
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	pa.logger.Infof("Completed analysis for URL: %s", url)

	return data, nil
}

func (pa *PageAnalyzer) isInaccessible(url string) bool {
	resp, err := pa.client.Get(url)
	if err != nil {
		return true
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 400 && resp.StatusCode < 600
}
