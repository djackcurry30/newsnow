package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Host               string   `mapstructure:"HOST"`
	Port               int      `mapstructure:"PORT"`
	DatabaseURL        string   `mapstructure:"DATABASE_URL"`
	GClientID          string   `mapstructure:"G_CLIENT_ID"`
	GClientSecret      string   `mapstructure:"G_CLIENT_SECRET"`
	GoogleClientID     string   `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string   `mapstructure:"GOOGLE_CLIENT_SECRET"`
	JWTSecret          string   `mapstructure:"JWT_SECRET"`
	InitTable          bool     `mapstructure:"INIT_TABLE"`
	EnableCache        bool     `mapstructure:"ENABLE_CACHE"`
	BaseURL            string   `mapstructure:"BASE_URL"`
	ProductHuntToken   string   `mapstructure:"PRODUCTHUNT_API_TOKEN"`
	AllowedEmails      []string `mapstructure:"ALLOWED_EMAILS"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if host := os.Getenv("HOST"); host != "" {
		cfg.Host = host
	}
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.DatabaseURL = dbURL
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.JWTSecret = jwtSecret
	}
	if clientID := os.Getenv("G_CLIENT_ID"); clientID != "" {
		cfg.GClientID = clientID
	}
	if clientSecret := os.Getenv("G_CLIENT_SECRET"); clientSecret != "" {
		cfg.GClientSecret = clientSecret
	}
	if googleClientID := os.Getenv("GOOGLE_CLIENT_ID"); googleClientID != "" {
		cfg.GoogleClientID = googleClientID
	}
	if googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET"); googleClientSecret != "" {
		cfg.GoogleClientSecret = googleClientSecret
	}
	if initTable := os.Getenv("INIT_TABLE"); initTable != "" {
		cfg.InitTable = initTable == "true"
	}
	if enableCache := os.Getenv("ENABLE_CACHE"); enableCache != "" {
		cfg.EnableCache = enableCache != "false"
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	return &cfg, nil
}

func (c *Config) GetDSN() string {
	if c.DatabaseURL == "" {
		return ""
	}

	// PostgreSQL DSN format: postgres://user:password@host:port/dbname?sslmode=disable
	return c.DatabaseURL
}

func (c *Config) IsLoginEnabled() bool {
	hasGitHubOAuth := c.GClientID != "" && c.GClientSecret != ""
	hasGoogleOAuth := c.GoogleClientID != "" && c.GoogleClientSecret != ""
	return (hasGitHubOAuth || hasGoogleOAuth) && c.JWTSecret != ""
}

func (c *Config) IsGitHubOAuthEnabled() bool {
	return c.GClientID != "" && c.GClientSecret != ""
}

func (c *Config) IsGoogleOAuthEnabled() bool {
	return c.GoogleClientID != "" && c.GoogleClientSecret != ""
}

func (c *Config) GetServerAddr() string {
	if c.Port == 0 {
		c.Port = 8080
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) GetOAuthRedirectURL() string {
	return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s", c.GClientID)
}

func (c *Config) GetAllowedEmails() []string {
	return c.AllowedEmails
}
