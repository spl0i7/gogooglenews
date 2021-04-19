package gogooglenews

import (
	"net/http"
	"time"
)

type GoogleNews interface {
	TopNews() ([]News, error)
	TopicHeadlines(topic string) ([]News, error)
	GeoHeadlines(geo string) ([]News, error)
	Search(query string, from, to *time.Time) ([]News, error)
}

type GoogleNewsOpt struct {
	Lang       string
	Country    string
	HttpClient *http.Client
}

type News struct {
	Title       string
	Link        string
	Time        time.Time
	Description string
	Source      string
	SourceUrl   string
}
