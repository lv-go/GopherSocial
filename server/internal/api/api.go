package api

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
	"github.com/sikozonpc/social/internal/app"
	auth2 "github.com/sikozonpc/social/internal/auth/auth"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/handlers"
	middleware2 "github.com/sikozonpc/social/internal/ratelimiter/middleware"
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

type Application struct {
	Config        config.Config
	Store         store.Storage
	CacheStorage  cache.Storage
	Logger        *zap.SugaredLogger
	Mailer        mailer.Client
	Authenticator auth.Authenticator
	RateLimiter   ratelimiter.Limiter
}

func Mount() http.Handler {
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

	if app.Config.RateLimiter.Enabled {
		r.Use(middleware2.RateLimiterMiddleware)
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		// Operations
		r.Get("/health", handlers.HealthCheckHandler)
		r.With(auth2.BasicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.Config.Addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(auth2.AuthTokenMiddleware)
			r.Post("/", handlers.CreatePostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(handlers.PostsContextMiddleware)
				r.Get("/", handlers.GetPostHandler)

				r.Patch("/", handlers.CheckPostOwnership("moderator", handlers.UpdatePostHandler))
				r.Delete("/", handlers.CheckPostOwnership("admin", handlers.DeletePostHandler))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", handlers.ActivateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(auth2.AuthTokenMiddleware)

				r.Get("/", handlers.GetUserHandler)
				r.Put("/follow", handlers.FollowUserHandler)
				r.Put("/unfollow", handlers.UnfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(auth2.AuthTokenMiddleware)
				r.Get("/feed", handlers.GetUserFeedHandler)
			})
		})

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/User", auth2.RegisterUserHandler)
			r.Post("/token", auth2.CreateTokenHandler)
		})
	})

	return r
}

func Run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = app.Config.Version
	docs.SwaggerInfo.Host = app.Config.ApiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.Config.Addr,
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

		app.Logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.Logger.Infow("server has started", "addr", app.Config.Addr, "env", app.Config.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.Logger.Infow("server has stopped", "addr", app.Config.Addr, "env", app.Config.Env)

	return nil
}
