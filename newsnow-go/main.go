package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"newsnow-go/internal/config"
	"newsnow-go/internal/fetcher/sources"
	"newsnow-go/internal/handler"
	"newsnow-go/internal/middleware"
	"newsnow-go/internal/repository"
	"newsnow-go/internal/service"
	jwtsvc "newsnow-go/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var db *repository.Database
	var cacheRepo *repository.CacheRepository
	var userRepo *repository.UserRepository

	if cfg.EnableCache {
		db, err = repository.NewDatabase(cfg.DatabaseURL)
		if err != nil {
			log.Printf("Warning: Failed to connect to database: %v", err)
		} else {
			_, err = db.Init()
			if err != nil {
				log.Printf("Warning: Failed to init database: %v", err)
			} else {
				cacheRepo = repository.NewCacheRepository(db.GetDB())
				userRepo = repository.NewUserRepository(db.GetDB())

				if cfg.InitTable {
					if err := cacheRepo.Init(); err != nil {
						log.Printf("Warning: Failed to init cache table: %v", err)
					}
					if err := userRepo.Init(); err != nil {
						log.Printf("Warning: Failed to init user table: %v", err)
					}
				}
			}
		}
	}

	cacheService := service.NewCacheService(cacheRepo)
	userService := service.NewUserService(userRepo)
	jwtService := jwtsvc.New(cfg.JWTSecret)
	sourceHolder := sources.NewSourceHolder()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Auth(cfg))

	r.GET("/api/version", handler.NewVersionHandler().Handle)
	r.GET("/api/config", handler.NewConfigHandler(cfg).Handle)
	r.GET("/api/enable-login", handler.NewEnableLoginHandler(cfg).Handle)
	r.GET("/api/user-agent", handler.NewUserAgentHandler().Handle)

	latestHandler := handler.NewLatestHandler(cfg, cacheService, sourceHolder)
	allHandler := handler.NewAllHandler(cfg, cacheService, sourceHolder)
	sHandler := handler.NewSHandler(cfg, cacheService, sourceHolder)
	proxyHandler := handler.NewProxyHandler(cfg, cacheService, sourceHolder)
	loginHandler := handler.NewLoginHandler(cfg, jwtService, userService)
	oauthHandler := handler.NewOAuthHandler(cfg, jwtService, userService)
	meHandler := handler.NewMeHandler(cfg, jwtService, userService)

	r.GET("/api/latest/:id", latestHandler.Handle)
	r.GET("/api/all", allHandler.Handle)
	r.GET("/api/s/:s", sHandler.Handle)
	r.GET("/api/proxy/:s", proxyHandler.Handle)

	authGroup := r.Group("/api")
	authGroup.POST("/login", loginHandler.Handle)
	authGroup.GET("/oauth/:provider", oauthHandler.Handle)
	authGroup.GET("/me", meHandler.Handle)

	r.GET("/api/mcp", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "newsnow",
			"description": "Get the latest news from multiple sources",
			"tools": []gin.H{
				{
					"name":        "latest",
					"description": "Get the latest news from a specific source",
					"parameters": gin.H{
						"type": "object",
						"properties": gin.H{
							"source_id": gin.H{
								"type":        "string",
								"description": "The source ID to fetch news from",
							},
						},
						"required": []string{"source_id"},
					},
				},
			},
		})
	})

	// Python execution API
	pythonHandler, err := handler.NewPythonHandler()
	if err != nil {
		log.Printf("Warning: Failed to initialize Python handler: %v", err)
	} else {
		pythonGroup := r.Group("/api/python")
		{
			pythonGroup.POST("/execute", pythonHandler.Execute)
			pythonGroup.POST("/execute-with-globals", pythonHandler.ExecuteWithGlobals)
			pythonGroup.POST("/execute-file", pythonHandler.ExecuteFile)

			// 环境管理API
			pythonGroup.GET("/environments", pythonHandler.ListEnvironments)
			pythonGroup.GET("/environments/:id", pythonHandler.GetEnvironment)
			pythonGroup.POST("/environments/cleanup", pythonHandler.CleanupEnvironment)
			pythonGroup.POST("/environments/cleanup-all", pythonHandler.CleanupAllEnvironments)
			pythonGroup.POST("/environments/cleanup-expired", pythonHandler.CleanupExpiredEnvironments)
		}
		defer pythonHandler.Cleanup()
	}

	addr := cfg.GetServerAddr()
	log.Printf("Starting server on %s", addr)

	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if cacheRepo != nil {
		cacheRepo.Close()
	}
	if db != nil {
		db.Close()
	}

	log.Println("Server shutdown complete")
}
