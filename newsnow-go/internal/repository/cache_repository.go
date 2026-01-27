package repository

import (
	"encoding/json"
	"newsnow-go/internal/model"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
	"time"

	"gorm.io/gorm"
)

type CacheRepository struct {
	db *gorm.DB
}

func NewCacheRepository(db *gorm.DB) *CacheRepository {
	return &CacheRepository{db: db}
}

func (r *CacheRepository) Init() error {
	return r.db.AutoMigrate(&model.Cache{})
}

func (r *CacheRepository) Set(key string, value []types.NewsItem) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	cache := model.Cache{
		ID:      key,
		Updated: now,
		Data:    string(data),
	}

	return r.db.Save(&cache).Error
}

func (r *CacheRepository) Get(key string) (*types.CacheInfo, error) {
	var cache model.Cache
	err := r.db.Where("id = ?", key).First(&cache).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var items []types.NewsItem
	if err := json.Unmarshal([]byte(cache.Data), &items); err != nil {
		return nil, err
	}

	return &types.CacheInfo{
		ID:      types.SourceID(cache.ID),
		Updated: cache.Updated,
		Items:   items,
	}, nil
}

func (r *CacheRepository) GetEntire(keys []string) ([]types.CacheInfo, error) {
	if len(keys) == 0 {
		return []types.CacheInfo{}, nil
	}

	var caches []model.Cache
	placeholders := ""
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		if i > 0 {
			placeholders += " OR "
		}
		placeholders += "id = ?"
		args[i] = key
	}

	err := r.db.Where(placeholders, args...).Find(&caches).Error
	if err != nil {
		return nil, err
	}

	result := make([]types.CacheInfo, 0, len(caches))
	for _, cache := range caches {
		var items []types.NewsItem
		if err := json.Unmarshal([]byte(cache.Data), &items); err != nil {
			utils.Warn("Failed to parse cache data for " + cache.ID)
			continue
		}
		result = append(result, types.CacheInfo{
			ID:      types.SourceID(cache.ID),
			Updated: cache.Updated,
			Items:   items,
		})
	}

	return result, nil
}

func (r *CacheRepository) Delete(key string) error {
	return r.db.Where("id = ?", key).Delete(&model.Cache{}).Error
}

func (r *CacheRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
