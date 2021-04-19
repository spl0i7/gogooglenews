package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	strip "github.com/grokify/html-strip-tags-go"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Media   string   `xml:"media,attr"`
	Channel struct {
		Text          string `xml:",chardata"`
		Generator     string `xml:"generator"`
		Title         string `xml:"title"`
		Link          string `xml:"link"`
		Language      string `xml:"language"`
		WebMaster     string `xml:"webMaster"`
		Copyright     string `xml:"copyright"`
		LastBuildDate string `xml:"lastBuildDate"`
		Description   string `xml:"description"`
		Item          []struct {
			Text  string `xml:",chardata"`
			Title string `xml:"title"`
			Link  string `xml:"link"`
			Guid  struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
			PubDate     string `xml:"pubDate"`
			Description string `xml:"description"`
			Source      struct {
				Text string `xml:",chardata"`
				URL  string `xml:"url,attr"`
			} `xml:"source"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (r rss) ToNews() []News {

	var result []News

	for _, i := range r.Channel.Item {

		desc := strip.StripTags(i.Description)
		desc = strings.ReplaceAll(desc, "&nbsp;", " ")

		parsed, _ := time.Parse(time.RFC1123, i.PubDate)
		result = append(result, News{
			Title:       i.Title,
			Link:        i.Link,
			Time:        parsed,
			Description: desc,
			Source:      i.Source.Text,
			SourceUrl:   i.Source.URL,
		})
	}

	return result

}

const rssNewsEndpoint = "https://news.google.com/rss"

type rssGoogleNews struct {
	httpClient *http.Client
	lang       string
	country    string
}

func (r rssGoogleNews) doRequest(requestUrl *url.URL) ([]News, error) {
	resp, err := r.httpClient.Get(requestUrl.String())
	if err != nil {
		return nil, err
	}

	var response rss
	if err := xml.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.ToNews(), nil
}
func (r rssGoogleNews) newRequestUrl() *url.URL {
	requestUrl, _ := url.Parse(rssNewsEndpoint)
	query := requestUrl.Query()
	query.Set("hl", fmt.Sprintf("%s-%s", r.lang, r.country))
	query.Set("gl", r.country)
	query.Set("ceid", fmt.Sprintf("%s:%s", r.country, r.lang))
	requestUrl.RawQuery = query.Encode()
	return requestUrl

}

func (r rssGoogleNews) TopNews() ([]News, error) {
	requestUrl := r.newRequestUrl()
	return r.doRequest(requestUrl)
}

func (r rssGoogleNews) TopicHeadlines(topic string) ([]News, error) {
	requestUrl := r.newRequestUrl()
	requestUrl.Path = path.Join(requestUrl.Path, "headlines", "section", "topic", topic)
	return r.doRequest(requestUrl)

}

func (r rssGoogleNews) GeoHeadlines(geo string) ([]News, error) {
	requestUrl := r.newRequestUrl()
	requestUrl.Path = path.Join(requestUrl.Path, "headlines", "section", "geo", geo)
	return r.doRequest(requestUrl)

}

func (r rssGoogleNews) Search(query string, from, to *time.Time) ([]News, error) {
	requestUrl := r.newRequestUrl()
	requestUrl.Path = path.Join(requestUrl.Path, "search")
	if from != nil {
		query += "+after:" + from.Format("2006-01-02")
	}
	if to != nil {
		query += "+before:" + to.Format("2006-01-02")
	}
	requestUrl.RawQuery = "q=" + query + "&" + requestUrl.Query().Encode()

	return r.doRequest(requestUrl)
}

func NewGoogleNews(opt GoogleNewsOpt) (GoogleNews, error) {
	if len(opt.Country) == 0 {
		return nil, errors.New("country not provided in opt")
	}
	if len(opt.Lang) == 0 {
		return nil, errors.New("language not provided in opt")
	}

	if opt.HttpClient == nil {
		opt.HttpClient = http.DefaultClient
	}

	return &rssGoogleNews{
		httpClient: opt.HttpClient,
		lang:       opt.Lang,
		country:    opt.Country,
	}, nil
}
