package app

import (
	"expvar"
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/go-redis/redis/v8"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/db"
	"github.com/sikozonpc/social/internal/mailer"
	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/sikozonpc/social/internal/store"
	"github.com/sikozonpc/social/internal/store/cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Config        config.Config
	Store         store.Storage
	CacheStorage  cache.Storage
	Logger        *zap.SugaredLogger
	Mailer        mailer.Client
	Authenticator auth.Authenticator
	RateLimiter   ratelimiter.Limiter
)

func Setup() {
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
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	slog.Debug("config loaded", "config", Config)

	// Logger
	Logger = zap.Must(zap.NewProduction()).Sugar()
	defer Logger.Sync()

	// Main Database
	_db, err := db.New(
		Config.Db.Addr,
		Config.Db.MaxOpenConns,
		Config.Db.MaxIdleConns,
		Config.Db.MaxIdleTime,
	)
	if err != nil {
		Logger.Fatal(err)
	}

	defer _db.Close()
	Logger.Info("database connection pool established")

	// Cache
	var rdb *redis.Client
	if Config.RedisCfg.Enabled {
		rdb = cache.NewRedisClient(Config.RedisCfg.Addr, Config.RedisCfg.Pw, Config.RedisCfg.Db)
		Logger.Info("redis cache connection established")

		defer rdb.Close()
	}

	// Rate limiter
	RateLimiter = ratelimiter.NewFixedWindowLimiter(
		Config.RateLimiter.RequestsPerTimeFrame,
		Config.RateLimiter.TimeFrame,
	)

	// Mailer
	// mailer := mailer.NewSendgrid(cfg.mail.SendGrid.ApiKey, cfg.mail.FromEmail)
	Mailer, err = mailer.NewMailTrapClient(Config.Mail.MailTrap.ApiKey, Config.Mail.FromEmail)
	if err != nil {
		Logger.Fatal(err)
	}

	// Authenticator
	Authenticator = auth.NewJWTAuthenticator(
		Config.Auth.Token.Secret,
		Config.Auth.Token.Iss,
		Config.Auth.Token.Iss,
	)

	Store = store.NewStorage(_db)
	CacheStorage = cache.NewRedisStorage(rdb)

	// Metrics collected
	expvar.NewString("version").Set(Config.Version)
	expvar.Publish("database", expvar.Func(func() any {
		return _db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

}
