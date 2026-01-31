package handler

import (
	"encoding/json"
	"net/http"

	"newsnow-go/internal/config"
	"newsnow-go/internal/service"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type SyncHandler struct {
	cfg         *config.Config
	userService *service.UserService
}

func NewSyncHandler(cfg *config.Config, userService *service.UserService) *SyncHandler {
	return &SyncHandler{
		cfg:         cfg,
		userService: userService,
	}
}

func (h *SyncHandler) Handle(c *gin.Context) {
	disabledLogin := !h.cfg.IsLoginEnabled()
	if disabledLogin {
		c.JSON(http.StatusUpgradeRequired, gin.H{
			"message": "Server not configured, disable login",
		})
		return
	}

	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userContext := userValue.(UserContext)
	userID := userContext.ID

	switch c.Request.Method {
	case "GET":
		h.handleGet(c, userID)
	case "POST":
		h.handlePost(c, userID)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Method not allowed"})
	}
}

func (h *SyncHandler) handleGet(c *gin.Context, userID string) {
	result, err := h.userService.GetUserData(userID)
	if err != nil {
		utils.Error("Failed to get user data: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user data"})
		return
	}

	if result == nil {
		c.JSON(http.StatusOK, gin.H{
			"data":        nil,
			"updatedTime": nil,
		})
		return
	}

	var data interface{}
	if len(result.Data) > 0 {
		if err := json.Unmarshal(result.Data, &data); err != nil {
			data = nil
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        data,
		"updatedTime": result.Updated,
	})
}

func (h *SyncHandler) handlePost(c *gin.Context, userID string) {
	var req struct {
		UpdatedTime int64                  `json:"updatedTime"`
		Data        map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if err := h.verifyMetadata(req.Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid metadata: " + err.Error()})
		return
	}

	if err := h.userService.SetUserData(userID, req.Data, req.UpdatedTime); err != nil {
		utils.Error("Failed to set user data: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"updatedTime": req.UpdatedTime,
	})
}

func (h *SyncHandler) verifyMetadata(data map[string]interface{}) error {
	if data == nil {
		return nil
	}

	if columns, ok := data["columns"].([]interface{}); ok {
		for _, col := range columns {
			if colMap, ok := col.(map[string]interface{}); ok {
				if id, ok := colMap["id"].(string); ok {
					if !isValidColumnID(id) {
						return nil
					}
				}
				if sources, ok := colMap["sources"].([]interface{}); ok {
					for _, src := range sources {
						if srcID, ok := src.(string); ok {
							if !isValidSourceID(srcID) {
								return nil
							}
						}
					}
				}
			}
		}
	}

	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		if action, ok := metadata["action"].(string); ok {
			if action != "" && action != "add" && action != "delete" && action != "update" {
				return nil
			}
		}
	}

	return nil
}

func isValidColumnID(id string) bool {
	validColumns := []string{"china", "tech", "finance", "world", "video", "forum", "social"}
	for _, col := range validColumns {
		if col == id {
			return true
		}
	}
	return false
}

func isValidSourceID(id string) bool {
	return types.SourceMap[types.SourceID(id)] != nil
}
