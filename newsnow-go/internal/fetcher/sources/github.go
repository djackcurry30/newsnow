package sources

import (
	"newsnow-go/internal/fetcher"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GithubTrendingFetcher struct {
	fetcher.BaseFetcher
}

func NewGithubTrendingFetcher() *GithubTrendingFetcher {
	return &GithubTrendingFetcher{
		BaseFetcher: fetcher.BaseFetcher{ID: "github"},
	}
}

func (f *GithubTrendingFetcher) Fetch() ([]types.NewsItem, error) {
	url := "https://github.com/trending?spoken_language_code="
	html, err := utils.Fetch(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	items := make([]types.NewsItem, 0)
	baseURL := "https://github.com"

	doc.Find("main .Box div[data-hpc] > article").Each(func(i int, s *goquery.Selection) {
		titleLink := s.Find(">h2 a")
		title := strings.ReplaceAll(titleLink.Text(), "\n", "")
		title = strings.TrimSpace(title)
		
		href, exists := titleLink.Attr("href")
		if !exists || href == "" || title == "" {
			return
		}

		url := baseURL + href
		
		star := s.Find("[href$=stargazers]").Text()
		star = strings.ReplaceAll(star, "\n", "")
		star = strings.TrimSpace(star)
		
		desc := s.Find(">p").Text()
		desc = strings.ReplaceAll(desc, "\n", "")
		desc = strings.TrimSpace(desc)

		items = append(items, types.NewsItem{
			ID:   href,
			Title: title,
			URL:   url,
			Extra: &types.NewsExtra{
				Info:  "☆ " + star,
				Hover: desc,
			},
		})
	})

	return items, nil
}

type GithubHolder struct {
	Trending *GithubTrendingFetcher
}

func NewGithubHolder() *GithubHolder {
	trending := NewGithubTrendingFetcher()
	return &GithubHolder{
		Trending: trending,
	}
}

func (h *GithubHolder) GetFetchers() map[string]fetcher.SourceFetcher {
	return map[string]fetcher.SourceFetcher{
		"github":                  h.Trending,
		"github-trending-today": h.Trending,
	}
}
