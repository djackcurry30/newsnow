package sources

import (
	"fmt"
	"newsnow-go/internal/fetcher"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
	"strings"
	"time"
)

type Kr36QuickFetcher struct {
	fetcher.BaseFetcher
}

func NewKr36QuickFetcher() *Kr36QuickFetcher {
	return &Kr36QuickFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "36kr-quick"},
	}
}

func (f *Kr36QuickFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://www.36kr.com/newsflashes"
	
	html, err := utils.Fetch(url)
	if err != nil {
		return nil, err
	}
	
	items := make([]types.NewsItem, 0)
	
	// 使用正则表达式匹配新闻快讯项
	// 匹配格式: <div class="newsflash-item">...<a class="item-title" href="...">标题</a>...<span class="time">时间</span>...</div>
	re := `<div class="newsflash-item"[^>]*>[\s\S]*?<a class="item-title"[^>]*href="([^"]+)"[^>]*>([^<]+)</a>[\s\S]*?<span class="time"[^>]*>([^<]+)</span>`
	matches := utils.FindAllStringSubmatch(re, html)
	
	for _, match := range matches {
		if len(match) >= 4 {
			href := match[1]
			title := match[2]
			relativeDate := match[3]
			
			if href != "" && title != "" && relativeDate != "" {
				item := types.NewsItem{
					ID:    href,
					Title: strings.TrimSpace(title),
					URL:   "https://www.36kr.com" + href,
				}
				items = append(items, item)
			}
		}
	}
	
	return items, nil
}

type Kr36RenqiFetcher struct {
	fetcher.BaseFetcher
}

func NewKr36RenqiFetcher() *Kr36RenqiFetcher {
	return &Kr36RenqiFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "36kr-renqi"},
	}
}

func (f *Kr36RenqiFetcher) Fetch() ([]types.NewsItem, error) {
	// 使用当前日期格式化
	formatted := fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	url := fmt.Sprintf("https://36kr.com/hot-list/renqi/%s/1", formatted)
	
	html, err := utils.Fetch(url)
	if err != nil {
		return nil, err
	}
	
	items := make([]types.NewsItem, 0)
	
	// 使用正则表达式匹配文章项
	// 匹配格式: <a class="article-item-title weight-bold" href="...">标题</a>
	re := `<a class="article-item-title weight-bold"[^>]*href="([^"]+)"[^>]*>([^<]+)</a>`
	matches := utils.FindAllStringSubmatch(re, html)
	
	for _, match := range matches {
		if len(match) >= 3 {
			href := match[1]
			title := match[2]
			
			if href != "" && title != "" {
				item := types.NewsItem{
					ID:    strings.TrimPrefix(href, "/a/"),
					Title: strings.TrimSpace(title),
					URL:   href,
				}
				items = append(items, item)
			}
		}
	}
	
	return items, nil
}

type Kr36Holder struct {
	Quick *Kr36QuickFetcher
	Renqi *Kr36RenqiFetcher
}

func NewKr36Holder() *Kr36Holder {
	return &Kr36Holder{
		Quick: NewKr36QuickFetcher(),
		Renqi: NewKr36RenqiFetcher(),
	}
}

func (h *Kr36Holder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"36kr":          h.Quick,
		"36kr-quick":    h.Quick,
		"36kr-renqi":    h.Renqi,
	}
}
