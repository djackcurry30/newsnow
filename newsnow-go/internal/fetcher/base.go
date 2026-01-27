package fetcher

import (
	"net/url"
	"sort"
	"strconv"
	"time"

	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
)

type SourceFetcher interface {
	Fetch() ([]types.NewsItem, error)
	GetID() string
}

type BaseFetcher struct {
	ID string
}

func (b *BaseFetcher) GetID() string {
	return b.ID
}

type RSSFetcher struct {
	BaseFetcher
	URL        string
	HiddenDate bool
}

func (f *RSSFetcher) Fetch() ([]types.NewsItem, error) {
	channel, err := utils.ParseRSS(f.URL)
	if err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(channel.Items))
	for _, item := range channel.Items {
		pubDate := ""
		if !f.HiddenDate && item.PubDate != "" {
			pubDate = utils.ParseDate(item.PubDate).Format(time.RFC3339)
		}
		items = append(items, types.NewsItem{
			ID:      item.Guid,
			Title:   item.Title,
			URL:     item.Link,
			PubDate: pubDate,
		})
	}

	return items, nil
}

type RSSHubFetcher struct {
	BaseFetcher
	Route       string
	Options     RSSHubOptions
	HiddenDate  bool
}

type RSSHubOptions struct {
	Sorted *bool
	Limit  *int
}

const RSSHubBase = "https://rsshub.rssforever.com"

func (f *RSSHubFetcher) Fetch() ([]types.NewsItem, error) {
	baseURL := RSSHubBase
	u, err := url.Parse(f.Route)
	if err != nil {
		return nil, err
	}

	fullURL := baseURL + f.Route
	q := u.Query()

	sorted := true
	if f.Options.Sorted != nil {
		sorted = *f.Options.Sorted
	}
	q.Set("format", "json")

	limit := 20
	if f.Options.Limit != nil {
		limit = *f.Options.Limit
	}
	q.Set("limit", strconv.Itoa(limit))

	if sorted {
		q.Set("sorted", "true")
	}

	fullURL += "?" + q.Encode()

	resp, err := utils.FetchRSSHub(fullURL)
	if err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		pubDate := ""
		if !f.HiddenDate && item.DatePublished != "" {
			pubDate = utils.ParseDate(item.DatePublished).Format(time.RFC3339)
		}
		items = append(items, types.NewsItem{
			ID:      item.Id,
			Title:   item.Title,
			URL:     item.Url,
			PubDate: pubDate,
		})
	}

	return items, nil
}

func DefineSource(fetchers interface{}) interface{} {
	return fetchers
}

func DefineRSSSource(url string, hiddenDate bool) SourceFetcher {
	return &RSSFetcher{
		BaseFetcher: BaseFetcher{},
		URL:         url,
		HiddenDate:  hiddenDate,
	}
}

func DefineRSSHubSource(route string, options RSSHubOptions, hiddenDate bool) SourceFetcher {
	return &RSSHubFetcher{
		BaseFetcher: BaseFetcher{},
		Route:       route,
		Options:     options,
		HiddenDate:  hiddenDate,
	}
}

func ProxySource(proxyURL string, fetcher SourceFetcher, isCloudflare bool) SourceFetcher {
	if isCloudflare {
		return &ProxyFetcher{
			BaseFetcher: BaseFetcher{ID: fetcher.GetID()},
			ProxyURL:    proxyURL,
		}
	}
	return fetcher
}

type ProxyFetcher struct {
	BaseFetcher
	ProxyURL string
}

func (f *ProxyFetcher) Fetch() ([]types.NewsItem, error) {
	var result struct {
		Items []types.NewsItem `json:"items"`
	}
	err := utils.FetchJSON(f.ProxyURL, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

func SortByDate(items []types.NewsItem) {
	sort.Slice(items, func(i, j int) bool {
		dateI := utils.ParseDate(items[i].PubDate)
		dateJ := utils.ParseDate(items[j].PubDate)
		return dateI.After(dateJ)
	})
}
