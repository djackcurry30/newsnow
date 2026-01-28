package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"newsnow-go/internal/python"

	"github.com/gin-gonic/gin"
)

// ExecutePythonRequest 执行Python请求
type ExecutePythonRequest struct {
	Code             string                 `json:"code" binding:"required"`
	Requirements     []string               `json:"requirements,omitempty"`
	RequirementsFile string                 `json:"requirements_file,omitempty"`
	Globals          map[string]interface{} `json:"globals,omitempty"`
	Timeout          int                    `json:"timeout,omitempty"` // 秒

	// 资源管理选项
	EnvID         string               `json:"env_id,omitempty"`         // 复用已有环境ID
	CleanupPolicy python.CleanupPolicy `json:"cleanup_policy,omitempty"` // 清理策略: auto, manual, ttl
	TTL           int                  `json:"ttl,omitempty"`            // TTL秒数（cleanup_policy=ttl时有效）
	ReuseExisting bool                 `json:"reuse_existing,omitempty"` // 是否复用匹配的现有环境
}

// ExecutePythonFileRequest 执行Python文件请求
type ExecutePythonFileRequest struct {
	Filepath string                 `json:"filepath" binding:"required"`
	Globals  map[string]interface{} `json:"globals,omitempty"`
}

// CleanupEnvironmentRequest 清理环境请求
type CleanupEnvironmentRequest struct {
	EnvID string `json:"env_id" binding:"required"`
}

// EnvironmentResponse 环境信息响应
type EnvironmentResponse struct {
	ID               string               `json:"id"`
	Path             string               `json:"path"`
	WorkingDir       string               `json:"working_dir"`
	Requirements     []string             `json:"requirements"`
	RequirementsFile string               `json:"requirements_file"`
	CreatedAt        time.Time            `json:"created_at"`
	LastUsedAt       time.Time            `json:"last_used_at"`
	CleanupPolicy    python.CleanupPolicy `json:"cleanup_policy"`
	TTL              time.Duration        `json:"ttl,omitempty"`
}

// ExecutePythonResponse 执行响应
type ExecutePythonResponse struct {
	Success bool        `json:"success"`
	Stdout  string      `json:"stdout,omitempty"`
	Stderr  string      `json:"stderr,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	EnvID   string      `json:"env_id,omitempty"` // 返回环境ID（如果是手动清理模式）
	Error   string      `json:"error,omitempty"`
}

// ListEnvironmentsResponse 环境列表响应
type ListEnvironmentsResponse struct {
	Success      bool                  `json:"success"`
	Environments []EnvironmentResponse `json:"environments,omitempty"`
	Count        int                   `json:"count"`
	Error        string                `json:"error,omitempty"`
}

// PythonHandler Python处理器
type PythonHandler struct {
	executor *python.PythonExecutor
}

// NewPythonHandler 创建新的Python处理器
func NewPythonHandler() (*PythonHandler, error) {
	executor, err := python.NewPythonExecutor()
	if err != nil {
		return nil, err
	}

	return &PythonHandler{
		executor: executor,
	}, nil
}

// Execute 执行 Python 代码字符串（支持虚拟环境和依赖安装）
func (h *PythonHandler) Execute(c *gin.Context) {
	var req ExecutePythonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ExecutePythonResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// 设置超时
	timeout := 30 * time.Second
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	// 设置TTL
	var ttl time.Duration
	if req.TTL > 0 {
		ttl = time.Duration(req.TTL) * time.Second
	}

	// 使用新的 ExecuteWithOptions 方法
	opts := python.ExecutionOptions{
		Code:             req.Code,
		Requirements:     req.Requirements,
		RequirementsFile: req.RequirementsFile,
		Globals:          req.Globals,
		Timeout:          timeout,
		EnvID:            req.EnvID,
		CleanupPolicy:    req.CleanupPolicy,
		TTL:              ttl,
		ReuseExisting:    req.ReuseExisting,
	}

	result, err := h.executor.ExecuteWithOptions(opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ExecutePythonResponse{
			Success: false,
			Stdout:  "",
			Stderr:  "",
			Error:   err.Error(),
		})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ExecutePythonResponse{
			Success: false,
			Stdout:  result.Stdout,
			Stderr:  result.Stderr,
			EnvID:   result.EnvID,
			Error:   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ExecutePythonResponse{
		Success: true,
		Stdout:  result.Stdout,
		Stderr:  result.Stderr,
		EnvID:   result.EnvID,
	})
}

// ExecuteWithGlobals 执行 Python 代码并传递全局变量
func (h *PythonHandler) ExecuteWithGlobals(c *gin.Context) {
	var req ExecutePythonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ExecutePythonResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// 设置超时
	timeout := 30 * time.Second
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	// 设置TTL
	var ttl time.Duration
	if req.TTL > 0 {
		ttl = time.Duration(req.TTL) * time.Second
	}

	// 使用新的 ExecuteWithOptions 方法
	opts := python.ExecutionOptions{
		Code:             req.Code,
		Requirements:     req.Requirements,
		RequirementsFile: req.RequirementsFile,
		Globals:          req.Globals,
		Timeout:          timeout,
		EnvID:            req.EnvID,
		CleanupPolicy:    req.CleanupPolicy,
		TTL:              ttl,
		ReuseExisting:    req.ReuseExisting,
	}

	result, err := h.executor.ExecuteWithOptions(opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ExecutePythonResponse{
			Success: false,
			Stdout:  "",
			Stderr:  "",
			Error:   err.Error(),
		})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ExecutePythonResponse{
			Success: false,
			Stdout:  result.Stdout,
			Stderr:  result.Stderr,
			EnvID:   result.EnvID,
			Error:   result.Error.Error(),
		})
		return
	}

	// 尝试解析结果
	var res interface{}
	if req.Globals != nil && len(req.Globals) > 0 {
		// 尝试解析 JSON 结果
		var jsonResult map[string]interface{}
		if err := json.Unmarshal([]byte(result.Stdout), &jsonResult); err == nil {
			if r, ok := jsonResult["result"]; ok {
				res = r
			}
		}
	}

	if res == nil {
		res = result.Stdout
	}

	c.JSON(http.StatusOK, ExecutePythonResponse{
		Success: true,
		Stdout:  result.Stdout,
		Stderr:  result.Stderr,
		Result:  res,
		EnvID:   result.EnvID,
	})
}

// ExecuteFile 执行 Python 文件
func (h *PythonHandler) ExecuteFile(c *gin.Context) {
	var req ExecutePythonFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ExecutePythonResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	stdout, stderr, err := h.executor.ExecuteFile(req.Filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ExecutePythonResponse{
			Success: false,
			Stdout:  stdout,
			Stderr:  stderr,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ExecutePythonResponse{
		Success: true,
		Stdout:  stdout,
		Stderr:  stderr,
	})
}

// ListEnvironments 列出所有虚拟环境
func (h *PythonHandler) ListEnvironments(c *gin.Context) {
	envs := h.executor.ListEnvironments()

	var responseEnvs []EnvironmentResponse
	for _, env := range envs {
		responseEnvs = append(responseEnvs, EnvironmentResponse{
			ID:               env.ID,
			Path:             env.Path,
			WorkingDir:       env.WorkingDir,
			Requirements:     env.Requirements,
			RequirementsFile: env.RequirementsFile,
			CreatedAt:        env.CreatedAt,
			LastUsedAt:       env.LastUsedAt,
			CleanupPolicy:    env.CleanupPolicy,
			TTL:              env.TTL,
		})
	}

	c.JSON(http.StatusOK, ListEnvironmentsResponse{
		Success:      true,
		Environments: responseEnvs,
		Count:        len(responseEnvs),
	})
}

// GetEnvironment 获取指定虚拟环境信息
func (h *PythonHandler) GetEnvironment(c *gin.Context) {
	envID := c.Param("id")
	if envID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "environment id is required",
		})
		return
	}

	env, err := h.executor.GetEnvironment(envID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"environment": EnvironmentResponse{
			ID:               env.ID,
			Path:             env.Path,
			WorkingDir:       env.WorkingDir,
			Requirements:     env.Requirements,
			RequirementsFile: env.RequirementsFile,
			CreatedAt:        env.CreatedAt,
			LastUsedAt:       env.LastUsedAt,
			CleanupPolicy:    env.CleanupPolicy,
			TTL:              env.TTL,
		},
	})
}

// CleanupEnvironment 手动清理指定虚拟环境
func (h *PythonHandler) CleanupEnvironment(c *gin.Context) {
	var req CleanupEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.executor.CleanupEnvironment(req.EnvID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "environment " + req.EnvID + " cleaned up successfully",
	})
}

// CleanupAllEnvironments 清理所有虚拟环境
func (h *PythonHandler) CleanupAllEnvironments(c *gin.Context) {
	if err := h.executor.CleanupAllEnvironments(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "all environments cleaned up successfully",
	})
}

// CleanupExpiredEnvironments 清理过期虚拟环境
func (h *PythonHandler) CleanupExpiredEnvironments(c *gin.Context) {
	if err := h.executor.CleanupExpiredEnvironments(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "expired environments cleaned up successfully",
	})
}

// Cleanup 清理 Python 解释器
func (h *PythonHandler) Cleanup() {
	python.Cleanup()
}
