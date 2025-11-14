package auth

import (
	"context"
	"net/http"

	"github.com/sikozonpc/social/internal/store"
)

type userKey string

const userCtx userKey = "User"

func GetUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

func SetUserInContext(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), userCtx, user)
	return r.WithContext(ctx)
}
