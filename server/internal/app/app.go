package app

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/handlers"
	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/sikozonpc/social/internal/repositories"
	"go.uber.org/zap"
)

var (
	Config                 config.Config
	Logger                 *zap.SugaredLogger
	AuthClient             auth.Client
	AuthMiddlewares        auth.Middlewares
	RateLimiter            ratelimiter.RateLimiter
	UsersHandlers          handlers.UsersHandlers
	PostsHandlers          handlers.PostsHandlers
	HealthCheckHandlers    handlers.HealthCheckHandler
	FeedHandlers           handlers.FeedsHandlers
	RateLimiterMiddlewares ratelimiter.Middlewares
	PostRepository         *repositories.PostsRepository
	CommentsRepository     *repositories.CommentsRepository
	FollowersRepository    *repositories.FollowersRepository
	RolesRepository        *repositories.RolesRepository
)

func Setup(ctx context.Context, configPath string) {
	Config = config.Setup(configPath)

	// Auth Client
	AuthClient = auth.Setup(ctx)

	// Logger
	Logger = zap.Must(zap.NewProduction()).Sugar()
	defer Logger.Sync()

	// Cache
	var rdb *redis.Client
	if Config.RedisCfg.Enabled {
		rdb = repositories.SetupRedisClient(Config.RedisCfg)
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

	repositories.SetupGormDB(Config.GormDBConfig)

	// Repositories
	PostRepository = repositories.NewPostsRepository()
	CommentsRepository = repositories.NewCommentsRepository()
	FollowersRepository = repositories.NewFollowersRepository()
	RolesRepository = repositories.NewRolesRepository()

	// Middlewares
	RateLimiterMiddlewares = ratelimiter.NewMiddlewares(Config.RateLimiter, RateLimiter)
	AuthMiddlewares = auth.NewMiddlewares(
		RateLimiter,
		AuthClient,
		Logger,
		Config,
	)

	// Handlers
	HealthCheckHandlers = handlers.NewHealthCheckHandler(Config)
	UsersHandlers = handlers.NewUsersHandlers(AuthMiddlewares, FollowersRepository)
	PostsHandlers = handlers.NewPostsHandlers(PostRepository, CommentsRepository, RolesRepository)
	FeedHandlers = handlers.NewFeedHandlers(PostRepository)

	// Metrics collected
	//expvar.NewString("version").Set(Config.Version)
	//expvar.Publish("goroutines", expvar.Func(func() any {
	//	return runtime.NumGoroutine()
	//}))
}
