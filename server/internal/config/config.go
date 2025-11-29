package config

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/spf13/viper"
)

type Config struct {
	Addr         string
	Db           DbConfig
	Env          string
	Version      string
	ApiURL       string
	Mail         MailConfig
	FrontendURL  string
	Auth         AuthConfig
	RedisCfg     repositories.RedisConfig
	RateLimiter  ratelimiter.Config
	GormDBConfig repositories.GormDBConfig
}

type AuthConfig struct {
	Basic BasicConfig
	Token TokenConfig
}

type TokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
}

type BasicConfig struct {
	User string
	Pass string
}

type MailConfig struct {
	SendGrid  SendGridConfig
	MailTrap  MailTrapConfig
	FromEmail string
	Exp       time.Duration
}

type MailTrapConfig struct {
	ApiKey string
}

type SendGridConfig struct {
	ApiKey string
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func Setup() Config {
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "" {
		appEnv = "." + appEnv
	}
	if appEnv == "" || appEnv == "dev" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		viper.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	// Set the name of the config file (without extension)
	viper.SetConfigName("config" + appEnv)
	// Set the type of the config file
	viper.SetConfigType("yaml")
	// Add the path where Viper should look for the config file
	viper.AddConfigPath(".")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var cfg Config
	// Unmarshal the config into the config struct
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	//slog.Debug("config loaded", "config", cfg)
	return cfg
}
