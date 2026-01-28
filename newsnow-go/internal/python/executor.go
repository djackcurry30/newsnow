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

type PythonExecutor struct {
	mu         sync.Mutex
	timeout    time.Duration
	workingDir string
}

// ExecutionOptions 执行选项
type ExecutionOptions struct {
	Code             string                 `json:"code"`
	Requirements     []string               `json:"requirements,omitempty"`
	RequirementsFile string                 `json:"requirements_file,omitempty"` // requirements.txt 文件路径
	Globals          map[string]interface{} `json:"globals,omitempty"`
	Timeout          time.Duration          `json:"timeout,omitempty"`
}

// NewPythonExecutor 创建一个新的 Python 执行器
func NewPythonExecutor() (*PythonExecutor, error) {
	return &PythonExecutor{
		timeout:    30 * time.Second,
		workingDir: "",
	}, nil
}

// SetTimeout 设置执行超时时间
func (e *PythonExecutor) SetTimeout(timeout time.Duration) {
	e.timeout = timeout
}

// SetWorkingDir 设置工作目录
func (e *PythonExecutor) SetWorkingDir(dir string) {
	e.workingDir = dir
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

// ExecuteWithOptions 使用选项执行 Python 代码（支持虚拟环境和依赖安装）
func (e *PythonExecutor) ExecuteWithOptions(opts ExecutionOptions) (string, string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	timeout := e.timeout
	if opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	// 如果有依赖或 requirements 文件，创建虚拟环境
	needVirtualEnv := len(opts.Requirements) > 0 || opts.RequirementsFile != ""
	var pythonExec string
	if needVirtualEnv {
		// 创建临时目录作为虚拟环境
		venvDir, err := os.MkdirTemp("", "python_venv_*")
		if err != nil {
			return "", "", fmt.Errorf("failed to create temp directory for venv: %w", err)
		}
		defer os.RemoveAll(venvDir)

		venvPath := filepath.Join(venvDir, "venv")

		// 创建虚拟环境
		if err := e.createVirtualEnv(venvPath); err != nil {
			return "", "", err
		}

		// 安装依赖（列表形式）
		if len(opts.Requirements) > 0 {
			if err := e.installRequirements(venvPath, opts.Requirements); err != nil {
				return "", "", err
			}
		}

		// 从 requirements.txt 文件安装依赖
		if opts.RequirementsFile != "" {
			if err := e.installRequirementsFromFile(venvPath, opts.RequirementsFile); err != nil {
				return "", "", err
			}
		}

		pythonExec = e.getVenvPythonPath(venvPath)
	} else {
		pythonExec = getPythonPath()
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "python_script_*.py")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp file: %w", err)
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
		return "", "", fmt.Errorf("failed to write code to temp file: %w", err)
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
		return stdout.String(), stderr.String(), fmt.Errorf("execution timeout")
	}

	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("execution failed: %w", err)
	}

	return stdout.String(), stderr.String(), nil
}

// Execute 执行 Python 代码字符串（向后兼容）
func (e *PythonExecutor) Execute(code string) (string, string, error) {
	return e.ExecuteWithOptions(ExecutionOptions{
		Code: code,
	})
}

// ExecuteWithInput 执行 Python 代码并传递输入数据
func (e *PythonExecutor) ExecuteWithInput(code string, input string) (string, string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

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
	stdout, stderr, err := e.ExecuteWithOptions(ExecutionOptions{
		Code:    code,
		Globals: globals,
	})
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w, stderr: %s", err, stderr)
	}

	// 尝试解析结果
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &result); err == nil {
		if res, ok := result["result"]; ok {
			return res, nil
		}
	}

	return strings.TrimSpace(stdout), nil
}

// ExecuteFile 执行 Python 文件（向后兼容）
func (e *PythonExecutor) ExecuteFile(filepath string) (string, string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

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
	e.mu.Lock()
	defer e.mu.Unlock()

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
	e.mu.Lock()
	defer e.mu.Unlock()

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
	// 清理临时文件等
}

// Cleanup 全局清理函数
func Cleanup() {
	// 全局清理
}
