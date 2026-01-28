package python

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	pythonPath string
	once       sync.Once
)

// getPythonPath 获取 Python 可执行文件路径
func getPythonPath() string {
	once.Do(func() {
		// 尝试查找 python3 或 python
		for _, cmd := range []string{"python3", "python"} {
			if path, err := exec.LookPath(cmd); err == nil {
				pythonPath = path
				return
			}
		}
		// 默认使用 python3
		pythonPath = "python3"
	})
	return pythonPath
}

// CleanupPolicy 清理策略
type CleanupPolicy string

const (
	// CleanupPolicyAuto 执行后自动清理
	CleanupPolicyAuto CleanupPolicy = "auto"
	// CleanupPolicyManual 手动清理
	CleanupPolicyManual CleanupPolicy = "manual"
	// CleanupPolicyTTL 按TTL自动清理
	CleanupPolicyTTL CleanupPolicy = "ttl"
)

// VirtualEnv 虚拟环境信息
type VirtualEnv struct {
	ID               string        `json:"id"`
	Path             string        `json:"path"`
	WorkingDir       string        `json:"working_dir"`
	Requirements     []string      `json:"requirements"`
	RequirementsFile string        `json:"requirements_file"`
	CreatedAt        time.Time     `json:"created_at"`
	LastUsedAt       time.Time     `json:"last_used_at"`
	CleanupPolicy    CleanupPolicy `json:"cleanup_policy"`
	TTL              time.Duration `json:"ttl,omitempty"`
	mu               sync.RWMutex
}

// IsExpired 检查虚拟环境是否过期
func (v *VirtualEnv) IsExpired() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.CleanupPolicy == CleanupPolicyManual {
		return false
	}
	if v.CleanupPolicy == CleanupPolicyAuto {
		return true // 标记为需要清理
	}
	if v.CleanupPolicy == CleanupPolicyTTL && v.TTL > 0 {
		return time.Since(v.LastUsedAt) > v.TTL
	}
	return false
}

// UpdateLastUsed 更新最后使用时间
func (v *VirtualEnv) UpdateLastUsed() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.LastUsedAt = time.Now()
}

// ExecutionOptions 执行选项
type ExecutionOptions struct {
	Code             string                 `json:"code"`
	Requirements     []string               `json:"requirements,omitempty"`
	RequirementsFile string                 `json:"requirements_file,omitempty"`
	Globals          map[string]interface{} `json:"globals,omitempty"`
	Timeout          time.Duration          `json:"timeout,omitempty"`

	// 新增：资源管理选项
	EnvID         string        `json:"env_id,omitempty"`         // 复用已有环境ID
	CleanupPolicy CleanupPolicy `json:"cleanup_policy,omitempty"` // 清理策略
	TTL           time.Duration `json:"ttl,omitempty"`            // TTL清理模式的过期时间
	ReuseExisting bool          `json:"reuse_existing,omitempty"` // 是否复用匹配的现有环境
}

// ExecutionResult 执行结果
type ExecutionResult struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	EnvID  string `json:"env_id,omitempty"` // 返回环境ID（如果是手动清理模式）
	Error  error  `json:"-"`
}

// PythonExecutor Python执行器
type PythonExecutor struct {
	mu           sync.RWMutex
	timeout      time.Duration
	workingDir   string
	environments map[string]*VirtualEnv // 管理持久化的虚拟环境
	envCounter   int64
}

// NewPythonExecutor 创建一个新的 Python 执行器
func NewPythonExecutor() (*PythonExecutor, error) {
	executor := &PythonExecutor{
		timeout:      30 * time.Second,
		workingDir:   "",
		environments: make(map[string]*VirtualEnv),
	}

	// 启动自动清理协程
	go executor.autoCleanupLoop()

	return executor, nil
}

// SetTimeout 设置执行超时时间
func (e *PythonExecutor) SetTimeout(timeout time.Duration) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.timeout = timeout
}

// SetWorkingDir 设置工作目录
func (e *PythonExecutor) SetWorkingDir(dir string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.workingDir = dir
}

// autoCleanupLoop 自动清理循环
func (e *PythonExecutor) autoCleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		e.CleanupExpiredEnvironments()
	}
}

// CleanupExpiredEnvironments 清理过期的虚拟环境
func (e *PythonExecutor) CleanupExpiredEnvironments() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var toDelete []string
	for id, env := range e.environments {
		if env.IsExpired() {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		env := e.environments[id]
		delete(e.environments, id)

		// 异步删除目录
		go func(path string) {
			if err := os.RemoveAll(path); err != nil {
				fmt.Printf("Failed to cleanup environment at %s: %v\n", path, err)
			}
		}(env.Path)
	}

	return nil
}

// CleanupEnvironment 手动清理指定虚拟环境
func (e *PythonExecutor) CleanupEnvironment(envID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	env, exists := e.environments[envID]
	if !exists {
		return fmt.Errorf("environment %s not found", envID)
	}

	delete(e.environments, envID)

	// 删除虚拟环境目录
	if err := os.RemoveAll(env.Path); err != nil {
		return fmt.Errorf("failed to cleanup environment: %w", err)
	}

	return nil
}

// CleanupAllEnvironments 清理所有虚拟环境
func (e *PythonExecutor) CleanupAllEnvironments() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var lastErr error
	for id, env := range e.environments {
		delete(e.environments, id)
		if err := os.RemoveAll(env.Path); err != nil {
			lastErr = err
			fmt.Printf("Failed to cleanup environment %s: %v\n", id, err)
		}
	}

	return lastErr
}

// GetEnvironment 获取虚拟环境信息
func (e *PythonExecutor) GetEnvironment(envID string) (*VirtualEnv, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	env, exists := e.environments[envID]
	if !exists {
		return nil, fmt.Errorf("environment %s not found", envID)
	}

	return env, nil
}

// ListEnvironments 列出所有虚拟环境
func (e *PythonExecutor) ListEnvironments() []*VirtualEnv {
	e.mu.RLock()
	defer e.mu.RUnlock()

	envs := make([]*VirtualEnv, 0, len(e.environments))
	for _, env := range e.environments {
		envs = append(envs, env)
	}

	return envs
}

// findReusableEnvironment 查找可复用的虚拟环境
func (e *PythonExecutor) findReusableEnvironment(requirements []string, requirementsFile string) *VirtualEnv {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, env := range e.environments {
		// 检查依赖是否匹配
		if requirementsFile != "" && env.RequirementsFile == requirementsFile {
			return env
		}
		if len(requirements) > 0 && sliceEqual(env.Requirements, requirements) {
			return env
		}
	}

	return nil
}

// sliceEqual 检查两个字符串切片是否相等
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// createVirtualEnv 创建虚拟环境
func (e *PythonExecutor) createVirtualEnv(venvPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, getPythonPath(), "-m", "venv", venvPath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("virtual environment creation timeout")
	}
	if err != nil {
		return fmt.Errorf("failed to create virtual environment: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

// installRequirements 在虚拟环境中安装依赖
func (e *PythonExecutor) installRequirements(venvPath string, requirements []string) error {
	if len(requirements) == 0 {
		return nil
	}

	pipPath := filepath.Join(venvPath, "bin", "pip")
	if _, err := os.Stat(pipPath); os.IsNotExist(err) {
		pipPath = filepath.Join(venvPath, "Scripts", "pip.exe") // Windows
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	args := append([]string{"install"}, requirements...)
	cmd := exec.CommandContext(ctx, pipPath, args...)

	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("package installation timeout")
	}
	if err != nil {
		return fmt.Errorf("failed to install packages: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

// installRequirementsFromFile 从 requirements.txt 文件安装依赖
func (e *PythonExecutor) installRequirementsFromFile(venvPath string, requirementsFile string) error {
	if requirementsFile == "" {
		return nil
	}

	// 检查文件是否存在
	if _, err := os.Stat(requirementsFile); os.IsNotExist(err) {
		return fmt.Errorf("requirements file not found: %s", requirementsFile)
	}

	pipPath := filepath.Join(venvPath, "bin", "pip")
	if _, err := os.Stat(pipPath); os.IsNotExist(err) {
		pipPath = filepath.Join(venvPath, "Scripts", "pip.exe") // Windows
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, pipPath, "install", "-r", requirementsFile)

	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("package installation from file timeout")
	}
	if err != nil {
		return fmt.Errorf("failed to install packages from file: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

// getVenvPythonPath 获取虚拟环境中的 Python 路径
func (e *PythonExecutor) getVenvPythonPath(venvPath string) string {
	pythonPath := filepath.Join(venvPath, "bin", "python")
	if _, err := os.Stat(pythonPath); os.IsNotExist(err) {
		pythonPath = filepath.Join(venvPath, "Scripts", "python.exe") // Windows
	}
	return pythonPath
}

// generateEnvID 生成环境ID
func (e *PythonExecutor) generateEnvID() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.envCounter++
	return fmt.Sprintf("env_%d_%d", time.Now().Unix(), e.envCounter)
}

// ExecuteWithOptions 使用选项执行 Python 代码（支持虚拟环境和依赖安装）
func (e *PythonExecutor) ExecuteWithOptions(opts ExecutionOptions) (*ExecutionResult, error) {
	timeout := e.timeout
	if opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	// 确定清理策略
	cleanupPolicy := opts.CleanupPolicy
	if cleanupPolicy == "" {
		cleanupPolicy = CleanupPolicyAuto // 默认自动清理
	}

	var pythonExec string
	var envID string
	var env *VirtualEnv

	// 检查是否需要虚拟环境
	needVirtualEnv := len(opts.Requirements) > 0 || opts.RequirementsFile != ""

	if needVirtualEnv {
		// 检查是否复用已有环境
		if opts.EnvID != "" {
			// 使用指定的环境ID
			existingEnv, err := e.GetEnvironment(opts.EnvID)
			if err != nil {
				return nil, fmt.Errorf("failed to get environment %s: %w", opts.EnvID, err)
			}
			env = existingEnv
			env.UpdateLastUsed()
			pythonExec = e.getVenvPythonPath(env.Path)
			envID = env.ID
		} else if opts.ReuseExisting {
			// 查找可复用的环境
			existingEnv := e.findReusableEnvironment(opts.Requirements, opts.RequirementsFile)
			if existingEnv != nil {
				env = existingEnv
				env.UpdateLastUsed()
				pythonExec = e.getVenvPythonPath(env.Path)
				envID = env.ID
			}
		}

		// 如果没有找到可复用的环境，创建新环境
		if env == nil {
			// 创建临时目录作为虚拟环境
			venvDir, err := os.MkdirTemp("", "python_venv_*")
			if err != nil {
				return nil, fmt.Errorf("failed to create temp directory for venv: %w", err)
			}

			venvPath := filepath.Join(venvDir, "venv")

			// 创建虚拟环境
			if err := e.createVirtualEnv(venvPath); err != nil {
				os.RemoveAll(venvDir)
				return nil, err
			}

			// 安装依赖（列表形式）
			if len(opts.Requirements) > 0 {
				if err := e.installRequirements(venvPath, opts.Requirements); err != nil {
					os.RemoveAll(venvDir)
					return nil, err
				}
			}

			// 从 requirements.txt 文件安装依赖
			if opts.RequirementsFile != "" {
				if err := e.installRequirementsFromFile(venvPath, opts.RequirementsFile); err != nil {
					os.RemoveAll(venvDir)
					return nil, err
				}
			}

			pythonExec = e.getVenvPythonPath(venvPath)

			// 根据清理策略决定是否保存环境
			if cleanupPolicy == CleanupPolicyManual || cleanupPolicy == CleanupPolicyTTL {
				envID = e.generateEnvID()
				env = &VirtualEnv{
					ID:               envID,
					Path:             venvPath,
					WorkingDir:       e.workingDir,
					Requirements:     opts.Requirements,
					RequirementsFile: opts.RequirementsFile,
					CreatedAt:        time.Now(),
					LastUsedAt:       time.Now(),
					CleanupPolicy:    cleanupPolicy,
					TTL:              opts.TTL,
				}

				e.mu.Lock()
				e.environments[envID] = env
				e.mu.Unlock()
			} else {
				// 自动清理模式：延迟删除
				defer os.RemoveAll(venvDir)
			}
		}
	} else {
		pythonExec = getPythonPath()
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "python_script_*.py")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// 构建代码
	var codeBuilder strings.Builder

	// 如果有全局变量，添加它们
	if len(opts.Globals) > 0 {
		codeBuilder.WriteString("import json\n")
		for key, value := range opts.Globals {
			switch v := value.(type) {
			case string:
				codeBuilder.WriteString(fmt.Sprintf("%s = %q\n", key, v))
			case int, int64, float64, bool:
				codeBuilder.WriteString(fmt.Sprintf("%s = %v\n", key, v))
			case []interface{}, map[string]interface{}:
				jsonBytes, _ := json.Marshal(v)
				codeBuilder.WriteString(fmt.Sprintf("%s = json.loads(%q)\n", key, string(jsonBytes)))
			default:
				codeBuilder.WriteString(fmt.Sprintf("%s = %q\n", key, fmt.Sprintf("%v", v)))
			}
		}
		codeBuilder.WriteString("\n")
	}

	codeBuilder.WriteString(opts.Code)

	// 写入代码
	if _, err := tmpFile.WriteString(codeBuilder.String()); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to write code to temp file: %w", err)
	}
	tmpFile.Close()

	// 执行代码
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, pythonExec, tmpFile.Name())

	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return &ExecutionResult{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
			EnvID:  envID,
			Error:  fmt.Errorf("execution timeout"),
		}, nil
	}

	if err != nil {
		return &ExecutionResult{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
			EnvID:  envID,
			Error:  fmt.Errorf("execution failed: %w", err),
		}, nil
	}

	return &ExecutionResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		EnvID:  envID,
		Error:  nil,
	}, nil
}

// Execute 执行 Python 代码字符串（向后兼容）
func (e *PythonExecutor) Execute(code string) (string, string, error) {
	result, err := e.ExecuteWithOptions(ExecutionOptions{
		Code: code,
	})
	if err != nil {
		return "", "", err
	}
	return result.Stdout, result.Stderr, result.Error
}

// ExecuteWithInput 执行 Python 代码并传递输入数据
func (e *PythonExecutor) ExecuteWithInput(code string, input string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, getPythonPath(), "-c", code)

	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}

	// 设置 stdin
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return stdout.String(), stderr.String(), fmt.Errorf("execution timeout")
	}

	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("execution failed: %w", err)
	}

	return stdout.String(), stderr.String(), nil
}

// ExecuteWithGlobals 执行 Python 代码并传递全局变量（向后兼容）
func (e *PythonExecutor) ExecuteWithGlobals(code string, globals map[string]interface{}) (interface{}, error) {
	result, err := e.ExecuteWithOptions(ExecutionOptions{
		Code:    code,
		Globals: globals,
	})
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("execution failed: %w, stderr: %s", result.Error, result.Stderr)
	}

	// 尝试解析结果
	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result.Stdout), &jsonResult); err == nil {
		if res, ok := jsonResult["result"]; ok {
			return res, nil
		}
	}

	return strings.TrimSpace(result.Stdout), nil
}

// ExecuteFile 执行 Python 文件（向后兼容）
func (e *PythonExecutor) ExecuteFile(filepath string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, getPythonPath(), filepath)

	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return stdout.String(), stderr.String(), fmt.Errorf("execution timeout")
	}

	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("execution failed: %w", err)
	}

	return stdout.String(), stderr.String(), nil
}

// ExecuteScript 执行 Python 脚本（从字符串，向后兼容）
func (e *PythonExecutor) ExecuteScript(script string, args ...string) (string, string, error) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "python_script_*.py")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入脚本
	if _, err := tmpFile.WriteString(script); err != nil {
		tmpFile.Close()
		return "", "", fmt.Errorf("failed to write script to temp file: %w", err)
	}
	tmpFile.Close()

	// 执行脚本
	return e.ExecuteFileWithArgs(tmpFile.Name(), args...)
}

// ExecuteFileWithArgs 执行 Python 文件并传递参数
func (e *PythonExecutor) ExecuteFileWithArgs(filepath string, args ...string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	args = append([]string{filepath}, args...)
	cmd := exec.CommandContext(ctx, getPythonPath(), args...)

	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return stdout.String(), stderr.String(), fmt.Errorf("execution timeout")
	}

	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("execution failed: %w", err)
	}

	return stdout.String(), stderr.String(), nil
}

// InstallPackage 安装 Python 包（向后兼容，安装到系统环境）
func (e *PythonExecutor) InstallPackage(packageName string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, getPythonPath(), "-m", "pip", "install", packageName)

	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return stdout.String(), stderr.String(), fmt.Errorf("installation timeout")
	}

	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("installation failed: %w", err)
	}

	return stdout.String(), stderr.String(), nil
}

// Cleanup 清理资源
func (e *PythonExecutor) Cleanup() {
	e.CleanupAllEnvironments()
}

// Cleanup 全局清理函数
func Cleanup() {
	// 全局清理
}
