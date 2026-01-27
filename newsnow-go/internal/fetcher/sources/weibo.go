package sources

import (
	"newsnow-go/internal/fetcher"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WeiboSearchRes struct {
	Statuses []WeiboStatus `json:"statuses"`
}

type WeiboStatus struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	TextRaw     string `json:"text_raw"`
	CreatedAt   string `json:"created_at"`
	User        WeiboUser `json:"user"`
	RetweetedStatus *WeiboStatus `json:"retweeted_status"`
}

type WeiboUser struct {
	ID       string `json:"id"`
	ScreenName string `json:"screen_name"`
}

type WeiboHotRes struct {
	List []WeiboHotItem `json:"list"`
}

type WeiboHotItem struct {
	Word       string `json:"word"`
	Url        string `json:"url"`
	Num        int    `json:"num"`
	Heat       int    `json:"heat"`
}

type WeiboFetcher struct {
	fetcher.BaseFetcher
}

func NewWeiboFetcher() *WeiboFetcher {
	return &WeiboFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "weibo"},
	}
}

func (f *WeiboFetcher) Fetch() ([]types.NewsItem, error) {
	baseURL := "https://s.weibo.com"
	url := baseURL + "/top/summary?cate=realtimehot"
	
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Cookie": "SUB=_2AkMWIuNSf8NxqwJRmP8dy2rhaoV2ygrEieKgfhKJJRMxHRl-yT9jqk86tRB6PaLNvQZR6zYUcYVT1zSjoSreQHidcUq7",
		"Referer": url,
	}
	
	html, err := utils.Fetch(url, utils.FetchOptions{Headers: headers})
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0)

	rows := doc.Find("#pl_top_realtimehot table tbody tr").Slice(1, -1)
	rows.Each(func(i int, s *goquery.Selection) {
		link := s.Find("td.td-02 a").FilterFunction(func(i int, el *goquery.Selection) bool {
			href, exists := el.Attr("href")
			return exists && href != "" && !strings.Contains(href, "javascript:void(0);")
		}).First()

		if link.Length() > 0 {
			title := strings.TrimSpace(link.Text())
			href, _ := link.Attr("href")

			if title != "" && href != "" {
				flag := strings.TrimSpace(s.Find("td.td-03").Text())
				var icon string
				switch flag {
				case "新":
					icon = "https://simg.s.weibo.com/moter/flags/1_0.png"
				case "热":
					icon = "https://simg.s.weibo.com/moter/flags/2_0.png"
				case "爆":
					icon = "https://simg.s.weibo.com/moter/flags/4_0.png"
				}

				extra := &types.NewsExtra{}
				if icon != "" {
					extra.Icon = icon
				}

				items = append(items, types.NewsItem{
					ID:    title,
					Title: title,
					URL:   baseURL + href,
					Extra: extra,
				})
			}
		}
	})

	return items, nil
}

type WeiboHotFetcher struct {
	fetcher.BaseFetcher
}

func NewWeiboHotFetcher() *WeiboHotFetcher {
	return &WeiboHotFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "weibo-hot"},
	}
}

func (f *WeiboHotFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://weibo.com/ajax/side/hotSearch"
	
	headers := map[string]string{
		"Accept":       "application/json, text/plain, */*",
		"Referer":      "https://weibo.com/",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"X-Requested-With": "XMLHttpRequest",
	}
	
	var res struct {
		Data struct {
			Realtime []WeiboHotItem `json:"realtime"`
		} `json:"data"`
	}
	
	if err := utils.FetchJSON(url, &res, utils.FetchOptions{Headers: headers}); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data.Realtime))
	for i, item := range res.Data.Realtime {
		items = append(items, types.NewsItem{
			ID:    strconv.Itoa(i),
			Title: item.Word,
			URL:   item.Url,
			Extra: &types.NewsExtra{
				Info: item.Num,
			},
		})
	}

	return items, nil
}
