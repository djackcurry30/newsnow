package sources

import (
	"strings"

	"newsnow-go/internal/fetcher"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
)

type CoolapkItem struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	EditorTitle string `json:"editor_title"`
	URL         string `json:"url"`
	EntityType  string `json:"entityType"`
	PubDate     string `json:"pubDate"`
	Dateline    int    `json:"dateline"`
	TargetRow   struct {
		SubTitle string `json:"subTitle"`
	} `json:"targetRow"`
}

type CoolapkRes struct {
	Data []CoolapkItem `json:"data"`
}

type CoolapkFetcher struct {
	fetcher.BaseFetcher
}

func NewCoolapkFetcher() *CoolapkFetcher {
	return &CoolapkFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "coolapk"},
	}
}

func (f *CoolapkFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://api.coolapk.com/v6/page/dataList?url=%2Ffeed%2FstatList%3FcacheExpires%3D300%26statType%3Dday%26sortField%3Ddetailnum%26title%3D%E4%BB%8A%E6%97%A5%E7%83%AD%E9%97%A8&title=%E4%BB%8A%E6%97%A5%E7%83%AD%E9%97%A8&subTitle=&page=1"

	headers := map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"X-App-Id":         "com.coolapk.market",
		"X-Sdk-Int":        "29",
		"X-Sdk-Locale":     "zh-CN",
		"X-App-Version":    "11.0",
		"X-Api-Version":    "11",
		"X-App-Code":       "2101202",
		"User-Agent":       "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36",
	}

	var res CoolapkRes
	if err := utils.FetchJSON(url, &res, utils.FetchOptions{Headers: headers}); err != nil {
		return nil, err
	}

	if len(res.Data) == 0 {
		return []types.NewsItem{}, nil
	}

	items := make([]types.NewsItem, 0, len(res.Data))
	for _, item := range res.Data {
		if item.ID == "" {
			continue
		}

		title := item.EditorTitle
		if title == "" {
			lines := strings.Split(item.Message, "\n")
			if len(lines) > 0 {
				title = lines[0]
			} else {
				title = item.Message
			}
		}

		items = append(items, types.NewsItem{
			ID:   item.ID,
			Title: title,
			URL:  "https://www.coolapk.com" + item.URL,
			Extra: &types.NewsExtra{
				Info: item.TargetRow.SubTitle,
			},
		})
	}

	return items, nil
}
