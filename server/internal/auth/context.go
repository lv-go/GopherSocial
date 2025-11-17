package auth

import (
	"context"
	"net/http"

	"github.com/sikozonpc/social/internal/models"
)

type userKey string

const userCtx userKey = "User"

func GetUserFromContext(r *http.Request) *models.User {
	user, _ := r.Context().Value(userCtx).(*models.User)
	return user
}

func SetUserInContext(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userCtx, user)
	return r.WithContext(ctx)
}
