Inspired by [pygooglenews](https://github.com/kotartemiy/pygooglenews)

# About

A golang based wrapper of the Google News RSS feed.

Top stories, topic related news feeds, geolocation news feed, and an extensive full text search feed.

This work is more of a collection of all things I could find out about how Google News functions.

# Installation

`go get github.com/spl0i7/gogooglenews`

# Quickstart

```go

googleNews, err := NewGoogleNews(GoogleNewsOpt{
    Lang:    "en",
    Country: "IN",
})

```

## Top Stories

```go
news, err = googleNews.TopNews()
```

## Stories by Topic

```go
news, err = googleNews.TopicHeadlines("BUSINESS")

```
## Geolocation Specific Stories

```go
news, err = googleNews.GeoHeadlines("bangalore")

```

## Stories by a Query Search

```go

to := time.Now()
from := to.AddDate(0, 0, -10)

news, err = googleNews.Search("Bitcoin", &from, &to)
```

# Proxy

If you make frequent calls to the rss service, its very much possible that Google might blacklist your IP.

To bypass this, ideally you should be making requests behind a proxy. To use a proxy pass the http client in options as following.

```go
proxyUrl, err := url.Parse("http://proxyIp:proxyPort")
myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

googleNews, err := NewGoogleNews(GoogleNewsOpt{
    Lang:    "en",
    Country: "IN",
    HttpClient: myClient
})

```
