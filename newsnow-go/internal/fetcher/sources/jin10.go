package sources

import (
	"encoding/json"
	"newsnow-go/internal/fetcher"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Jin10Item struct {
	ID        string     `json:"id"`
	Time      string     `json:"time"`
	Type      int        `json:"type"`
	Data      Jin10Data  `json:"data"`
	Important int        `json:"important"`
	Tags      []string   `json:"tags"`
	Channel   []int      `json:"channel"`
	Remark    []interface{} `json:"remark"`
}

type Jin10Data struct {
	Pic       string `json:"pic"`
	Title     string `json:"title"`
	Source    string `json:"source"`
	Content   string `json:"content"`
	SourceLink string `json:"source_link"`
	VipTitle  string `json:"vip_title"`
	Lock      bool   `json:"lock"`
	VipLevel  int    `json:"vip_level"`
	VipDesc   string `json:"vip_desc"`
}

type Jin10Fetcher struct {
	fetcher.BaseFetcher
}

func NewJin10Fetcher() *Jin10Fetcher {
	return &Jin10Fetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "jin10"},
	}
}

func (f *Jin10Fetcher) Fetch() ([]types.NewsItem, error) {
	timestamp := time.Now().UnixMilli()
	url := "https://www.jin10.com/flash_newest.js?t=" + strconv.FormatInt(timestamp, 10)

	rawData, err := utils.Fetch(url)
	if err != nil {
		return nil, err
	}

	jsonStr := regexp.MustCompile(`^var\s+newest\s*=\s*`).ReplaceAllString(rawData, "")
	jsonStr = regexp.MustCompile(`;*$`).ReplaceAllString(jsonStr, "")
	jsonStr = strings.TrimSpace(jsonStr)

	var data []Jin10Item
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0)
	for _, k := range data {
		if k.Data.Title == "" && k.Data.Content == "" {
			continue
		}
		if contains(k.Channel, 5) {
			continue
		}

		text := k.Data.Title
		if text == "" {
			text = k.Data.Content
		}
		text = strings.ReplaceAll(text, "</?b>", "")

		re := regexp.MustCompile(`^【([^】]*)】(.*)$`)
		matches := re.FindStringSubmatch(text)
		
		title := text
		desc := ""
		if len(matches) >= 3 {
			title = matches[1]
			desc = matches[2]
		}

		pubDate := ""
		if k.Time != "" {
			pubDate = utils.ParseDate(k.Time, "Asia/Shanghai").Format(time.RFC3339)
		}

		items = append(items, types.NewsItem{
			ID:      k.ID,
			Title:   title,
			URL:     "https://flash.jin10.com/detail/" + k.ID,
			PubDate: pubDate,
			Extra: &types.NewsExtra{
				Hover: desc,
				Info:  f.getImportantMark(k.Important),
			},
		})
	}

	return items, nil
}

func (f *Jin10Fetcher) getImportantMark(important int) interface{} {
	if important > 0 {
		return "☆"
	}
	return nil
}

func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	default:
		return ""
	}
}
