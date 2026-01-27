package service

import (
	"newsnow-go/internal/repository"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
	"time"
)

type CacheService struct {
	repo *repository.CacheRepository
	TTL  time.Duration
}

func NewCacheService(repo *repository.CacheRepository) *CacheService {
	return &CacheService{
		repo: repo,
		TTL:  30 * time.Minute,
	}
}

func (s *CacheService) Get(sourceID types.SourceID) (*types.CacheInfo, error) {
	if s.repo == nil {
		return nil, nil
	}
	return s.repo.Get(string(sourceID))
}

func (s *CacheService) Set(sourceID types.SourceID, items []types.NewsItem) error {
	if s.repo == nil {
		return nil
	}
	return s.repo.Set(string(sourceID), items)
}

func (s *CacheService) GetMultiple(sourceIDs []types.SourceID) ([]types.CacheInfo, error) {
	if s.repo == nil {
		return []types.CacheInfo{}, nil
	}
	
	keys := make([]string, len(sourceIDs))
	for i, id := range sourceIDs {
		keys[i] = string(id)
	}
	return s.repo.GetEntire(keys)
}

func (s *CacheService) ShouldUseCache(cache *types.CacheInfo, source *types.Source, latest bool, disabledLogin bool, hasUser bool) (useCache bool, isSuccess bool) {
	if cache == nil {
		return false, false
	}

	now := time.Now().UnixMilli()
	interval := int64(source.Interval)
	ttl := int64(s.TTL.Milliseconds())

	if now-cache.Updated < interval {
		return true, true
	}

	if now-cache.Updated < ttl {
		if !latest || (disabledLogin || !hasUser) {
			return true, false
		}
	}

	return false, false
}

func (s *CacheService) UpdateCache(sourceID types.SourceID, items []types.NewsItem) {
	if s.repo == nil {
		return
	}
	
	go func() {
		if err := s.repo.Set(string(sourceID), items); err != nil {
			utils.Error("Failed to update cache: " + err.Error())
		}
	}()
}
