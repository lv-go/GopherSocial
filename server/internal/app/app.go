package app

import (
	"context"
	"expvar"
	"runtime"

	"github.com/go-redis/redis/v8"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/db"
	"github.com/sikozonpc/social/internal/handlers"
	"github.com/sikozonpc/social/internal/mailer"
	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/internal/store"
	"github.com/sikozonpc/social/internal/store/cache"
	"go.uber.org/zap"
)

var (
	Config                 config.Config
	Store                  store.Storage
	CacheStorage           cache.Storage
	Logger                 *zap.SugaredLogger
	Mailer                 mailer.Client
	Authenticator          auth.Authenticator
	AuthHandlers           auth.Handlers
	AuthMiddlewares        auth.Middlewares
	RateLimiter            ratelimiter.RateLimiter
	UsersHandlers          handlers.UsersHandlers
	PostsHandlers          handlers.PostsHandlers
	HealthCheckHandlers    handlers.HealthCheckHandler
	FeedHandlers           handlers.FeedsHandlers
	RateLimiterMiddlewares ratelimiter.Middlewares
	UsersRepository        *repositories.UsersRepository
	PostRepository         *repositories.PostsRepository
	CommentsRepository     *repositories.CommentsRepository
	FollowersRepository    *repositories.FollowersRepository
	RolesRepository        *repositories.RolesRepository
)

func Setup(ctx context.Context) {
	Config = config.Setup()

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

	go func() {
		<-ctx.Done()
		Logger.Info("closing database connection pool")
		err := _db.Close()
		if err != nil {
			Logger.Error(err)
		}
	}()
	Logger.Info("database connection pool established")

	// Cache
	var rdb *redis.Client
	if Config.RedisCfg.Enabled {
		rdb = cache.NewRedisClient(Config.RedisCfg.Addr, Config.RedisCfg.Pw, Config.RedisCfg.Db)
		Logger.Info("redis cache connection established")

		go func() {
			<-ctx.Done()
			Logger.Info("closing redis cache connection")
			if err := rdb.Close(); err != nil {
				Logger.Error("error closing redis cache connection", "error", err)
			}
		}()
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

	RateLimiterMiddlewares = ratelimiter.NewMiddlewares(Config.RateLimiter, RateLimiter)
	// Handlers
	AuthHandlers = auth.NewHandlers(
		Store,
		UsersRepository,
		Mailer,
		Logger,
		Config,
		Authenticator,
	)
	AuthMiddlewares = auth.NewMiddlewares(
		Store,
		CacheStorage,
		UsersRepository,
		RateLimiter,
		Authenticator,
		Logger,
		Config,
	)

	HealthCheckHandlers = handlers.NewHealthCheckHandler(Config)
	UsersHandlers = handlers.NewUsersHandlers(AuthMiddlewares, UsersRepository, FollowersRepository)
	PostsHandlers = handlers.NewPostsHandlers(UsersRepository, PostRepository, CommentsRepository, RolesRepository)
	FeedHandlers = handlers.NewFeedHandlers(Store)

	// Metrics collected
	expvar.NewString("version").Set(Config.Version)
	expvar.Publish("database", expvar.Func(func() any {
		return _db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
}
