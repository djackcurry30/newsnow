package sources

import (
	"time"

	"newsnow-go/internal/fetcher"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
)

type V2EXRes struct {
	Version      string       `json:"version"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	HomePageURL  string       `json:"home_page_url"`
	FeedURL      string       `json:"feed_url"`
	Icon         string       `json:"icon"`
	Favicon      string       `json:"favicon"`
	Items        []V2EXItem   `json:"items"`
}

type V2EXItem struct {
	URL           string `json:"url"`
	DateModified  string `json:"date_modified,omitempty"`
	DatePublished string `json:"date_published"`
	ContentHTML   string `json:"content_html"`
	Title         string `json:"title"`
	ID            string `json:"id"`
}

type V2EXFetcher struct {
	fetcher.BaseFetcher
}

func NewV2EXFetcher() *V2EXFetcher {
	return &V2EXFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "v2ex"},
	}
}

func (f *V2EXFetcher) Fetch() ([]types.NewsItem, error) {
	categories := []string{"create", "ideas", "programmer", "share"}
	results := make([]types.NewsItem, 0)
	resultChan := make(chan []types.NewsItem, len(categories))
	errorChan := make(chan error, len(categories))

	for _, cat := range categories {
		go func(category string) {
			url := "https://www.v2ex.com/feed/" + category + ".json"
			var res V2EXRes
			if err := utils.FetchJSON(url, &res); err != nil {
				errorChan <- err
				return
			}

			items := make([]types.NewsItem, 0, len(res.Items))
			for _, item := range res.Items {
				date := item.DateModified
				if date == "" {
					date = item.DatePublished
				}
				pubDate := ""
				if date != "" {
					pubDate = utils.ParseDate(date).Format(time.RFC3339)
				}
				items = append(items, types.NewsItem{
					ID:      item.ID,
					Title:   item.Title,
					URL:     item.URL,
					PubDate: pubDate,
					Extra: &types.NewsExtra{
						Date: pubDate,
					},
				})
			}
			resultChan <- items
		}(cat)
	}

	for i := 0; i < len(categories); i++ {
		select {
		case err := <-errorChan:
			return nil, err
		case items := <-resultChan:
			results = append(results, items...)
		}
	}

	fetcher.SortByDate(results)
	return results, nil
}

type V2EXHolder struct {
	Share *V2EXFetcher
}

func NewV2EXHolder() *V2EXHolder {
	share := NewV2EXFetcher()
	return &V2EXHolder{
		Share: share,
	}
}

func (h *V2EXHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"v2ex":       h.Share,
		"v2ex-share": h.Share,
	}
}
