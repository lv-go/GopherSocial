package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/store"
	"github.com/sikozonpc/social/internal/utils"
)

type UsersHandlers struct {
	authMiddlewares auth.Middlewares
	store           store.Storage
}

func NewUsersHandlers(
	authMiddlewares auth.Middlewares,
	store store.Storage,
) UsersHandlers {
	return UsersHandlers{
		authMiddlewares: authMiddlewares,
		store:           store,
	}
}

// GetUserHandler godoc
//
//	@Summary		Fetches a User profile
//	@Description	Fetches a User profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (h *UsersHandlers) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.authMiddlewares.GetUser(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			utils.NotFoundResponse(w, r, err)
			return
		default:
			utils.InternalServerError(w, r, err)
			return
		}
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

// FollowUserHandler godoc
//
//	@Summary		Follows a User
//	@Description	Follows a User by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (h *UsersHandlers) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := auth.GetUserFromContext(r)

	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := h.store.Followers.Follow(ctx, followedID, followerUser.ID); err != nil {
		switch err {
		case store.ErrConflict:
			utils.ConflictResponse(w, r, err)
			return
		default:
			utils.InternalServerError(w, r, err)
			return
		}
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

// UnfollowUserHandler gdoc
//
//	@Summary		Unfollow a User
//	@Description	Unfollow a User by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (h *UsersHandlers) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := auth.GetUserFromContext(r)

	unfollowedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := h.store.Followers.Unfollow(ctx, followerUser.ID, unfollowedID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

// ActivateUserHandler godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (h *UsersHandlers) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := h.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			utils.NotFoundResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, "")
}
