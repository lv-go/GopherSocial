package auth

import (
	"net/http"

	"github.com/sikozonpc/social/internal/store"
)

type userKey string

const userCtx userKey = "User"

func GetUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
