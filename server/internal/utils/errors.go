package utils

import (
	"log/slog"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONMessage(w, http.StatusInternalServerError, "the server encountered a problem")
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	slog.Warn("forbidden", "method", r.Method, "path", r.URL.Path)

	WriteJSONMessage(w, http.StatusForbidden, "forbidden")
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONMessage(w, http.StatusBadRequest, err.Error())
}

func ConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONMessage(w, http.StatusConflict, err.Error())
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONMessage(w, http.StatusNotFound, "not found")
}

func UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJSONMessage(w, http.StatusUnauthorized, "unauthorized")
}

func UnauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	WriteJSONMessage(w, http.StatusUnauthorized, "unauthorized")
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	slog.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	WriteJSONMessage(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
