package sources

import (
	"newsnow-go/internal/fetcher"
)

type SourceHolder struct {
	fetchers map[string]fetcher.SourceFetcher
}

func NewSourceHolder() *SourceHolder {
	h := &SourceHolder{
		fetchers: make(map[string]fetcher.SourceFetcher),
	}

	v2exHolder := NewV2EXHolder()
	for id, f := range v2exHolder.GetFetchers() {
		h.fetchers[id] = f
	}

	githubHolder := NewGithubHolder()
	for id, f := range githubHolder.GetFetchers() {
		h.fetchers[id] = f
	}

	kr36Holder := NewKr36Holder()
	for id, f := range kr36Holder.GetFetchers() {
		h.fetchers[id] = f
	}

	bilibiliHolder := NewBilibiliHolder()
	for id, f := range bilibiliHolder.GetFetchers() {
		h.fetchers[id] = f
	}

	wallstreetcnHolder := NewWallStreetCNHolder()
	for id, f := range wallstreetcnHolder.GetFetchers() {
		h.fetchers[id] = f
	}

	h.fetchers["zhihu"] = NewZhihuFetcher()
	h.fetchers["weibo"] = NewWeiboFetcher()
	h.fetchers["weibo-hot"] = NewWeiboHotFetcher()
	h.fetchers["jin10"] = NewJin10Fetcher()
	h.fetchers["coolapk"] = NewCoolapkFetcher()
	h.fetchers["ithome"] = NewIthomeFetcher()
	h.fetchers["sspai"] = NewSspaiFetcher()
	h.fetchers["juejin"] = NewJuejinFetcher()
	h.fetchers["solidot"] = NewSolidotFetcher()
	h.fetchers["hackernews"] = NewHackerNewsFetcher()
	h.fetchers["producthunt"] = NewProductHuntFetcher()
	h.fetchers["baidu"] = NewBaiduFetcher()
	h.fetchers["douyin"] = NewDouyinFetcher()
	h.fetchers["hupu"] = NewHupuFetcher()
	h.fetchers["tieba"] = NewTiebaFetcher()
	h.fetchers["toutiao"] = NewToutiaoFetcher()
	h.fetchers["thepaper"] = NewThePaperFetcher()
	h.fetchers["zaobao"] = NewZaobaoFetcher()
	clsHolder := NewCLSHolder()
	for id, f := range clsHolder.GetFetchers() {
		h.fetchers[id] = f
	}
	xueqiuHolder := NewXueqiuHolder()
	for id, f := range xueqiuHolder.GetFetchers() {
		h.fetchers[id] = f
	}
	h.fetchers["gelonghui"] = NewGeLongHuiFetcher()
	fastbullHolder := NewFastBullHolder()
	for id, f := range fastbullHolder.GetFetchers() {
		h.fetchers[id] = f
	}
	mktnewsHolder := NewMKTNewsHolder()
	for id, f := range mktnewsHolder.GetFetchers() {
		h.fetchers[id] = f
	}
	pcbetaHolder := NewPCBetaHolder()
	for id, f := range pcbetaHolder.GetFetchers() {
		h.fetchers[id] = f
	}
	h.fetchers["freebuf"] = NewFreebufFetcher()
	h.fetchers["steam"] = NewSteamFetcher()
	tencentHolder := NewTencentHolder()
	for id, f := range tencentHolder.GetFetchers() {
		h.fetchers[id] = f
	}
	qqvideoHolder := NewQQVideoHolder()
	for id, f := range qqvideoHolder.GetFetchers() {
		h.fetchers[id] = f
	}
	iqiyiHolder := NewIqiyiHolder()
	for id, f := range iqiyiHolder.GetFetchers() {
		h.fetchers[id] = f
	}
	h.fetchers["douban"] = NewDoubanFetcher()
	h.fetchers["nowcoder"] = NewNowCoderFetcher()
	h.fetchers["ifeng"] = NewIfengFetcher()
	chongbuluoHolder := NewChongBuluoHolder()
	for id, f := range chongbuluoHolder.GetFetchers() {
		h.fetchers[id] = f
	}
	h.fetchers["kuaishou"] = NewKuaishouFetcher()
	h.fetchers["kaopu"] = NewKaopuFetcher()
	h.fetchers["cankaoxiaoxi"] = NewCankaoXiaoXiFetcher()
	h.fetchers["sputniknewscn"] = NewSputnikNewsCNFetcher()

	return h
}

func (h *SourceHolder) Get(id string) fetcher.SourceFetcher {
	return h.fetchers[id]
}

func (h *SourceHolder) All() map[string]fetcher.SourceFetcher {
	return h.fetchers
}
