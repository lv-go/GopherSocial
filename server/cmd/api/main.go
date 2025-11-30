package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/sikozonpc/social/internal/api"
	app "github.com/sikozonpc/social/internal/app"
)

const version = "1.1.0"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gohpers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Setup(ctx, "")

	mux := api.Mount()

	slog.Info("Shutting down server", "error", api.Run(mux))
}
