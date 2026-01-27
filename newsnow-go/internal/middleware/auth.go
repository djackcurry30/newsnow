package middleware

import (
	"net/http"
	"strings"

	"newsnow-go/internal/config"
	"newsnow-go/pkg/jwt"
	"newsnow-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserContext struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func Auth(cfg *config.Config) gin.HandlerFunc {
	jwtService := jwt.New(cfg.JWTSecret)
	
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		
		if !strings.HasPrefix(path, "/api") {
			c.Next()
			return
		}

		disabledLogin := !cfg.IsLoginEnabled()
		c.Set("disabledLogin", disabledLogin)

		if disabledLogin {
			if !isPublicAPI(path) {
				c.JSON(http.StatusUpgradeRequired, gin.H{
					"message": "Server not configured, disable login",
				})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		if requiresAuth(path) {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				tokenString := strings.TrimPrefix(authHeader, "Bearer")
				tokenString = strings.TrimSpace(tokenString)
				
				claims, err := jwtService.VerifyToken(tokenString)
				if err == nil && claims != nil {
					c.Set("user", UserContext{
						ID:   claims.ID,
						Type: claims.Type,
					})
				} else {
					if strings.HasPrefix(path, "/api/me") {
						c.JSON(http.StatusUnauthorized, gin.H{
							"message": "JWT verification failed",
						})
						c.Abort()
						return
					}
					utils.Warn("JWT verification failed")
				}
			} else if strings.HasPrefix(path, "/api/me") {
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

func isPublicAPI(path string) bool {
	publicAPIs := []string{
		"/api/s",
		"/api/proxy",
		"/api/latest",
		"/api/mcp",
		"/api/enable-login",
		"/api/version",
	}
	for _, api := range publicAPIs {
		if strings.HasPrefix(path, api) {
			return true
		}
	}
	return false
}

func requiresAuth(path string) bool {
	protectedAPIs := []string{
		"/api/s",
		"/api/me",
	}
	for _, api := range protectedAPIs {
		if strings.HasPrefix(path, api) {
			return true
		}
	}
	return false
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}
