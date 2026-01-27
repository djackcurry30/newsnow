package sources

import (
	"encoding/json"
	"fmt"
	"newsnow-go/internal/fetcher"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
)

type BaiduHotItem struct {
	Keyword string `json:"Keyword"`
	Url     string `json:"Url"`
}

type BaiduHotRes struct {
	Data struct {
		BaiduHotList []BaiduHotItem `json:"BaiduHotList"`
	} `json:"data"`
}

type BaiduFetcher struct {
	fetcher.BaseFetcher
}

func NewBaiduFetcher() *BaiduFetcher {
	return &BaiduFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "baidu"},
	}
}

type BaiduRes struct {
	Data struct {
		Cards []struct {
			Content []struct {
				IsTop  *bool  `json:"isTop"`
				Word   string `json:"word"`
				RawURL string `json:"rawUrl"`
				Desc   string `json:"desc"`
			} `json:"content"`
		} `json:"cards"`
	} `json:"data"`
}

func (f *BaiduFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://top.baidu.com/board?tab=realtime"

	html, err := utils.Fetch(url)
	if err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0)

	// 使用正则表达式提取JSON数据
	re := `(?s)<!--s-data:(.*?)-->`
	matches := utils.FindAllStringSubmatch(re, html)
	if len(matches) == 0 || len(matches[0]) < 2 {
		return items, nil
	}

	jsonStr := matches[0][1]

	// 解析JSON
	var res BaiduRes
	if err := json.Unmarshal([]byte(jsonStr), &res); err != nil {
		return items, nil
	}

	// 提取热搜项
	if len(res.Data.Cards) == 0 || len(res.Data.Cards[0].Content) == 0 {
		return items, nil
	}

	for _, item := range res.Data.Cards[0].Content {
		// 跳过置顶项
		if item.IsTop != nil && *item.IsTop {
			continue
		}

		if item.Word == "" || item.RawURL == "" {
			continue
		}

		extra := &types.NewsExtra{}
		if item.Desc != "" {
			extra.Hover = item.Desc
		}

		items = append(items, types.NewsItem{
			ID:    item.RawURL,
			Title: item.Word,
			URL:   item.RawURL,
			Extra: extra,
		})
	}

	return items, nil
}

type DouyinFetcher struct {
	fetcher.BaseFetcher
}

func NewDouyinFetcher() *DouyinFetcher {
	return &DouyinFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "douyin"},
	}
}

type DouyinRes struct {
	Data struct {
		WordList []struct {
			SentenceID string `json:"sentence_id"`
			Word       string `json:"word"`
			EventTime  string `json:"event_time"`
			HotValue   string `json:"hot_value"`
		} `json:"word_list"`
	} `json:"data"`
}

func (f *DouyinFetcher) Fetch() ([]types.NewsItem, error) {
	// 先获取cookie
	loginURL := "https://login.douyin.com/"
	_, err := utils.Fetch(loginURL)
	if err != nil {
		return nil, err
	}

	// 使用API获取热搜数据
	url := "https://www.douyin.com/aweme/v1/web/hot/search/list/?device_platform=webapp&aid=6383&channel=channel_pc_web&detail_list=1"

	var res DouyinRes
	if err := utils.FetchJSON(url, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data.WordList))
	for _, item := range res.Data.WordList {
		items = append(items, types.NewsItem{
			ID:    item.SentenceID,
			Title: item.Word,
			URL:   fmt.Sprintf("https://www.douyin.com/hot/%s", item.SentenceID),
		})
	}

	return items, nil
}

type HupuFetcher struct {
	fetcher.BaseFetcher
}

func NewHupuFetcher() *HupuFetcher {
	return &HupuFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "hupu"},
	}
}

func (f *HupuFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://bbs.hupu.com/topic-daily-hot"

	html, err := utils.Fetch(url)
	if err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0)

	// 使用正则表达式匹配热门话题项
	// 匹配格式: <a href="/xxx.html" ... class="p-title" ...>标题</a>
	re := `<a href="(/[0-9]+\.html)"[^>]*class="p-title"[^>]*>([^<]+)</a>`
	matches := utils.FindAllStringSubmatch(re, html)

	for _, match := range matches {
		if len(match) >= 3 {
			path := match[1]
			title := match[2]

			item := types.NewsItem{
				ID:    path,
				Title: title,
				URL:   "https://bbs.hupu.com" + path,
			}
			items = append(items, item)
		}
	}

	return items, nil
}

type TiebaFetcher struct {
	fetcher.BaseFetcher
}

func NewTiebaFetcher() *TiebaFetcher {
	return &TiebaFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "tieba"},
	}
}

type TiebaRes struct {
	Data struct {
		BangTopic struct {
			TopicList []struct {
				TopicID    int64  `json:"topic_id"`
				TopicName  string `json:"topic_name"`
				CreateTime int64  `json:"create_time"`
				TopicURL   string `json:"topic_url"`
			} `json:"topic_list"`
		} `json:"bang_topic"`
	} `json:"data"`
}

func (f *TiebaFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://tieba.baidu.com/hottopic/browse/topicList"

	var res TiebaRes
	if err := utils.FetchJSON(url, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data.BangTopic.TopicList))

	for _, item := range res.Data.BangTopic.TopicList {
		newsItem := types.NewsItem{
			ID:    fmt.Sprintf("%d", item.TopicID),
			Title: item.TopicName,
			URL:   item.TopicURL,
		}
		items = append(items, newsItem)
	}

	return items, nil
}

type ToutiaoFetcher struct {
	fetcher.BaseFetcher
}

func NewToutiaoFetcher() *ToutiaoFetcher {
	return &ToutiaoFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "toutiao"},
	}
}

type ToutiaoRes struct {
	Data []struct {
		ClusterIDStr string `json:"ClusterIdStr"`
		Title        string `json:"Title"`
		HotValue     string `json:"HotValue"`
		Image        struct {
			URL string `json:"url"`
		} `json:"Image"`
		LabelURI *struct {
			URL string `json:"url"`
		} `json:"LabelUri"`
	} `json:"data"`
}

func (f *ToutiaoFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://www.toutiao.com/hot-event/hot-board/?origin=toutiao_pc"

	var res ToutiaoRes
	if err := utils.FetchJSON(url, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data))
	for _, item := range res.Data {
		extra := &types.NewsExtra{}
		if item.LabelURI != nil {
			extra.Icon = item.LabelURI.URL
		}

		items = append(items, types.NewsItem{
			ID:    item.ClusterIDStr,
			Title: item.Title,
			URL:   fmt.Sprintf("https://www.toutiao.com/trending/%s/", item.ClusterIDStr),
			Extra: extra,
		})
	}

	return items, nil
}

type WeiboHotHolder struct {
	Search  *WeiboFetcher
	Hotlist *WeiboHotFetcher
}

func NewWeiboHotHolder() *WeiboHotHolder {
	return &WeiboHotHolder{
		Search:  NewWeiboFetcher(),
		Hotlist: NewWeiboHotFetcher(),
	}
}

func (h *WeiboHotHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"weibo":     h.Search,
		"weibo-hot": h.Hotlist,
	}
}
