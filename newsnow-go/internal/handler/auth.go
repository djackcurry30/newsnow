package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"newsnow-go/internal/config"
	jwtsvc "newsnow-go/pkg/jwt"
	"newsnow-go/pkg/md5"
	"newsnow-go/internal/service"
	"newsnow-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserContext struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type LoginHandler struct {
	cfg        *config.Config
	jwt        *jwtsvc.JWT
	userService *service.UserService
}

func NewLoginHandler(cfg *config.Config, jwtService *jwtsvc.JWT, userService *service.UserService) *LoginHandler {
	return &LoginHandler{
		cfg:        cfg,
		jwt:        jwtService,
		userService: userService,
	}
}

func (h *LoginHandler) Handle(c *gin.Context) {
	disabledLogin := !h.cfg.IsLoginEnabled()
	c.Set("disabledLogin", disabledLogin)
	
	if disabledLogin {
		c.JSON(http.StatusUpgradeRequired, gin.H{
			"message": "Server not configured, disable login",
		})
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing email or password"})
		return
	}

	allowedEmails := h.cfg.GetAllowedEmails()
	if len(allowedEmails) > 0 {
		found := false
		for _, email := range allowedEmails {
			if strings.EqualFold(email, req.Email) {
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Email not allowed"})
			return
		}
	}

	passwordHash := md5.Sum(req.Password)
	secret := h.cfg.JWTSecret + string(passwordHash)
	
	now := time.Now().UnixMilli()
	secretHash := md5.SumBase64(secret)
	userID := md5.Sum(req.Email + secretHash)
	
	userType := "password"
	if err := h.userService.AddUser(userID, req.Email, userType); err != nil {
		utils.Error("Failed to add user: " + err.Error())
	}

	token, err := h.jwt.GenerateToken(userID, userType, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"user":     userID,
		"type":     userType,
		"email":    req.Email,
		"create":   true,
		"disabled": disabledLogin,
		"now":      now,
	})
}

type OAuthHandler struct {
	cfg        *config.Config
	jwt        *jwtsvc.JWT
	userService *service.UserService
}

func NewOAuthHandler(cfg *config.Config, jwtService *jwtsvc.JWT, userService *service.UserService) *OAuthHandler {
	return &OAuthHandler{
		cfg:        cfg,
		jwt:        jwtService,
		userService: userService,
	}
}

func (h *OAuthHandler) Handle(c *gin.Context) {
	provider := c.Param("provider")

	if provider != "github" && provider != "google" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid provider"})
		return
	}

	// Check if provider is configured
	if provider == "github" && !h.cfg.IsGitHubOAuthEnabled() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "GitHub OAuth not configured"})
		return
	}
	if provider == "google" && !h.cfg.IsGoogleOAuthEnabled() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Google OAuth not configured"})
		return
	}

	code := c.Query("code")
	if code == "" {
		state := c.Query("state")
		var authURL string

		if provider == "github" {
			redirectURI := h.cfg.BaseURL + "/api/oauth/github"
			authURL = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user",
				h.cfg.GClientID, url.QueryEscape(redirectURI))
			if state != "" {
				authURL += "&state=" + url.QueryEscape(state)
			}
		} else {
			redirectURI := h.cfg.BaseURL + "/api/oauth/google"
			authURL = fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid profile email",
				h.cfg.GoogleClientID, url.QueryEscape(redirectURI))
			if state != "" {
				authURL += "&state=" + url.QueryEscape(state)
			}
		}

		c.Redirect(http.StatusFound, authURL)
		return
	}

	tokenResponse, err := h.exchangeCode(provider, code)
	if err != nil {
		utils.Error("OAuth token exchange failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to exchange code"})
		return
	}

	userInfo, err := h.getUserInfo(provider, tokenResponse)
	if err != nil {
		utils.Error("Failed to get user info: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user info"})
		return
	}

	userID := userInfo["id"].(string)
	email := userInfo["email"].(string)
	
	if err := h.userService.AddUser(userID, email, provider); err != nil {
		utils.Error("Failed to add user: " + err.Error())
	}

	token, err := h.jwt.GenerateToken(userID, provider, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	redirectURL := h.cfg.BaseURL + "#access_token=" + token + "&provider=" + provider + "&create=true&disabled=" + strconv.FormatBool(!h.cfg.IsLoginEnabled())
	c.Redirect(http.StatusFound, redirectURL)
}

func (h *OAuthHandler) exchangeCode(provider, code string) (map[string]string, error) {
	var tokenURL, clientID, clientSecret string

	if provider == "github" {
		tokenURL = "https://github.com/login/oauth/access_token"
		clientID = h.cfg.GClientID
		clientSecret = h.cfg.GClientSecret
	} else {
		tokenURL = "https://oauth2.googleapis.com/token"
		clientID = h.cfg.GoogleClientID
		clientSecret = h.cfg.GoogleClientSecret
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", h.cfg.BaseURL+"/api/oauth/"+provider)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if _, ok := result["error"]; ok {
		return nil, nil
	}

	token := result["access_token"].(string)
	return map[string]string{"access_token": token}, nil
}

func (h *OAuthHandler) getUserInfo(provider string, tokenResponse map[string]string) (map[string]interface{}, error) {
	token := tokenResponse["access_token"]
	
	var userURL string
	if provider == "github" {
		userURL = "https://api.github.com/user"
	} else {
		userURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	}

	req, _ := http.NewRequest("GET", userURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "NewsNow")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

type MeHandler struct {
	cfg        *config.Config
	jwt        *jwtsvc.JWT
	userService *service.UserService
}

func NewMeHandler(cfg *config.Config, jwtService *jwtsvc.JWT, userService *service.UserService) *MeHandler {
	return &MeHandler{
		cfg:        cfg,
		jwt:        jwtService,
		userService: userService,
	}
}

func (h *MeHandler) Handle(c *gin.Context) {
	disabledLogin := !h.cfg.IsLoginEnabled()
	c.Set("disabledLogin", disabledLogin)
	
	if disabledLogin {
		c.JSON(http.StatusUpgradeRequired, gin.H{
			"message": "Server not configured, disable login",
		})
		return
	}

	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "JWT verification failed"})
		return
	}

	userContext := userValue.(UserContext)
	
	user, err := h.userService.GetUser(userContext.ID)
	if err != nil {
		utils.Error("Failed to get user: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      user.ID,
		"email":   user.Email,
		"type":    user.Type,
		"data":    user.Data,
		"created": user.Created,
		"updated": user.Updated,
	})
}
