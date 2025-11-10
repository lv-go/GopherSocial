package app

import (
	"expvar"
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/go-redis/redis/v8"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/db"
	"github.com/sikozonpc/social/internal/mailer"
	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/sikozonpc/social/internal/store"
	"github.com/sikozonpc/social/internal/store/cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const version = "1.1.0"

func Main() {

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

	// Unmarshal the config into the config struct
	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	slog.Debug("config loaded", "config", cfg)

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Main Database
	db, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	// Cache
	var rdb *redis.Client
	if cfg.RedisCfg.Enabled {
		rdb = cache.NewRedisClient(cfg.RedisCfg.Addr, cfg.RedisCfg.Pw, cfg.RedisCfg.Db)
		logger.Info("redis cache connection established")

		defer rdb.Close()
	}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.RateLimiter.RequestsPerTimeFrame,
		cfg.RateLimiter.TimeFrame,
	)

	// Mailer
	// mailer := mailer.NewSendgrid(cfg.Mail.SendGrid.ApiKey, cfg.Mail.FromEmail)
	mailtrap, err := mailer.NewMailTrapClient(cfg.Mail.MailTrap.ApiKey, cfg.Mail.FromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.Auth.Token.Secret,
		cfg.Auth.Token.Iss,
		cfg.Auth.Token.Iss,
	)

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailtrap,
		authenticator: jwtAuthenticator,
		rateLimiter:   rateLimiter,
	}

	// Metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
