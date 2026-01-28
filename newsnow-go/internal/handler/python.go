package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"newsnow-go/internal/python"

	"github.com/gin-gonic/gin"
)

type ExecutePythonRequest struct {
	Code             string                 `json:"code" binding:"required"`
	Requirements     []string               `json:"requirements,omitempty"`
	RequirementsFile string                 `json:"requirements_file,omitempty"` // requirements.txt 文件路径
	Globals          map[string]interface{} `json:"globals,omitempty"`
	Timeout          int                    `json:"timeout,omitempty"` // 秒
}

type ExecutePythonFileRequest struct {
	Filepath string                 `json:"filepath" binding:"required"`
	Globals  map[string]interface{} `json:"globals,omitempty"`
}

type ExecutePythonResponse struct {
	Success bool        `json:"success"`
	Stdout  string      `json:"stdout,omitempty"`
	Stderr  string      `json:"stderr,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PythonHandler struct {
	executor *python.PythonExecutor
}

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

	// 使用新的 ExecuteWithOptions 方法
	opts := python.ExecutionOptions{
		Code:             req.Code,
		Requirements:     req.Requirements,
		RequirementsFile: req.RequirementsFile,
		Globals:          req.Globals,
		Timeout:          timeout,
	}

	stdout, stderr, err := h.executor.ExecuteWithOptions(opts)
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

	// 使用新的 ExecuteWithOptions 方法
	opts := python.ExecutionOptions{
		Code:             req.Code,
		Requirements:     req.Requirements,
		RequirementsFile: req.RequirementsFile,
		Globals:          req.Globals,
		Timeout:          timeout,
	}

	stdout, stderr, err := h.executor.ExecuteWithOptions(opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ExecutePythonResponse{
			Success: false,
			Stdout:  stdout,
			Stderr:  stderr,
			Error:   err.Error(),
		})
		return
	}

	// 尝试解析结果
	var result interface{}
	if req.Globals != nil && len(req.Globals) > 0 {
		// 尝试解析 JSON 结果
		var jsonResult map[string]interface{}
		if err := json.Unmarshal([]byte(stdout), &jsonResult); err == nil {
			if res, ok := jsonResult["result"]; ok {
				result = res
			}
		}
	}

	if result == nil {
		result = stdout
	}

	c.JSON(http.StatusOK, ExecutePythonResponse{
		Success: true,
		Stdout:  stdout,
		Stderr:  stderr,
		Result:  result,
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

// Cleanup 清理 Python 解释器
func (h *PythonHandler) Cleanup() {
	python.Cleanup()
}
