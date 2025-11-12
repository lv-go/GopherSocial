package handlers

import (
	"net/http"

	"github.com/sikozonpc/social/internal/app"
	"github.com/sikozonpc/social/internal/utils"
)

// HealthCheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.Config.Env,
		"version": app.Config.Version,
	}

	utils.WriteJSON(w, http.StatusOK, data)
}
