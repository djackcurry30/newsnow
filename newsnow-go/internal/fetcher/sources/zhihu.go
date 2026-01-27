package sources

import (
	"newsnow-go/internal/fetcher"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
	"regexp"
)

type ZhihuRes struct {
	Data []ZhihuItem `json:"data"`
}

type ZhihuItem struct {
	Type            string          `json:"type"`
	StyleType       string          `json:"style_type"`
	FeedSpecific    ZhihuFeedSpec   `json:"feed_specific"`
	Target          ZhihuTarget     `json:"target"`
}

type ZhihuFeedSpec struct {
	AnswerCount int `json:"answer_count"`
}

type ZhihuTarget struct {
	TitleArea   ZhihuTextArea   `json:"title_area"`
	ExcerptArea ZhihuTextArea   `json:"excerpt_area"`
	ImageArea   ZhihuImageArea  `json:"image_area"`
	MetricsArea ZhihuTextArea   `json:"metrics_area"`
	LabelArea   ZhihuLabelArea  `json:"label_area"`
	Link        ZhihuLink       `json:"link"`
}

type ZhihuTextArea struct {
	Text string `json:"text"`
}

type ZhihuImageArea struct {
	URL string `json:"url"`
}

type ZhihuLabelArea struct {
	Type  string `json:"type"`
	Trend int    `json:"trend"`
}

type ZhihuLink struct {
	URL string `json:"url"`
}

type ZhihuFetcher struct {
	fetcher.BaseFetcher
}

func NewZhihuFetcher() *ZhihuFetcher {
	return &ZhihuFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "zhihu"},
	}
}

func (f *ZhihuFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://www.zhihu.com/api/v3/feed/topstory/hot-list-web?limit=20&desktop=true"
	var res ZhihuRes
	headers := map[string]string{
		"Accept":             "application/json",
		"Referer":            "https://www.zhihu.com/",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}
	
	if err := utils.FetchJSON(url, &res, utils.FetchOptions{Headers: headers}); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data))
	for _, item := range res.Data {
		re := regexp.MustCompile(`(\d+)$`)
		id := re.FindStringSubmatch(item.Target.Link.URL)
		idStr := ""
		if len(id) > 1 {
			idStr = id[1]
		} else {
			idStr = item.Target.Link.URL
		}
		
		items = append(items, types.NewsItem{
			ID:    idStr,
			Title: item.Target.TitleArea.Text,
			URL:   item.Target.Link.URL,
			Extra: &types.NewsExtra{
				Info:  item.Target.MetricsArea.Text,
				Hover: item.Target.ExcerptArea.Text,
			},
		})
	}

	return items, nil
}
