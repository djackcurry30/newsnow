package utils

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

type RSSChannel struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	Items       []RSSItem `xml:"item"` // 修复标签，从 "channel>item" 改为 "item"
}

type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
	Content     string `xml:"content:encoded"`
}

type RSSFeed struct {
	Channel RSSChannel `xml:"channel"`
}

func ParseRSS(url string) (*RSSChannel, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, err
	}

	return &feed.Channel, nil
}

func ParseRSSFromBytes(data []byte) (*RSSChannel, error) {
	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}
	return &feed.Channel, nil
}

type RSSHubItem struct {
	Id            string `json:"id"`
	Url           string `json:"url"`
	Title         string `json:"title"`
	ContentHtml   string `json:"content_html"`
	DatePublished string `json:"date_published"`
}

type RSSHubResponse struct {
	Title       string       `json:"title"`
	HomePageURL string       `json:"home_page_url"`
	Description string       `json:"description"`
	Items       []RSSHubItem `json:"items"`
}

func ParseRSSHub(data []byte) (*RSSHubResponse, error) {
	var resp RSSHubResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func FetchRSSHub(url string) (*RSSHubResponse, error) {
	data, err := FetchBytes(url)
	if err != nil {
		return nil, err
	}
	return ParseRSSHub(data)
}
