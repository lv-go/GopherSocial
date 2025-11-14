package handlers

import (
	"net/http"

	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/utils"
)

type HealthCheckHandler struct {
	config config.Config
}

func NewHealthCheckHandler(config config.Config) HealthCheckHandler {
	return HealthCheckHandler{
		config: config,
	}
}

// HealthCheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (h *HealthCheckHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     h.config.Env,
		"version": h.config.Version,
	}

	utils.WriteJSON(w, http.StatusOK, data)
}
