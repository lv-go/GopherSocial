package app

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/sikozonpc/social/docs" // This is required to generate swagger docs
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/env"
	"github.com/sikozonpc/social/internal/mailer"
	"github.com/sikozonpc/social/internal/ratelimiter"
	"github.com/sikozonpc/social/internal/store"
	"github.com/sikozonpc/social/internal/store/cache"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	Addr        string
	DB          dbConfig
	Env         string
	ApiURL      string
	Mail        mailConfig
	FrontendURL string
	Auth        authConfig
	RedisCfg    redisConfig
	RateLimiter ratelimiter.Config
}

type redisConfig struct {
	Addr    string
	Pw      string
	Db      int
	Enabled bool
}

type authConfig struct {
	Basic basicConfig
	Token tokenConfig
}

type tokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
}

type basicConfig struct {
	User string
	Pass string
}

type mailConfig struct {
	SendGrid  sendGridConfig
	MailTrap  mailTrapConfig
	FromEmail string
	Exp       time.Duration
}

type mailTrapConfig struct {
	ApiKey string
}

type sendGridConfig struct {
	ApiKey string
}

type dbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5173")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	if app.config.RateLimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		// Operations
		r.Get("/health", app.healthCheckHandler)
		r.With(app.BasicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.Addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)

				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.ApiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "Addr", app.config.Addr, "Env", app.config.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "Addr", app.config.Addr, "Env", app.config.Env)

	return nil
}
