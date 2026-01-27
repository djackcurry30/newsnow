# NewsNow Go 后端项目实现方案

## 一、项目概述

本项目旨在使用 Go 语言完全复制 NewsNow 现有 TypeScript 后端的全部功能，将数据库从 SQLite 替换为 PostgreSQL，同时保持原有业务逻辑和数据结构的一致性。NewsNow 是一个功能完善的新闻聚合平台，支持 60 多个不同类型的新闻源，涵盖科技、财经、中国、国际等多个领域，提供实时新闻更新、热搜榜单、用户认证与数据同步等功能。

项目的核心价值在于提供一个高性能、高可用的后端服务，通过 Go 语言的并发优势和优秀的生态库，实现与原 TypeScript 版本完全一致的数据抓取、缓存管理和 API 服务能力。PostgreSQL 数据库的引入将提升数据存储的可靠性、可扩展性和查询性能，同时支持更丰富的数据类型和索引优化。

## 二、技术选型

### 2.1 核心框架

在 Web 框架选择上，我们需要在多个优秀的 Go Web 框架中进行权衡。Gin 是目前最流行的 Go Web 框架，以其高性能、低内存占用和简洁的 API 设计著称，其路由基于 Radix 树实现，支持中间件链式调用，社区活跃度高，文档完善。Echo 是另一个高性能框架，以其极简的设计和优秀的性能表现受到青睐，内置丰富的中间件支持，适合构建 RESTful API。Fiber 基于 Fasthttp 构建，是 Go 语言中最快的 Web 框架之一，API 设计与 Express.js 类似，对前端开发者友好。Hertz 是字节跳动开源的高性能 HTTP 框架，采用自研的网络库，支持 HTTP/1.1 和 HTTP/2 协议，适合大规模微服务场景。

综合考虑性能、生态完善度和学习曲线，本项目推荐使用 Gin 作为核心 Web 框架。Gin 的中间件机制与 Nitro 的 Handler 机制类似，便于移植原有逻辑；同时 GORM 或 pgx 作为数据库 ORM 层，可以很好地映射原有 db0 的数据库操作接口。

### 2.2 数据库驱动

PostgreSQL 数据库连接将使用 pgx 作为底层驱动，它是纯 Go 实现的 PostgreSQL 驱动，支持 PostgreSQL 的全部特性，性能优于传统的 libpq 驱动。在此基础上，我们可以选择使用 GORM 作为 ORM 层，它提供了丰富的模型定义、关联操作和查询构建功能，API 设计直观，支持自动迁移。另一个选择是使用 sqlx 作为更轻量的数据库操作层，保留原生 SQL 的灵活性同时提供结构体映射功能。考虑到项目需要精确复制原有数据库操作逻辑，且新闻源和用户数据的结构相对固定，推荐使用 GORM 作为主要数据库操作方式，必要时结合原生 SQL 实现特定查询。

### 2.3 HTTP 客户端

数据抓取模块需要强大的 HTTP 客户端支持。Go 标准库的 net/http 客户端功能基础，需要自行处理重试、超时和响应解析。Resty 是一个功能丰富的 HTTP 客户端库，支持链式调用、自动重试、JSON 序列化等特性，使用简便。Go-Resty 是 Resty 的改进版本，性能更优。Alpha 则是一个专为爬虫设计的 HTTP 客户端，内置代理支持、Cookie 管理、请求签名等功能，与原项目使用的 ofetch 和 cheerio 组合定位相似。综合考虑，本项目推荐使用 Alpha 进行网页内容抓取，Go-Resty 处理 API 请求，两者配合实现各新闻源的数据获取。

### 2.4 HTML 解析

HTML 内容解析是新闻源抓取的核心环节。goquery 提供了类似 jQuery 的 DOM 操作 API，使用 CSS 选择器即可提取页面元素，是 Go 语言中最流行的 HTML 解析库。Colly 基于回调函数设计，性能优异，适合大规模爬虫场景。Puppy 是另一个轻量级的 HTML 解析库，API 简洁。由于原项目大量使用 cheerio 进行 HTML 解析，goquery 的 jQuery 风格 API 将大大降低迁移成本，推荐作为主要的 HTML 解析工具。

### 2.5 其他依赖

认证方面，将使用 golang-jwt 库处理 JWT 令牌的生成和验证，该库是 Go 语言中最成熟的 JWT 实现，API 与 TypeScript 版本使用的 jose 库类似。MCP 服务器功能将使用官方 Go SDK 实现。日志记录使用 zerolog 或 zap，前者 API 简洁，后者性能更优。配置管理使用 viper，支持环境变量、配置文件等多种配置方式。日期处理使用 dateparse 或 carbon 库，处理各种日期格式的解析。

## 三、目录结构设计

### 3.1 整体目录结构

```
newsnow-go/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
├── internal/
│   ├── config/                  # 配置模块
│   │   ├── config.go            # 配置结构定义
│   │   └── loader.go            # 配置加载逻辑
│   ├── middleware/              # HTTP 中间件
│   │   ├── auth.go              # JWT 认证中间件
│   │   └── cors.go              # 跨域中间件
│   ├── handler/                 # HTTP 处理器
│   │   ├── latest.go            # 最新新闻接口
│   │   ├── login.go             # GitHub OAuth 登录
│   │   ├── oauth.go             # OAuth 回调处理
│   │   ├── me.go                # 用户数据管理
│   │   ├── sync.go              # 数据同步
│   │   ├── search.go            # 新闻搜索
│   │   ├── entire.go            # 批量获取
│   │   ├── enable_login.go      # 登录开关
│   │   ├── mcp.go               # MCP 服务器
│   │   └── version.go           # 版本信息
│   ├── service/                 # 业务逻辑层
│   │   ├── cache.go             # 缓存服务
│   │   ├── user.go              # 用户服务
│   │   ├── source.go            # 新闻源服务
│   │   └── oauth.go             # OAuth 服务
│   ├── repository/              # 数据访问层
│   │   ├── cache_repository.go  # 缓存数据访问
│   │   └── user_repository.go   # 用户数据访问
│   ├── model/                   # 数据模型
│   │   ├── news.go              # 新闻项模型
│   │   ├── source.go            # 新闻源模型
│   │   └── user.go              # 用户模型
│   ├── fetcher/                 # 数据抓取模块
│   │   ├── base.go              # 抓取器基类
│   │   ├── rss.go               # RSS 源抓取
│   │   ├── rsshub.go            # RSSHub 源抓取
│   │   └── sources/             # 各新闻源实现
│   │       ├── v2ex.go
│   │       ├── zhihu.go
│   │       ├── github.go
│   │       └── ...60+ 源文件
│   ├── mcp/                     # MCP 服务器
│   │   ├── server.go            # MCP 服务端
│   │   └── tools.go             # MCP 工具定义
│   ├── utils/                   # 工具函数
│   │   ├── date.go              # 日期处理
│   │   ├── fetch.go             # HTTP 请求
│   │   ├── logger.go            # 日志封装
│   │   ├── crypto.go            # 加密工具
│   │   ├── rss.go               # RSS 解析
│   │   └── constant.go          # 常量定义
│   └── types/                   # 类型定义
│       └── types.go             # 共享类型
├── pkg/
│   ├── jwt/                     # JWT 工具包
│   └── md5/                     # MD5 工具包
├── scripts/                     # 构建脚本
├── config/                      # 配置文件
│   └── config.yaml.example      # 配置示例
├── migrations/                  # 数据库迁移
├── test/                        # 测试文件
├── docs/                        # 文档
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

### 3.2 目录设计说明

cmd 目录存放程序入口文件，遵循 Go 项目最佳实践，将 main.go 放在独立的子目录中可以支持未来添加其他命令（如 CLI 工具）。internal 目录包含项目的内部实现，外部无法导入，保证代码的封装性。其中 config 目录负责配置管理，middleware 目录存放 HTTP 中间件，handler 目录包含所有 API 端点的处理函数，service 目录封装业务逻辑，repository 目录处理数据库操作，model 目录定义数据模型，fetcher 目录是核心的数据抓取模块，mcp 目录实现 MCP 服务器功能，utils 目录提供通用工具函数。

pkg 目录用于存放可复用的公共包，如 JWT 和 MD5 工具，这部分代码未来可以独立提取使用。migrations 目录存放数据库迁移文件，支持版本化管理数据库结构。scripts 目录存放构建和部署脚本，提高开发和部署效率。

## 四、数据库设计

### 4.1 数据库 Schema 设计

PostgreSQL 数据库将创建两个主要表：cache 表用于存储新闻缓存数据，user 表用于存储用户信息。cache 表的设计需要支持高效的键值查询和批量获取操作，user 表则需要支持用户数据的读写和索引查询。

```sql
-- 缓存表：存储新闻源的缓存数据
CREATE TABLE cache (
    id TEXT PRIMARY KEY,
    updated BIGINT NOT NULL,
    data TEXT NOT NULL
);

-- 用户表：存储用户信息和配置数据
CREATE TABLE "user" (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL,
    data TEXT NOT NULL DEFAULT '',
    type TEXT NOT NULL DEFAULT 'github',
    created BIGINT NOT NULL,
    updated BIGINT NOT NULL
);

-- 用户 ID 索引，加速查询
CREATE INDEX idx_user_id ON "user"(id);
```

cache 表使用 TEXT 类型存储 id、data 字段，保持与原 SQLite 实现的数据类型兼容性。updated 字段使用 BIGINT 存储时间戳，与 JavaScript 的 Date.now() 返回的毫秒级时间戳一致。data 字段存储 JSON 格式的新闻数据，使用 TEXT 类型可以容纳任意长度的 JSON 字符串。

user 表的 id 字段存储用户唯一标识，在 GitHub OAuth 认证场景下为 GitHub 用户 ID。type 字段标识用户认证类型，目前支持 github，未来可以扩展其他 OAuth 提供商。data 字段存储用户的配置数据，包括订阅的新闻源、自定义栏目布局等，以 JSON 格式存储。

### 4.2 GORM 模型定义

```go
package model

import (
    "time"
)

type Cache struct {
    ID      string    `gorm:"primaryKey;column:id"`
    Updated int64     `gorm:"column:updated;not null"`
    Data    string    `gorm:"column:data;not null"`
}

func (Cache) TableName() string {
    return "cache"
}

type User struct {
    ID      string    `gorm:"primaryKey;column:id"`
    Email   string    `gorm:"column:email;not null"`
    Data    string    `gorm:"column:data;not null;default:''"`
    Type    string    `gorm:"column:type;not null;default:'github'"`
    Created int64     `gorm:"column:created;not null"`
    Updated int64     `gorm:"column:updated;not null"`
}

func (User) TableName() string {
    return "user"
}
```

模型定义使用 GORM 的标签语法指定数据库列名、主键和约束条件。TableName 方法用于自定义表名，user 表需要使用双引号包裹，因为 user 是 PostgreSQL 的保留关键字。

### 4.3 数据库初始化

数据库初始化需要检查表是否存在，如果不存在则创建。GORM 的 AutoMigrate 功能可以自动创建表结构，但在生产环境中通常建议使用显式的迁移脚本以确保数据安全。

```go
package repository

import (
    "newsnow-go/internal/model"
    "newsnow-go/pkg/logger"

    "gorm.io/gorm"
)

type Database struct {
    db *gorm.DB
}

func NewDatabase(dsn string) (*Database, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return &Database{db: db}, nil
}

func (d *Database) Init() error {
    // 创建缓存表
    if err := d.db.AutoMigrate(&model.Cache{}); err != nil {
        return err
    }
    // 创建用户表
    if err := d.db.AutoMigrate(&model.User{}); err != nil {
        return err
    }
    logger.Success("Database tables initialized")
    return nil
}

func (d *Database) GetDB() *gorm.DB {
    return d.db
}
```

## 五、核心模块实现

### 5.1 配置模块

配置模块负责加载和管理应用配置，支持环境变量和配置文件两种方式。配置结构需要包含数据库连接信息、GitHub OAuth 凭证、JWT 密钥等关键参数。

```go
package config

import (
    "fmt"
    "os"
    "strconv"

    "github.com/spf13/viper"
)

type Config struct {
    // Server
    Host string `mapstructure:"HOST"`
    Port int    `mapstructure:"PORT"`

    // Database
    DatabaseURL string `mapstructure:"DATABASE_URL"`

    // GitHub OAuth
    GClientID     string `mapstructure:"G_CLIENT_ID"`
    GClientSecret string `mapstructure:"G_CLIENT_SECRET"`

    // JWT
    JWTSecret string `mapstructure:"JWT_SECRET"`

    // Feature flags
    InitTable   bool `mapstructure:"INIT_TABLE"`
    EnableCache bool `mapstructure:"ENABLE_CACHE"`

    // MCP
    BaseURL string `mapstructure:"BASE_URL"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("config/")
    viper.AddConfigPath(".")

    // 加载配置文件
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }

    // 环境变量覆盖
    viper.AutomaticEnv()

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    // 环境变量优先级更高
    if host := os.Getenv("HOST"); host != "" {
        cfg.Host = host
    }
    if port := os.Getenv("PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            cfg.Port = p
        }
    }

    return &cfg, nil
}

func (c *Config) GetDSN() string {
    return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
        c.DatabaseURL /* 解析 host、dbname 等 */)
}
```

配置模块使用 viper 库实现配置管理，支持 YAML 配置文件和环境变量两种配置方式。环境变量的优先级高于配置文件，这样可以方便地在不同环境（开发、测试、生产）中使用不同的配置。GetDSN 方法用于构建 PostgreSQL 连接字符串。

### 5.2 日志模块

日志模块封装日志输出，提供不同级别的日志记录功能。参照原项目使用 consola 的设计，Go 版本使用 zerolog 或 slog 实现类似的日志输出格式。

```go
package logger

import (
    "os"
    "time"

    "github.com/rs/zerolog"
)

var log zerolog.Logger

func Init() {
    zerolog.TimeFieldFormat = time.RFC3339
    zerolog.SetGlobalLevel(zerolog.DebugLevel)

    log = zerolog.New(os.Stdout).
        With().
        Timestamp().
        Caller().
        Logger()
}

func Success(v interface{}) {
    log.Info().Msgf("%v", v)
}

func Info(v interface{}) {
    log.Info().Msgf("%v", v)
}

func Warn(v interface{}) {
    log.Warn().Msgf("%v", v)
}

func Error(v interface{}) {
    log.Error().Msgf("%v", v)
}

func Errorf(format string, v ...interface{}) {
    log.Error().Msgf(format, v...)
}
```

### 5.3 认证中间件

认证中间件负责验证 JWT Token 并将用户信息注入请求上下文。该中间件需要支持可选认证模式，即某些 API 端点可以在未登录状态下返回缓存数据。

```go
package middleware

import (
    "errors"
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

type UserContext struct {
    ID   string `json:"id"`
    Type string `json:"type"`
}

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查环境变量配置
        if os.Getenv("JWT_SECRET") == "" ||
           os.Getenv("G_CLIENT_ID") == "" ||
           os.Getenv("G_CLIENT_SECRET") == "" {
            c.Set("disabledLogin", true)
            // 部分接口仍需登录
            if strings.HasPrefix(c.Request.URL.Path, "/api/s") ||
               strings.HasPrefix(c.Request.URL.Path, "/api/me") {
                c.JSON(http.StatusUpgradeRequired, gin.H{
                    "message": "Server not configured, disable login",
                })
                c.Abort()
                return
            }
            c.Next()
            return
        }

        // 检查是否需要认证的路径
        if strings.HasPrefix(c.Request.URL.Path, "/api/s") ||
           strings.HasPrefix(c.Request.URL.Path, "/api/me") {
            authHeader := c.GetHeader("Authorization")
            if authHeader != "" {
                // 移除 "Bearer " 前缀
                tokenString := strings.TrimPrefix(authHeader, "Bearer ")
                tokenString = strings.TrimSpace(tokenString)

                token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, errors.New("unexpected signing method")
                    }
                    return []byte(os.Getenv("JWT_SECRET")), nil
                })

                if err == nil && token.Valid {
                    if claims, ok := token.Claims.(jwt.MapClaims); ok {
                        if id, ok := claims["id"].(string); ok {
                            c.Set("user", UserContext{
                                ID:   id,
                                Type: claims["type"].(string),
                            })
                        }
                    }
                } else {
                    if strings.HasPrefix(c.Request.URL.Path, "/api/me") {
                        c.JSON(http.StatusUnauthorized, gin.H{
                            "message": "JWT verification failed",
                        })
                        c.Abort()
                        return
                    }
                    logger.Warn("JWT verification failed")
                }
            } else if strings.HasPrefix(c.Request.URL.Path, "/api/me") {
                c.JSON(http.StatusUnauthorized, gin.H{
                    "message": "JWT verification failed",
                })
                c.Abort()
                return
            }
        }

        c.Next()
    }
}
```

### 5.4 缓存服务

缓存服务封装缓存数据的读写操作，支持单条获取、批量获取和写入功能。缓存策略与原项目保持一致，使用 TTL 和刷新间隔两层控制。

```go
package service

import (
    "encoding/json"
    "newsnow-go/internal/model"
    "newsnow-go/internal/repository"
    "newsnow-go/shared"
    "time"
)

type CacheService struct {
    repo *repository.CacheRepository
    TTL  time.Duration
}

func NewCacheService(repo *repository.CacheRepository) *CacheService {
    return &CacheService{
        repo: repo,
        TTL:  30 * time.Minute, // 默认 30 分钟 TTL
    }
}

func (s *CacheService) Get(sourceID string) (*CacheInfo, error) {
    cache, err := s.repo.Get(sourceID)
    if err != nil {
        return nil, err
    }
    if cache == nil {
        return nil, nil
    }

    var items []shared.NewsItem
    if err := json.Unmarshal([]byte(cache.Data), &items); err != nil {
        return nil, err
    }

    return &CacheInfo{
        ID:      cache.ID,
        Updated: cache.Updated,
        Items:   items,
    }, nil
}

func (s *CacheService) Set(sourceID string, items []shared.NewsItem) error {
    data, err := json.Marshal(items)
    if err != nil {
        return err
    }
    return s.repo.Set(sourceID, string(data), time.Now().UnixMilli())
}

type CacheInfo struct {
    ID      string
    Updated int64
    Items   []shared.NewsItem
}
```

### 5.5 新闻源抓取器基类

新闻源抓取器基类定义统一的接口规范，所有具体新闻源的抓取器都需要实现该接口。

```go
package fetcher

import (
    "newsnow-go/shared"
)

type SourceFetcher interface {
    Fetch() ([]shared.NewsItem, error)
    GetID() string
}

type BaseFetcher struct {
    ID string
}

func (b *BaseFetcher) GetID() string {
    return b.ID
}

// defineSource 用于定义一个或多个新闻源
func defineSource(fetchers interface{}) interface{} {
    return fetchers
}

// defineRSSSource 从 RSS 源获取数据
func defineRSSSource(url string, hiddenDate bool) SourceFetcher {
    return &RSSFetcher{
        URL:        url,
        HiddenDate: hiddenDate,
    }
}

// defineRSSHubSource 从 RSSHub 获取数据
func defineRSSHubSource(route string, options RSSHubOptions, hiddenDate bool) SourceFetcher {
    return &RSSHubFetcher{
        Route:      route,
        Options:    options,
        HiddenDate: hiddenDate,
    }
}

type RSSHubOptions struct {
    Sorted *bool  // 是否排序，默认 true
    Limit  *int   // 限制数量，默认 20
}

type RSSFetcher struct {
    BaseFetcher
    URL        string
    HiddenDate bool
}

func (f *RSSFetcher) Fetch() ([]shared.NewsItem, error) {
    // 实现 RSS 解析逻辑
}

type RSSHubFetcher struct {
    BaseFetcher
    Route      string
    Options    RSSHubOptions
    HiddenDate bool
}

func (f *RSSHubFetcher) Fetch() ([]shared.NewsItem, error) {
    // 实现 RSSHub 解析逻辑
}
```

### 5.6 API 处理器

API 处理器实现各个 API 端点的业务逻辑。以下是 latest 接口的实现示例。

```go
package handler

type LatestHandler struct {
    cacheService *service.CacheService
    fetchers     map[string]fetcher.SourceFetcher
}

func NewLatestHandler(cacheService *service.CacheService, fetchers map[string]fetcher.SourceFetcher) *LatestHandler {
    return &LatestHandler{
        cacheService: cacheService,
        fetchers:     fetchers,
    }
}

func (h *LatestHandler) Handle(c *gin.Context) {
    query := c.Query("id")
    latestParam := c.Query("latest")
    latest := latestParam != "" && latestParam != "false"

    // 验证 source ID
    if query == "" || !shared.Sources[query] {
        // 检查是否有重定向
        if redirectID := shared.Sources[query].Redirect; redirectID != "" {
            query = redirectID
        } else {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "Invalid source id",
            })
            return
        }
    }

    // 尝试从缓存获取
    if h.cacheService != nil {
        cache, err := h.cacheService.Get(query)
        if err == nil && cache != nil {
            interval := shared.Sources[query].Interval
            now := time.Now().UnixMilli()

            // 检查是否在刷新间隔内
            if now-cache.Updated < interval {
                c.JSON(http.StatusOK, shared.SourceResponse{
                    Status:      "success",
                    ID:          query,
                    UpdatedTime: now,
                    Items:       cache.Items,
                })
                return
            }

            // 检查 TTL 缓存
            if now-cache.Updated < int64(h.cacheService.TTL.Milliseconds()) {
                if !latest || !h.hasUser(c) {
                    c.JSON(http.StatusOK, shared.SourceResponse{
                        Status:      "cache",
                        ID:          query,
                        UpdatedTime: cache.Updated,
                        Items:       cache.Items,
                    })
                    return
                }
            }
        }
    }

    // 获取新数据
    fetcher := h.fetchers[query]
    if fetcher == nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "Fetcher not found",
        })
        return
    }

    items, err := fetcher.Fetch()
    if err != nil {
        logger.Error(err)
        // 返回缓存数据（如果有）
        if cache, _ := h.cacheService.Get(query); cache != nil {
            c.JSON(http.StatusOK, shared.SourceResponse{
                Status:      "cache",
                ID:          query,
                UpdatedTime: cache.Updated,
                Items:       cache.Items,
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": err.Error(),
        })
        return
    }

    // 限制最多 30 条
    if len(items) > 30 {
        items = items[:30]
    }

    // 更新缓存
    if h.cacheService != nil {
        go h.cacheService.Set(query, items)
    }

    logger.Success("fetch " + query + " latest")
    c.JSON(http.StatusOK, shared.SourceResponse{
        Status:      "success",
        ID:          query,
        UpdatedTime: time.Now().UnixMilli(),
        Items:       items,
    })
}

func (h *LatestHandler) hasUser(c *gin.Context) bool {
    _, exists := c.Get("user")
    disabledLogin, _ := c.Get("disabledLogin")
    return exists && !disabledLogin.(bool)
}
```

### 5.7 OAuth 处理

GitHub OAuth 处理的完整流程包括接收回调、交换访问令牌、获取用户信息、创建或更新用户、生成 JWT 令牌并重定向前端。

```go
package handler

func (h *OAuthHandler) Callback(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Missing code parameter",
        })
        return
    }

    // 1. 交换访问令牌
    tokenResp, err := h.oauthService.ExchangeToken(code)
    if err != nil {
        logger.Error(err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to exchange token",
        })
        return
    }

    // 2. 获取用户信息
    userInfo, err := h.oauthService.GetUserInfo(tokenResp.AccessToken)
    if err != nil {
        logger.Error(err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to get user info",
        })
        return
    }

    // 3. 创建或更新用户
    userID := strconv.FormatInt(userInfo.ID, 10)
    email := userInfo.NotificationEmail
    if email == "" {
        email = userInfo.Email
    }
    h.userService.AddUser(userID, email, "github")

    // 4. 生成 JWT
    jwtToken, err := h.jwtService.GenerateToken(userID, "github")
    if err != nil {
        logger.Error(err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to generate token",
        })
        return
    }

    // 5. 重定向回前端
    params := url.Values{
        "login": []string{"github"},
        "jwt":   []string{jwtToken},
        "user":  []string{fmt.Sprintf(`{"avatar":"%s","name":"%s"}`, userInfo.AvatarURL, userInfo.Name)},
    }
    c.Redirect(http.StatusFound, "/?"+params.Encode())
}
```

### 5.8 MCP 服务器

MCP 服务器提供 get_hotest_latest_news 工具，允许外部系统通过 MCP 协议获取新闻数据。

```go
package mcp

import (
    "context"
    "encoding/json"
    "newsnow-go/internal/service"
    "newsnow-go/shared"

    "github.com/metoro-io/mcp-golang"
    "github.com/metoro-io/mcp-golang/transport/streamablehttp"
)

type MCPServer struct {
    server     *mcp.Server
    cacheService *service.CacheService
    fetchers   map[string]fetcher.SourceFetcher
}

func NewMCPServer(cacheService *service.CacheService, fetchers map[string]fetcher.SourceFetcher) *MCPServer {
    server := mcp.NewServer("NewsNow", "0.0.1", mcp.WithCapabilities(&mcp.ServerCapabilities{
        Logging: &mcp.LoggingCapabilities{},
    }))

    mcpserver := &MCPServer{
        server:       server,
        cacheService: cacheService,
        fetchers:     fetchers,
    }

    mcpserver.registerTools()

    return mcpserver
}

func (s *MCPServer) registerTools() {
    s.server.AddTool(mcp.NewTool("get_hotest_latest_news",
        mcp.WithDescription("get hotest or latest news from source by id, return count news"),
        mcp.WithParameter("id", mcp.String, mcp.Description("source id, e.g. v2ex, zhihu, github")),
        mcp.WithParameter("count", mcp.Int, mcp.Default(10), mcp.Description("count of news to return")),
        func(ctx context.Context, request mcp.Request) (mcp.Result, error) {
            params := request.Params.Arguments
            id := params["id"].(string)
            count := int(params["count"].(float64))
            if count < 1 {
                count = 10
            }

            fetcher := s.fetchers[id]
            if fetcher == nil {
                return mcp.NewErrorResult(nil), nil
            }

            items, err := fetcher.Fetch()
            if err != nil {
                return mcp.NewErrorResult(nil), err
            }

            if len(items) > count {
                items = items[:count]
            }

            content := make([]mcp.Content, len(items))
            for i, item := range items {
                content[i] = mcp.TextContent{
                    Type: "text",
                    Text: "[" + item.Title + "](" + item.URL + ")",
                }
            }

            return mcp.NewContentResult(content...), nil
        },
    ))
}

func (s *MCPServer) Run(addr string) error {
    transport := streamablehttp.NewStreamableHTTPServerTransport(addr)
    transport.SetMCPServer(s.server)
    return s.server.Serve(transport)
}
```

## 六、新闻源实现

### 6.1 新闻源列表

原项目支持 60 多个新闻源，需要逐一实现 Go 版本。以下是部分主要新闻源的列表和分类。

科技类新闻源包括 V2EX（v2ex、v2ex-share）、GitHub Trending（github、github-trending-today）、36氪（36kr、36kr-quick、36kr-renqi）、IT之家（ithome）、少数派（sspai）、稀土掘金（juejin）、Solidot（solidot）、Hacker News（hackernews）、Product Hunt（producthunt）、远景论坛（pcbeta、pcbeta-windows11）、Freebuf（freebuf）和酷安（coolapk）。

财经类新闻源包括 金十数据（jin10）、华尔街见闻（wallstreetcn、wallstreetcn-quick、wallstreetcn-news、wallstreetcn-hot）、财联社（cls、cls-telegraph、cls-depth、cls-hot）、雪球（xueqiu、xueqiu-hotstock）、格隆汇（gelonghui）、法布财经（fastbull、fastbull-express、fastbull-news）、MKTNews（mktnews、mktnews-flash）。

中国区新闻源包括 知乎（zhihu）、微博（weibo）、抖音（douyin）、虎扑（hupu）、百度贴吧（tieba）、今日头条（toutiao）、哔哩哔哩（bilibili、bilibili-hot-search、bilibili-hot-video、bilibili-ranking）、快手（kuaishou）、腾讯新闻（tencent、tencent-hot）、爱奇艺（iqiyi、iqiyi-hot-ranklist）、腾讯视频（qqvideo、qqvideo-tv-hotsearch）、豆瓣（douban）、百度热搜（baidu）、牛客（nowcoder）、凤凰网（ifeng）、虫部落（chongbuluo、chongbuluo-latest、chongbuluo-hot）。

国际区新闻源包括 联合早报（zaobao）、卫星通讯社（sputniknewscn）、参考消息（cankaoxiaoxi）、靠谱新闻（kaopu）、Steam（steam）。

### 6.2 新闻源实现示例

以下是 V2EX 新闻源的 Go 实现示例，展示了如何将 TypeScript 逻辑转换为 Go。

```go
package sources

import (
    "encoding/json"
    "newsnow-go/internal/fetcher"
    "newsnow-go/internal/utils"
    "newsnow-go/shared"
)

type V2EXRes struct {
    Version      string `json:"version"`
    Title        string `json:"title"`
    Description  string `json:"description"`
    HomePageURL  string `json:"home_page_url"`
    FeedURL      string `json:"feed_url"`
    Icon         string `json:"icon"`
    Favicon      string `json:"favicon"`
    Items        []V2EXItem `json:"items"`
}

type V2EXItem struct {
    URL           string `json:"url"`
    DateModified  string `json:"date_modified,omitempty"`
    DatePublished string `json:"date_published"`
    ContentHTML   string `json:"content_html"`
    Title         string `json:"title"`
    ID            string `json:"id"`
}

type V2EXFetcher struct {
    fetcher.BaseFetcher
}

func NewV2EXFetcher() *V2EXFetcher {
    return &V2EXFetcher{
        BaseFetcher: fetcher.BaseFetcher{ID: "v2ex"},
    }
}

func (f *V2EXFetcher) Fetch() ([]shared.NewsItem, error) {
    categories := []string{"create", "ideas", "programmer", "share"}
    results := make([]shared.NewsItem, 0)

    for _, cat := range categories {
        url := "https://www.v2ex.com/feed/" + cat + ".json"
        var res V2EXRes
        if err := utils.FetchJSON(url, &res); err != nil {
            return nil, err
        }
        for _, item := range res.Items {
            date := item.DateModified
            if date == "" {
                date = item.DatePublished
            }
            results = append(results, shared.NewsItem{
                ID:   item.ID,
                Title: item.Title,
                URL:   item.URL,
                PubDate: date,
                Extra: &shared.NewsExtra{
                    Date: date,
                },
            })
        }
    }

    // 按日期排序
    sortByDate(results)

    return results, nil
}

func sortByDate(items []shared.NewsItem) {
    // 实现日期排序逻辑
}
```

### 6.3 通用 RSS 源实现

对于使用 RSS 订阅的网站，可以使用通用的 RSS 抓取器实现。

```go
package fetcher

import (
    "encoding/xml"
    "newsnow-go/internal/utils"
    "newsnow-go/shared"
    "time"
)

type RSSFetcher struct {
    BaseFetcher
    URL        string
    HiddenDate bool
}

type RSSChannel struct {
    Title       string    `xml:"title"`
    Description string    `xml:"description"`
    Link        string    `xml:"link"`
    Image       RSSImage  `xml:"image"`
    Items       []RSSItem `xml:"item"`
}

type RSSImage struct {
    URL string `xml:"url"`
}

type RSSItem struct {
    Title       string `xml:"title"`
    Description string `xml:"description"`
    Link        string `xml:"link"`
    PubDate     string `xml:"pubDate"`
    GUID        string `xml:"guid"`
}

func (f *RSSFetcher) Fetch() ([]shared.NewsItem, error) {
    var channel RSSChannel
    if err := utils.FetchXML(f.URL, &channel); err != nil {
        return nil, err
    }

    items := make([]shared.NewsItem, 0, len(channel.Items))
    for _, item := range channel.Items {
        pubDate := ""
        if !f.HiddenDate && item.PubDate != "" {
            if t, err := time.Parse(time.RFC1123, item.PubDate); err == nil {
                pubDate = t.Format(time.RFC3339)
            }
        }
        items = append(items, shared.NewsItem{
            ID:      item.GUID,
            Title:   item.Title,
            URL:     item.Link,
            PubDate: pubDate,
        })
    }

    return items, nil
}
```

## 七、构建与部署

### 7.1 Makefile

```makefile
.PHONY: all build run dev test lint clean docker

all: build

build:
    go build -o bin/newsnow ./cmd/server

run: build
    ./bin/newsnow

dev:
    air

test:
    go test ./...

lint:
    golangci-lint run

clean:
    rm -rf bin/ coverage.out

docker:
    docker build -t newsnow-go:latest .

docker-run:
    docker run -p 8080:8080 -e DATABASE_URL="host=db user=postgres password=secret dbname=newsnow sslmode=disable" newsnow-go:latest
```

### 7.2 Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o newsnow ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/newsnow .
COPY config ./config

ENV HOST=0.0.0.0 PORT=8080
EXPOSE 8080

CMD ["./newsnow"]
```

### 7.3 docker-compose.yml

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - '8080:8080'
    environment:
      - HOST=0.0.0.0
      - PORT=8080
      - DATABASE_URL=postgresql://postgres:secret@db:5432/newsnow?sslmode=disable
      - G_CLIENT_ID=${G_CLIENT_ID}
      - G_CLIENT_SECRET=${G_CLIENT_SECRET}
      - JWT_SECRET=${JWT_SECRET}
      - INIT_TABLE=true
      - ENABLE_CACHE=true
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=secret
      - POSTGRES_DB=newsnow
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_data:
```

## 八、开发进度规划

### 第一阶段：基础设施（1-2 天）

本阶段完成项目骨架搭建和核心基础设施配置。主要任务包括初始化 Go Module、创建目录结构、实现配置加载模块、实现日志模块、实现数据库连接和 GORM 模型定义。本阶段的交付物是一个可以正常启动和关闭的基础项目框架。

### 第二阶段：核心服务（2-3 天）

本阶段实现系统的核心服务和业务逻辑层。主要任务包括实现缓存服务（CacheService）、实现用户服务（UserService）、实现认证中间件（JWT 验证）、实现 OAuth 服务（GitHub 登录流程）。本阶段完成后，应用具备用户认证和数据存储的基本能力。

### 第三阶段：数据抓取（3-5 天）

本阶段实现新闻源数据抓取模块，这是工作量最大的阶段。主要任务包括实现基础抓取器框架、实现 RSS 和 RSSHub 抓取器、实现各分类新闻源（共 60+ 个源）。考虑到新闻源数量较多，可以按照分类并行开发，预计需要 3-5 天完成全部抓取器的实现。

### 第四阶段：API 层（1-2 天）

本阶段实现所有 API 端点。主要任务包括实现 latest 接口、实现 search 接口、实现 entire 批量获取接口、实现用户数据同步接口、实现登录相关接口、实现 MCP 服务器接口。API 层相对简单，预计 1-2 天可以完成。

### 第五阶段：测试与优化（1-2 天）

本阶段进行全面的测试和性能优化。主要任务包括编写单元测试、集成测试、压力测试、性能调优、Bug 修复。测试覆盖率应达到 80% 以上，确保系统的稳定性和可靠性。

### 第六阶段：文档与部署（1 天）

本阶段完成项目文档编写和部署配置。主要任务包括完善 README 文档、编写 API 文档、配置 Docker 部署、编写 CI/CD 配置文件。本阶段完成后，项目具备生产环境部署能力。

## 九、总结

本实现方案详细规划了使用 Go 语言复制 NewsNow 后端功能的完整路径。通过合理的技术选型、清晰的目录结构设计、完善的数据库 schema 定义，以及模块化的实现思路，可以确保项目的高质量交付。PostgreSQL 数据库的引入将提升系统的数据存储能力和可靠性，而 Go 语言的高性能特性将为系统带来更好的响应速度和并发处理能力。整个项目预计需要 9-15 个工作日完成开发，经过测试和优化后可以投入生产使用。
