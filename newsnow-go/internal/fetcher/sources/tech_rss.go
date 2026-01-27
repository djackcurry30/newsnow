package sources

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"newsnow-go/internal/fetcher"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
)

type IthomeFetcher struct {
	fetcher.BaseFetcher
}

func NewIthomeFetcher() *IthomeFetcher {
	return &IthomeFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "ithome"},
	}
}

func (f *IthomeFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://www.ithome.com/rss/", false).Fetch()
}

type SspaiFetcher struct {
	fetcher.BaseFetcher
}

func NewSspaiFetcher() *SspaiFetcher {
	return &SspaiFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "sspai"},
	}
}

func (f *SspaiFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://sspai.com/feed", false).Fetch()
}

type JuejinFetcher struct {
	fetcher.BaseFetcher
}

func NewJuejinFetcher() *JuejinFetcher {
	return &JuejinFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "juejin"},
	}
}

type JuejinRes struct {
	Data []struct {
		Content struct {
			Title     string `json:"title"`
			ContentID string `json:"content_id"`
		} `json:"content"`
	} `json:"data"`
}

func (f *JuejinFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://api.juejin.cn/content_api/v1/content/article_rank?category_id=1&type=hot&spider=0"

	var res JuejinRes
	if err := utils.FetchJSON(url, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data))

	for _, item := range res.Data {
		newsItem := types.NewsItem{
			ID:    item.Content.ContentID,
			Title: item.Content.Title,
			URL:   fmt.Sprintf("https://juejin.cn/post/%s", item.Content.ContentID),
		}
		items = append(items, newsItem)
	}

	return items, nil
}

type SolidotFetcher struct {
	fetcher.BaseFetcher
}

func NewSolidotFetcher() *SolidotFetcher {
	return &SolidotFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "solidot"},
	}
}

func (f *SolidotFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://www.solidot.org/index.rss", false).Fetch()
}

type HackerNewsFetcher struct {
	fetcher.BaseFetcher
}

func NewHackerNewsFetcher() *HackerNewsFetcher {
	return &HackerNewsFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "hackernews"},
	}
}

func (f *HackerNewsFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://news.ycombinator.com/rss", false).Fetch()
}

type ProductHuntFetcher struct {
	fetcher.BaseFetcher
}

func NewProductHuntFetcher() *ProductHuntFetcher {
	return &ProductHuntFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "producthunt"},
	}
}

func (f *ProductHuntFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://www.producthunt.com/feed", false).Fetch()
}

type BilibiliHotSearchFetcher struct {
	fetcher.BaseFetcher
}

func NewBilibiliHotSearchFetcher() *BilibiliHotSearchFetcher {
	return &BilibiliHotSearchFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "bilibili"},
	}
}

func (f *BilibiliHotSearchFetcher) Fetch() ([]types.NewsItem, error) {
	apiURL := "https://s.search.bilibili.com/main/hotword?limit=30"

	type WapRes struct {
		Code int `json:"code"`
		List []struct {
			HotID     int     `json:"hot_id"`
			Keyword   string  `json:"keyword"`
			ShowName  string  `json:"show_name"`
			Score     float64 `json:"score"`
			GotoType  int     `json:"goto_type"`
			GotoValue string  `json:"goto_value"`
			Icon      string  `json:"icon"`
			HeatScore float64 `json:"heat_score"`
		} `json:"list"`
	}

	var res WapRes
	if err := utils.FetchJSON(apiURL, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.List))
	for _, item := range res.List {
		items = append(items, types.NewsItem{
			ID:    fmt.Sprintf("%d", item.HotID),
			Title: item.ShowName,
			URL:   fmt.Sprintf("https://search.bilibili.com/all?keyword=%s", url.QueryEscape(item.Keyword)),
			Extra: &types.NewsExtra{
				Icon: item.Icon,
			},
		})
	}

	return items, nil
}

type BilibiliHotVideoFetcher struct {
	fetcher.BaseFetcher
}

func NewBilibiliHotVideoFetcher() *BilibiliHotVideoFetcher {
	return &BilibiliHotVideoFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "bilibili-hot-video"},
	}
}

func (f *BilibiliHotVideoFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://api.bilibili.com/x/web-interface/popular"

	type HotVideoRes struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		TTL     int    `json:"ttl"`
		Data    struct {
			List []struct {
				AID     int    `json:"aid"`
				Title   string `json:"title"`
				PubDate int64  `json:"pubdate"`
				Desc    string `json:"desc"`
				Owner   struct {
					Name string `json:"name"`
				} `json:"owner"`
				Stat struct {
					View int `json:"view"`
					Like int `json:"like"`
				} `json:"stat"`
				BVID string `json:"bvid"`
				Pic  string `json:"pic"`
			} `json:"list"`
		} `json:"data"`
	}

	var res HotVideoRes
	if err := utils.FetchJSON(url, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data.List))
	for _, video := range res.Data.List {
		pubDate := ""
		if video.PubDate > 0 {
			pubDate = time.Unix(video.PubDate, 0).Format(time.RFC3339)
		}

		info := fmt.Sprintf("%s · %d观看 · %d点赞", video.Owner.Name, video.Stat.View, video.Stat.Like)

		items = append(items, types.NewsItem{
			ID:      video.BVID,
			Title:   video.Title,
			URL:     fmt.Sprintf("https://www.bilibili.com/video/%s", video.BVID),
			PubDate: pubDate,
			Extra: &types.NewsExtra{
				Info:  info,
				Hover: video.Desc,
				Icon:  video.Pic,
			},
		})
	}

	return items, nil
}

type BilibiliRankingFetcher struct {
	fetcher.BaseFetcher
}

func NewBilibiliRankingFetcher() *BilibiliRankingFetcher {
	return &BilibiliRankingFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "bilibili-ranking"},
	}
}

func (f *BilibiliRankingFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://api.bilibili.com/x/web-interface/ranking/v2"

	type HotVideoRes struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		TTL     int    `json:"ttl"`
		Data    struct {
			List []struct {
				AID     int    `json:"aid"`
				Title   string `json:"title"`
				PubDate int64  `json:"pubdate"`
				Desc    string `json:"desc"`
				Owner   struct {
					Name string `json:"name"`
				} `json:"owner"`
				Stat struct {
					View int `json:"view"`
					Like int `json:"like"`
				} `json:"stat"`
				BVID string `json:"bvid"`
				Pic  string `json:"pic"`
			} `json:"list"`
		} `json:"data"`
	}

	var res HotVideoRes
	if err := utils.FetchJSON(url, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data.List))
	for _, video := range res.Data.List {
		pubDate := ""
		if video.PubDate > 0 {
			pubDate = time.Unix(video.PubDate, 0).Format(time.RFC3339)
		}

		info := fmt.Sprintf("%s · %d观看 · %d点赞", video.Owner.Name, video.Stat.View, video.Stat.Like)

		items = append(items, types.NewsItem{
			ID:      video.BVID,
			Title:   video.Title,
			URL:     fmt.Sprintf("https://www.bilibili.com/video/%s", video.BVID),
			PubDate: pubDate,
			Extra: &types.NewsExtra{
				Info:  info,
				Hover: video.Desc,
				Icon:  video.Pic,
			},
		})
	}

	return items, nil
}

type BilibiliHolder struct {
	HotSearch *BilibiliHotSearchFetcher
	HotVideo  *BilibiliHotVideoFetcher
	Ranking   *BilibiliRankingFetcher
}

func NewBilibiliHolder() *BilibiliHolder {
	return &BilibiliHolder{
		HotSearch: NewBilibiliHotSearchFetcher(),
		HotVideo:  NewBilibiliHotVideoFetcher(),
		Ranking:   NewBilibiliRankingFetcher(),
	}
}

func (h *BilibiliHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"bilibili":            h.HotSearch,
		"bilibili-hot-search": h.HotSearch,
		"bilibili-hot-video":  h.HotVideo,
		"bilibili-ranking":    h.Ranking,
	}
}

type ThePaperFetcher struct {
	fetcher.BaseFetcher
}

func NewThePaperFetcher() *ThePaperFetcher {
	return &ThePaperFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "thepaper"},
	}
}

type ThePaperRes struct {
	Data struct {
		HotNews []struct {
			ContID      string `json:"contId"`
			Name        string `json:"name"`
			PubTimeLong int64  `json:"pubTimeLong"`
		} `json:"hotNews"`
	} `json:"data"`
}

func (f *ThePaperFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://cache.thepaper.cn/contentapi/wwwIndex/rightSidebar"

	var res ThePaperRes
	if err := utils.FetchJSON(url, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data.HotNews))
	for _, item := range res.Data.HotNews {
		items = append(items, types.NewsItem{
			ID:    item.ContID,
			Title: item.Name,
			URL:   fmt.Sprintf("https://www.thepaper.cn/newsDetail_forward_%s", item.ContID),
		})
	}

	return items, nil
}

type ZaobaoFetcher struct {
	fetcher.BaseFetcher
}

func NewZaobaoFetcher() *ZaobaoFetcher {
	return &ZaobaoFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "zaobao"},
	}
}

func (f *ZaobaoFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://www.zaobao.com.sg/rss/feed", false).Fetch()
}

type WallStreetCNLiveFetcher struct {
	fetcher.BaseFetcher
}

func NewWallStreetCNLiveFetcher() *WallStreetCNLiveFetcher {
	return &WallStreetCNLiveFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "wallstreetcn"},
	}
}

func (f *WallStreetCNLiveFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://api-one.wallstcn.com/apiv1/content/lives?channel=global-channel&limit=30"

	type Item struct {
		URI          string `json:"uri"`
		ID           int    `json:"id"`
		Title        string `json:"title,omitempty"`
		ContentText  string `json:"content_text"`
		ContentShort string `json:"content_short"`
		DisplayTime  int64  `json:"display_time"`
		Type         string `json:"type,omitempty"`
	}

	type LiveRes struct {
		Data struct {
			Items []Item `json:"items"`
		} `json:"data"`
	}

	var res LiveRes
	if err := utils.FetchJSON(url, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data.Items))
	for _, item := range res.Data.Items {
		title := item.Title
		if title == "" {
			title = item.ContentText
		}

		pubDate := ""
		if item.DisplayTime > 0 {
			pubDate = time.Unix(item.DisplayTime, 0).Format(time.RFC3339)
		}

		items = append(items, types.NewsItem{
			ID:      fmt.Sprintf("%d", item.ID),
			Title:   title,
			URL:     item.URI,
			PubDate: pubDate,
			Extra: &types.NewsExtra{
				Date:  pubDate,
				Hover: item.ContentShort,
			},
		})
	}

	return items, nil
}

type WallStreetCNNewsFetcher struct {
	fetcher.BaseFetcher
}

func NewWallStreetCNNewsFetcher() *WallStreetCNNewsFetcher {
	return &WallStreetCNNewsFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "wallstreetcn-news"},
	}
}

func (f *WallStreetCNNewsFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://api-one.wallstcn.com/apiv1/content/information-flow?channel=global-channel&accept=article&limit=30"

	type Item struct {
		URI          string `json:"uri"`
		ID           int    `json:"id"`
		Title        string `json:"title,omitempty"`
		ContentText  string `json:"content_text"`
		ContentShort string `json:"content_short"`
		DisplayTime  int64  `json:"display_time"`
		Type         string `json:"type,omitempty"`
	}

	type ResourceItem struct {
		ResourceType string `json:"resource_type,omitempty"`
		Resource     Item   `json:"resource"`
	}

	type NewsRes struct {
		Data struct {
			Items []ResourceItem `json:"items"`
		} `json:"data"`
	}

	var res NewsRes
	if err := utils.FetchJSON(url, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0)
	for _, resourceItem := range res.Data.Items {
		if resourceItem.ResourceType == "theme" || resourceItem.ResourceType == "ad" || resourceItem.Resource.Type == "live" || resourceItem.Resource.URI == "" {
			continue
		}

		item := resourceItem.Resource
		title := item.Title
		if title == "" {
			title = item.ContentShort
		}

		pubDate := ""
		if item.DisplayTime > 0 {
			pubDate = time.Unix(item.DisplayTime, 0).Format(time.RFC3339)
		}

		items = append(items, types.NewsItem{
			ID:      fmt.Sprintf("%d", item.ID),
			Title:   title,
			URL:     item.URI,
			PubDate: pubDate,
			Extra: &types.NewsExtra{
				Date:  pubDate,
				Hover: item.ContentShort,
			},
		})
	}

	return items, nil
}

type WallStreetCNHotFetcher struct {
	fetcher.BaseFetcher
}

func NewWallStreetCNHotFetcher() *WallStreetCNHotFetcher {
	return &WallStreetCNHotFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "wallstreetcn-hot"},
	}
}

func (f *WallStreetCNHotFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://api-one.wallstcn.com/apiv1/content/articles/hot?period=all"

	type Item struct {
		URI         string `json:"uri"`
		ID          int    `json:"id"`
		Title       string `json:"title,omitempty"`
		DisplayTime int64  `json:"display_time"`
	}

	type HotRes struct {
		Data struct {
			DayItems []Item `json:"day_items"`
		} `json:"data"`
	}

	var res HotRes
	if err := utils.FetchJSON(url, &res); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Data.DayItems))
	for _, item := range res.Data.DayItems {
		items = append(items, types.NewsItem{
			ID:    fmt.Sprintf("%d", item.ID),
			Title: item.Title,
			URL:   item.URI,
		})
	}

	return items, nil
}

type WallStreetCNHolder struct {
	Live *WallStreetCNLiveFetcher
	News *WallStreetCNNewsFetcher
	Hot  *WallStreetCNHotFetcher
}

func NewWallStreetCNHolder() *WallStreetCNHolder {
	return &WallStreetCNHolder{
		Live: NewWallStreetCNLiveFetcher(),
		News: NewWallStreetCNNewsFetcher(),
		Hot:  NewWallStreetCNHotFetcher(),
	}
}

func (h *WallStreetCNHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"wallstreetcn":       h.Live,
		"wallstreetcn-quick": h.Live,
		"wallstreetcn-news":  h.News,
		"wallstreetcn-hot":   h.Hot,
	}
}

type CLSFetcher struct {
	fetcher.BaseFetcher
}

func NewCLSFetcher() *CLSFetcher {
	return &CLSFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "cls"},
	}
}

func (f *CLSFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://cls.cn/css cls/rss", false).Fetch()
}

type CLSDepthFetcher struct {
	fetcher.BaseFetcher
}

func NewCLSDepthFetcher() *CLSDepthFetcher {
	return &CLSDepthFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "cls-depth"},
	}
}

func (f *CLSDepthFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://cls.cn/css cls/rss", false).Fetch()
}

type CLSHotFetcher struct {
	fetcher.BaseFetcher
}

func NewCLSHotFetcher() *CLSHotFetcher {
	return &CLSHotFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "cls-hot"},
	}
}

func (f *CLSHotFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://cls.cn/css cls/rss", false).Fetch()
}

type CLSTelegraphFetcher struct {
	fetcher.BaseFetcher
}

func NewCLSTelegraphFetcher() *CLSTelegraphFetcher {
	return &CLSTelegraphFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "cls-telegraph"},
	}
}

func (f *CLSTelegraphFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://cls.cn/css cls/rss", false).Fetch()
}

type CLSHolder struct {
	Main      *CLSFetcher
	Depth     *CLSDepthFetcher
	Hot       *CLSHotFetcher
	Telegraph *CLSTelegraphFetcher
}

func NewCLSHolder() *CLSHolder {
	return &CLSHolder{
		Main:      NewCLSFetcher(),
		Depth:     NewCLSDepthFetcher(),
		Hot:       NewCLSHotFetcher(),
		Telegraph: NewCLSTelegraphFetcher(),
	}
}

func (h *CLSHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"cls":           h.Main,
		"cls-depth":     h.Depth,
		"cls-hot":       h.Hot,
		"cls-telegraph": h.Telegraph,
	}
}

type XueqiuFetcher struct {
	fetcher.BaseFetcher
}

func NewXueqiuFetcher() *XueqiuFetcher {
	return &XueqiuFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "xueqiu"},
	}
}

func (f *XueqiuFetcher) Fetch() ([]types.NewsItem, error) {
	sorted := true
	return fetcher.DefineRSSHubSource("/xueqiu/user/1019667", fetcher.RSSHubOptions{Sorted: &sorted}, false).Fetch()
}

type XueqiuHotstockFetcher struct {
	fetcher.BaseFetcher
}

func NewXueqiuHotstockFetcher() *XueqiuHotstockFetcher {
	return &XueqiuHotstockFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "xueqiu-hotstock"},
	}
}

func (f *XueqiuHotstockFetcher) Fetch() ([]types.NewsItem, error) {
	sorted := true
	return fetcher.DefineRSSHubSource("/xueqiu/user/1019667", fetcher.RSSHubOptions{Sorted: &sorted}, false).Fetch()
}

type XueqiuHolder struct {
	Main     *XueqiuFetcher
	Hotstock *XueqiuHotstockFetcher
}

func NewXueqiuHolder() *XueqiuHolder {
	return &XueqiuHolder{
		Main:     NewXueqiuFetcher(),
		Hotstock: NewXueqiuHotstockFetcher(),
	}
}

func (h *XueqiuHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"xueqiu":          h.Main,
		"xueqiu-hotstock": h.Hotstock,
	}
}

type GeLongHuiFetcher struct {
	fetcher.BaseFetcher
}

func NewGeLongHuiFetcher() *GeLongHuiFetcher {
	return &GeLongHuiFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "gelonghui"},
	}
}

func (f *GeLongHuiFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://www.gelonghui.com/rss", false).Fetch()
}

type FastBullExpressFetcher struct {
	fetcher.BaseFetcher
}

func NewFastBullExpressFetcher() *FastBullExpressFetcher {
	return &FastBullExpressFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "fastbull"},
	}
}

func (f *FastBullExpressFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://fastbull.org/feed", false).Fetch()
}

type FastBullNewsFetcher struct {
	fetcher.BaseFetcher
}

func NewFastBullNewsFetcher() *FastBullNewsFetcher {
	return &FastBullNewsFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "fastbull-news"},
	}
}

func (f *FastBullNewsFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://fastbull.org/feed", false).Fetch()
}

type FastBullHolder struct {
	Express *FastBullExpressFetcher
	News    *FastBullNewsFetcher
}

func NewFastBullHolder() *FastBullHolder {
	return &FastBullHolder{
		Express: NewFastBullExpressFetcher(),
		News:    NewFastBullNewsFetcher(),
	}
}

func (h *FastBullHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"fastbull":         h.Express,
		"fastbull-express": h.Express,
		"fastbull-news":    h.News,
	}
}

type MKTNewsFetcher struct {
	fetcher.BaseFetcher
}

func NewMKTNewsFetcher() *MKTNewsFetcher {
	return &MKTNewsFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "mktnews"},
	}
}

func (f *MKTNewsFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://mktnews.com/feed", false).Fetch()
}

type MKTNewsFlashFetcher struct {
	fetcher.BaseFetcher
}

func NewMKTNewsFlashFetcher() *MKTNewsFlashFetcher {
	return &MKTNewsFlashFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "mktnews-flash"},
	}
}

func (f *MKTNewsFlashFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://mktnews.com/feed", false).Fetch()
}

type MKTNewsHolder struct {
	Main  *MKTNewsFetcher
	Flash *MKTNewsFlashFetcher
}

func NewMKTNewsHolder() *MKTNewsHolder {
	return &MKTNewsHolder{
		Main:  NewMKTNewsFetcher(),
		Flash: NewMKTNewsFlashFetcher(),
	}
}

func (h *MKTNewsHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"mktnews":       h.Main,
		"mktnews-flash": h.Flash,
	}
}

type PCBetaFetcher struct {
	fetcher.BaseFetcher
}

func NewPCBetaFetcher() *PCBetaFetcher {
	return &PCBetaFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "pcbeta"},
	}
}

func (f *PCBetaFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://forum.pceva.com.cn/rss", false).Fetch()
}

type PCBetaWindows11Fetcher struct {
	fetcher.BaseFetcher
}

func NewPCBetaWindows11Fetcher() *PCBetaWindows11Fetcher {
	return &PCBetaWindows11Fetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "pcbeta-windows11"},
	}
}

func (f *PCBetaWindows11Fetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://forum.pceva.com.cn/rss", false).Fetch()
}

type PCBetaHolder struct {
	Main      *PCBetaFetcher
	Windows11 *PCBetaWindows11Fetcher
}

func NewPCBetaHolder() *PCBetaHolder {
	return &PCBetaHolder{
		Main:      NewPCBetaFetcher(),
		Windows11: NewPCBetaWindows11Fetcher(),
	}
}

func (h *PCBetaHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"pcbeta":           h.Main,
		"pcbeta-windows11": h.Windows11,
	}
}

type FreebufFetcher struct {
	fetcher.BaseFetcher
}

func NewFreebufFetcher() *FreebufFetcher {
	return &FreebufFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "freebuf"},
	}
}

func (f *FreebufFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://www.freebuf.com/articles/rss", false).Fetch()
}

type SteamFetcher struct {
	fetcher.BaseFetcher
}

func NewSteamFetcher() *SteamFetcher {
	return &SteamFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "steam"},
	}
}

func (f *SteamFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://store.steampowered.com/feeds/news/app/105600", false).Fetch()
}

type TencentFetcher struct {
	fetcher.BaseFetcher
}

func NewTencentFetcher() *TencentFetcher {
	return &TencentFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "tencent"},
	}
}

func (f *TencentFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://new.qq.com/rain/feed", false).Fetch()
}

type TencentHotFetcher struct {
	fetcher.BaseFetcher
}

func NewTencentHotFetcher() *TencentHotFetcher {
	return &TencentHotFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "tencent-hot"},
	}
}

func (f *TencentHotFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://new.qq.com/rain/feed", false).Fetch()
}

type TencentHolder struct {
	Main *TencentFetcher
	Hot  *TencentHotFetcher
}

func NewTencentHolder() *TencentHolder {
	return &TencentHolder{
		Main: NewTencentFetcher(),
		Hot:  NewTencentHotFetcher(),
	}
}

func (h *TencentHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"tencent":     h.Main,
		"tencent-hot": h.Hot,
	}
}

type QQVideoFetcher struct {
	fetcher.BaseFetcher
}

func NewQQVideoFetcher() *QQVideoFetcher {
	return &QQVideoFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "qqvideo"},
	}
}

func (f *QQVideoFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://v.qq.com/x/feeds/bangumi", false).Fetch()
}

type QQVideoTVHotsearchFetcher struct {
	fetcher.BaseFetcher
}

func NewQQVideoTVHotsearchFetcher() *QQVideoTVHotsearchFetcher {
	return &QQVideoTVHotsearchFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "qqvideo-tv-hotsearch"},
	}
}

func (f *QQVideoTVHotsearchFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://v.qq.com/x/feeds/bangumi", false).Fetch()
}

type QQVideoHolder struct {
	Main        *QQVideoFetcher
	TVHotsearch *QQVideoTVHotsearchFetcher
}

func NewQQVideoHolder() *QQVideoHolder {
	return &QQVideoHolder{
		Main:        NewQQVideoFetcher(),
		TVHotsearch: NewQQVideoTVHotsearchFetcher(),
	}
}

func (h *QQVideoHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"qqvideo":              h.Main,
		"qqvideo-tv-hotsearch": h.TVHotsearch,
	}
}

type IqiyiFetcher struct {
	fetcher.BaseFetcher
}

func NewIqiyiFetcher() *IqiyiFetcher {
	return &IqiyiFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "iqiyi"},
	}
}

type IqiyiVideoInfo struct {
	Creator []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"creator"`
	Contributor []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"contributor"`
	ShowDate    string `json:"showDate"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	EntityID    int64  `json:"entity_id"`
	BackImage   string `json:"back_image"`
	DisplayName string `json:"display_name"`
	Title       string `json:"title"`
	ShowTime    string `json:"show_time"`
	Desc        string `json:"desc"`
	RankPrefix  string `json:"rank_prefix"`
	PageURL     string `json:"page_url"`
}

type IqiyiResp struct {
	Code  int `json:"code"`
	Items []struct {
		Order int `json:"order"`
		Temp  any `json:"temp"`
		Video []struct {
			BasisDataUrls string           `json:"basisDataUrls"`
			BlockID       string           `json:"block_id"`
			CardSource    string           `json:"card_source"`
			Links         []any            `json:"links"`
			Data          []IqiyiVideoInfo `json:"data"`
			Adverts       []any            `json:"adverts"`
			Config        []any            `json:"config"`
		} `json:"video"`
	} `json:"items"`
}

func (f *IqiyiFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://mesh.if.iqiyi.com/portal/lw/v7/channel/card/videoTab?channelName=recommend&data_source=v7_rec_sec_hot_rank_list&tempId=85&count=30&block_id=hot_ranklist&device=14a4b5ba98e790dce6dc07482447cf48&from=webapp"

	headers := map[string]string{
		"Referer": "https://www.iqiyi.com",
	}

	var resp IqiyiResp
	if err := utils.FetchJSON(url, &resp, utils.FetchOptions{Headers: headers}); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0)

	// 提取视频数据
	if len(resp.Items) > 0 {
		for _, videoItem := range resp.Items[0].Video {
			for _, video := range videoItem.Data {
				extra := &types.NewsExtra{
					Info:  video.Desc,
					Hover: video.Description,
				}

				items = append(items, types.NewsItem{
					ID:      fmt.Sprintf("%d", video.EntityID),
					Title:   video.Title,
					URL:     video.PageURL,
					PubDate: video.ShowDate,
					Extra:   extra,
				})
			}
		}
	}

	return items, nil
}

type IqiyiHotRanklistFetcher struct {
	fetcher.BaseFetcher
}

func NewIqiyiHotRanklistFetcher() *IqiyiHotRanklistFetcher {
	return &IqiyiHotRanklistFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "iqiyi-hot-ranklist"},
	}
}

func (f *IqiyiHotRanklistFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://mesh.if.iqiyi.com/portal/lw/v7/channel/card/videoTab?channelName=recommend&data_source=v7_rec_sec_hot_rank_list&tempId=85&count=30&block_id=hot_ranklist&device=14a4b5ba98e790dce6dc07482447cf48&from=webapp"

	headers := map[string]string{
		"Referer": "https://www.iqiyi.com",
	}

	var resp IqiyiResp
	if err := utils.FetchJSON(url, &resp, utils.FetchOptions{Headers: headers}); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0)

	// 提取视频数据
	if len(resp.Items) > 0 {
		for _, videoItem := range resp.Items[0].Video {
			for _, video := range videoItem.Data {
				extra := &types.NewsExtra{
					Info:  video.Desc,
					Hover: video.Description,
				}

				items = append(items, types.NewsItem{
					ID:      fmt.Sprintf("%d", video.EntityID),
					Title:   video.Title,
					URL:     video.PageURL,
					PubDate: video.ShowDate,
					Extra:   extra,
				})
			}
		}
	}

	return items, nil
}

type IqiyiHolder struct {
	Main        *IqiyiFetcher
	HotRanklist *IqiyiHotRanklistFetcher
}

func NewIqiyiHolder() *IqiyiHolder {
	return &IqiyiHolder{
		Main:        NewIqiyiFetcher(),
		HotRanklist: NewIqiyiHotRanklistFetcher(),
	}
}

func (h *IqiyiHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"iqiyi":              h.Main,
		"iqiyi-hot-ranklist": h.HotRanklist,
	}
}

type DoubanFetcher struct {
	fetcher.BaseFetcher
}

func NewDoubanFetcher() *DoubanFetcher {
	return &DoubanFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "douban"},
	}
}

func (f *DoubanFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://m.douban.com/rexxar/api/v2/subject/recent_hot/movie"

	type MovieItem struct {
		Rating struct {
			Count     int     `json:"count"`
			Max       int     `json:"max"`
			StarCount float64 `json:"star_count"`
			Value     float64 `json:"value"`
		} `json:"rating"`
		Title string `json:"title"`
		Pic   struct {
			Large  string `json:"large"`
			Normal string `json:"normal"`
		} `json:"pic"`
		IsNew        bool   `json:"is_new"`
		URI          string `json:"uri"`
		EpisodesInfo string `json:"episodes_info"`
		CardSubtitle string `json:"card_subtitle"`
		Type         string `json:"type"`
		ID           string `json:"id"`
	}

	type HotMoviesRes struct {
		Category      string      `json:"category"`
		Tags          []struct{}  `json:"tags"`
		Items         []MovieItem `json:"items"`
		RecommendTags []struct{}  `json:"recommend_tags"`
		Total         int         `json:"total"`
		Type          string      `json:"type"`
	}

	var res HotMoviesRes
	headers := map[string]string{
		"Referer": "https://movie.douban.com/",
		"Accept":  "application/json, text/plain, */*",
	}

	if err := utils.FetchJSON(url, &res, utils.FetchOptions{Headers: headers}); err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0, len(res.Items))
	for _, movie := range res.Items {
		info := movie.CardSubtitle
		// 只取前三个信息
		parts := strings.Split(info, " / ")
		if len(parts) > 3 {
			info = strings.Join(parts[:3], " / ")
		}

		items = append(items, types.NewsItem{
			ID:    movie.ID,
			Title: movie.Title,
			URL:   "https://movie.douban.com/subject/" + movie.ID,
			Extra: &types.NewsExtra{
				Info:  info,
				Hover: movie.CardSubtitle,
			},
		})
	}

	return items, nil
}

type NowCoderFetcher struct {
	fetcher.BaseFetcher
}

func NewNowCoderFetcher() *NowCoderFetcher {
	return &NowCoderFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "nowcoder"},
	}
}

func (f *NowCoderFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://www.nowcoder.com/discuss/recent/good?type=0", false).Fetch()
}

type IfengFetcher struct {
	fetcher.BaseFetcher
}

func NewIfengFetcher() *IfengFetcher {
	return &IfengFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "ifeng"},
	}
}

func (f *IfengFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://www.ifeng.com/rss", false).Fetch()
}

type ChongBuluoLatestFetcher struct {
	fetcher.BaseFetcher
}

func NewChongBuluoLatestFetcher() *ChongBuluoLatestFetcher {
	return &ChongBuluoLatestFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "chongbuluo"},
	}
}

func (f *ChongBuluoLatestFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://chongbuluo.com/rss", false).Fetch()
}

type ChongBuluoHotFetcher struct {
	fetcher.BaseFetcher
}

func NewChongBuluoHotFetcher() *ChongBuluoHotFetcher {
	return &ChongBuluoHotFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "chongbuluo-hot"},
	}
}

func (f *ChongBuluoHotFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://chongbuluo.com/rss", false).Fetch()
}

type ChongBuluoHolder struct {
	Latest *ChongBuluoLatestFetcher
	Hot    *ChongBuluoHotFetcher
}

func NewChongBuluoHolder() *ChongBuluoHolder {
	return &ChongBuluoHolder{
		Latest: NewChongBuluoLatestFetcher(),
		Hot:    NewChongBuluoHotFetcher(),
	}
}

func (h *ChongBuluoHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"chongbuluo":        h.Latest,
		"chongbuluo-latest": h.Latest,
		"chongbuluo-hot":    h.Hot,
	}
}

type KuaishouFetcher struct {
	fetcher.BaseFetcher
}

func NewKuaishouFetcher() *KuaishouFetcher {
	return &KuaishouFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "kuaishou"},
	}
}

type KuaishouApolloState struct {
	DefaultClient map[string]interface{} `json:"defaultClient"`
}

type KuaishouHotRankData struct {
	Result      int    `json:"result"`
	Pcursor     string `json:"pcursor"`
	WebPageArea string `json:"webPageArea"`
	Items       []struct {
		Type      string `json:"type"`
		Generated bool   `json:"generated"`
		ID        string `json:"id"`
		Typename  string `json:"typename"`
	} `json:"items"`
}

type KuaishouHotItem struct {
	Name    string `json:"name"`
	IconURL string `json:"iconUrl"`
	TagType string `json:"tagType"`
}

func (f *KuaishouFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://www.kuaishou.com/?isHome=1"

	html, err := utils.Fetch(url)
	if err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0)

	// 提取window.__APOLLO_STATE__中的数据
	re := `window\.__APOLLO_STATE__\s*=\s*(\{.+?\});`
	matches := utils.FindAllStringSubmatch(re, html)
	if len(matches) == 0 || len(matches[0]) < 2 {
		return items, nil
	}

	jsonStr := matches[0][1]

	// 解析JSON数据
	var apolloState KuaishouApolloState
	if err := json.Unmarshal([]byte(jsonStr), &apolloState); err != nil {
		return items, nil
	}

	// 获取热榜数据ID
	hotRankQueryKey := "visionHotRank({\"page\":\"home\"})"
	rootQuery, ok := apolloState.DefaultClient["ROOT_QUERY"].(map[string]interface{})
	if !ok {
		return items, nil
	}

	hotRankQuery, ok := rootQuery[hotRankQueryKey].(map[string]interface{})
	if !ok {
		return items, nil
	}

	hotRankID, ok := hotRankQuery["id"].(string)
	if !ok {
		return items, nil
	}

	// 获取热榜列表数据
	hotRankDataJSON, err := json.Marshal(apolloState.DefaultClient[hotRankID])
	if err != nil {
		return items, nil
	}

	var hotRankData KuaishouHotRankData
	if err := json.Unmarshal(hotRankDataJSON, &hotRankData); err != nil {
		return items, nil
	}

	// 转换数据格式
	for _, item := range hotRankData.Items {
		// 从id中提取实际的热搜词
		hotSearchWord := item.ID

		// 获取具体的热榜项数据
		hotItemJSON, err := json.Marshal(apolloState.DefaultClient[item.ID])
		if err != nil {
			continue
		}

		var hotItem KuaishouHotItem
		if err := json.Unmarshal(hotItemJSON, &hotItem); err != nil {
			continue
		}

		// 跳过置顶项
		if hotItem.TagType == "置顶" {
			continue
		}

		newsItem := types.NewsItem{
			ID:    hotSearchWord,
			Title: hotItem.Name,
			URL:   fmt.Sprintf("https://www.kuaishou.com/search/video?searchKey=%s", hotItem.Name),
		}
		items = append(items, newsItem)
	}

	return items, nil
}

type KaopuFetcher struct {
	fetcher.BaseFetcher
}

func NewKaopuFetcher() *KaopuFetcher {
	return &KaopuFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "kaopu"},
	}
}

func (f *KaopuFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://www.kaopu001.com/rss", false).Fetch()
}

type CankaoXiaoXiFetcher struct {
	fetcher.BaseFetcher
}

func NewCankaoXiaoXiFetcher() *CankaoXiaoXiFetcher {
	return &CankaoXiaoXiFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "cankaoxiaoxi"},
	}
}

func (f *CankaoXiaoXiFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://www.cankaoxiaoxi.com/rss", false).Fetch()
}

type SputnikNewsCNFetcher struct {
	fetcher.BaseFetcher
}

func NewSputnikNewsCNFetcher() *SputnikNewsCNFetcher {
	return &SputnikNewsCNFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "sputniknewscn"},
	}
}

func (f *SputnikNewsCNFetcher) Fetch() ([]types.NewsItem, error) {
	return fetcher.DefineRSSSource("https://sputniknews.cn/rss", false).Fetch()
}
