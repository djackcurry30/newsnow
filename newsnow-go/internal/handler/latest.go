package handler

import (
	"newsnow-go/internal/config"
	"newsnow-go/internal/fetcher/sources"
	"newsnow-go/internal/service"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type LatestHandler struct {
	cfg     *config.Config
	cache   *service.CacheService
	sources *sources.SourceHolder
}

func NewLatestHandler(cfg *config.Config, cache *service.CacheService, sourceHolder *sources.SourceHolder) *LatestHandler {
	return &LatestHandler{
		cfg:     cfg,
		cache:   cache,
		sources: sourceHolder,
	}
}

func (h *LatestHandler) Handle(c *gin.Context) {
	sourceID := types.SourceID(c.Param("id"))
	source := types.SourceMap[sourceID]

	if source == nil {
		c.JSON(404, gin.H{"message": "Source not found"})
		return
	}

	latest := c.Query("latest") == "true"

	cache, _ := h.cache.Get(sourceID)
	disabledLogin := !h.cfg.IsLoginEnabled()

	hasUser := false
	if user, exists := c.Get("user"); exists {
		hasUser = user != nil
	}

	useCache, isSuccess := h.cache.ShouldUseCache(cache, source, latest, disabledLogin, hasUser)

	if useCache {
		if isSuccess {
			c.Header("X-Cache", "HIT")
		} else {
			c.Header("X-Cache", "STALE")
		}
		c.JSON(200, types.SourceResponse{
			Status: "success",
			ID:     sourceID,
			Items:  cache.Items,
		})
		return
	}

	f := h.sources.Get(string(sourceID))
	if f == nil {
		c.JSON(404, gin.H{"message": "Fetcher not found"})
		return
	}

	items, err := f.Fetch()
	if err != nil {
		utils.Error("Failed to fetch: " + err.Error())

		if cache != nil {
			c.Header("X-Cache", "STALE")
			c.JSON(200, types.SourceResponse{
				Status: "success",
				ID:     sourceID,
				Items:  cache.Items,
			})
			return
		}

		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	h.cache.UpdateCache(sourceID, items)

	c.Header("X-Cache", "MISS")
	c.JSON(200, types.SourceResponse{
		Status: "success",
		ID:     sourceID,
		Items:  items,
	})
}

type AllHandler struct {
	cfg     *config.Config
	cache   *service.CacheService
	sources *sources.SourceHolder
}

func NewAllHandler(cfg *config.Config, cache *service.CacheService, sourceHolder *sources.SourceHolder) *AllHandler {
	return &AllHandler{
		cfg:     cfg,
		cache:   cache,
		sources: sourceHolder,
	}
}

func (h *AllHandler) Handle(c *gin.Context) {
	rawIDs := c.Query("id")
	if rawIDs == "" {
		c.JSON(400, gin.H{"message": "Missing id parameter"})
		return
	}

	idList := strings.Split(rawIDs, ",")
	if len(idList) > 60 {
		c.JSON(400, gin.H{"message": "Too many IDs"})
		return
	}

	sourceIDs := make([]types.SourceID, 0, len(idList))
	for _, id := range idList {
		sourceID := types.SourceID(id)
		if types.SourceMap[sourceID] == nil {
			continue
		}
		sourceIDs = append(sourceIDs, sourceID)
	}

	if len(sourceIDs) == 0 {
		c.JSON(200, []types.CacheInfo{})
		return
	}

	caches, _ := h.cache.GetMultiple(sourceIDs)

	disabledLogin := !h.cfg.IsLoginEnabled()
	hasUser := false
	if user, exists := c.Get("user"); exists {
		hasUser = user != nil
	}

	cacheMap := make(map[types.SourceID]*types.CacheInfo)
	for i, cache := range caches {
		cacheMap[cache.ID] = &caches[i]
	}

	result := make([]types.CacheInfo, 0)
	for _, sourceID := range sourceIDs {
		cache, ok := cacheMap[sourceID]
		if !ok {
			cache = nil
		}

		source := types.SourceMap[sourceID]
		latest := true
		useCache, _ := h.cache.ShouldUseCache(cache, source, latest, disabledLogin, hasUser)

		if !useCache {
			f := h.sources.Get(string(sourceID))
			if f != nil {
				items, err := f.Fetch()
				if err == nil {
					h.cache.UpdateCache(sourceID, items)
					result = append(result, types.CacheInfo{
						ID:    sourceID,
						Items: items,
					})
					continue
				}
			}
		}

		if cache != nil {
			result = append(result, *cache)
		}
	}

	c.JSON(200, result)
}

type SHandler struct {
	cfg     *config.Config
	cache   *service.CacheService
	sources *sources.SourceHolder
}

func NewSHandler(cfg *config.Config, cache *service.CacheService, sourceHolder *sources.SourceHolder) *SHandler {
	return &SHandler{
		cfg:     cfg,
		cache:   cache,
		sources: sourceHolder,
	}
}

func (h *SHandler) Handle(c *gin.Context) {
	s := c.Param("s")
	parts := strings.SplitN(s, ",", 2)

	sourceID := types.SourceID(parts[0])
	source := types.SourceMap[sourceID]

	if source == nil {
		c.JSON(404, gin.H{"message": "Source not found"})
		return
	}

	cache, _ := h.cache.Get(sourceID)
	if cache != nil && len(cache.Items) > 0 {
		c.JSON(200, types.SourceResponse{
			Status: "success",
			ID:     sourceID,
			Items:  cache.Items,
		})
		return
	}

	f := h.sources.Get(string(sourceID))
	if f == nil {
		c.JSON(404, gin.H{"message": "Fetcher not found"})
		return
	}

	items, err := f.Fetch()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	h.cache.UpdateCache(sourceID, items)

	c.JSON(200, types.SourceResponse{
		Status: "success",
		ID:     sourceID,
		Items:  items,
	})
}

type ProxyHandler struct {
	cfg     *config.Config
	cache   *service.CacheService
	sources *sources.SourceHolder
}

func NewProxyHandler(cfg *config.Config, cache *service.CacheService, sourceHolder *sources.SourceHolder) *ProxyHandler {
	return &ProxyHandler{
		cfg:     cfg,
		cache:   cache,
		sources: sourceHolder,
	}
}

func (h *ProxyHandler) Handle(c *gin.Context) {
	s := c.Param("s")
	parts := strings.SplitN(s, ",", 2)

	sourceID := types.SourceID(parts[0])
	source := types.SourceMap[sourceID]

	if source == nil {
		c.JSON(404, gin.H{"message": "Source not found"})
		return
	}

	cache, _ := h.cache.Get(sourceID)
	if cache != nil && len(cache.Items) > 0 {
		c.JSON(200, types.SourceResponse{
			Status: "success",
			ID:     sourceID,
			Items:  cache.Items,
		})
		return
	}

	f := h.sources.Get(string(sourceID))
	if f == nil {
		c.JSON(404, gin.H{"message": "Fetcher not found"})
		return
	}

	items, err := f.Fetch()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	h.cache.UpdateCache(sourceID, items)

	c.JSON(200, types.SourceResponse{
		Status: "success",
		ID:     sourceID,
		Items:  items,
	})
}

type ConfigHandler struct {
	cfg *config.Config
}

func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

func (h *ConfigHandler) Handle(c *gin.Context) {
	c.JSON(200, gin.H{
		"server":      true,
		"login":       h.cfg.IsLoginEnabled(),
		"enableLogin": h.cfg.IsLoginEnabled(),
	})
}

type VersionHandler struct{}

func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

func (h *VersionHandler) Handle(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": "2.0.0-go",
	})
}

type EnableLoginHandler struct {
	cfg *config.Config
}

func NewEnableLoginHandler(cfg *config.Config) *EnableLoginHandler {
	return &EnableLoginHandler{cfg: cfg}
}

func (h *EnableLoginHandler) Handle(c *gin.Context) {
	c.JSON(200, gin.H{
		"enableLogin": h.cfg.IsLoginEnabled(),
	})
}

type UserAgentHandler struct{}

func NewUserAgentHandler() *UserAgentHandler {
	return &UserAgentHandler{}
}

func (h *UserAgentHandler) Handle(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.String(200, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
}
