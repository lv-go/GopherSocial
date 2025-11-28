package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/internal/store"
	"github.com/sikozonpc/social/internal/utils"
)

type UsersHandlers struct {
	authMiddlewares     auth.Middlewares
	usersRepository     *repositories.UsersRepository
	followersRepository *repositories.FollowersRepository
}

func NewUsersHandlers(
	authMiddlewares auth.Middlewares,
	usersRepository *repositories.UsersRepository,
	followersRepository *repositories.FollowersRepository,
) UsersHandlers {
	return UsersHandlers{
		authMiddlewares:     authMiddlewares,
		usersRepository:     usersRepository,
		followersRepository: followersRepository,
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
	userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.usersRepository.GetByID(r.Context(), uint(userID))
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

	if err := h.followersRepository.Follow(ctx, uint(followedID), followerUser.ID); err != nil {
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

	unfollowedID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := h.followersRepository.Unfollow(ctx, uint(unfollowedID), followerUser.ID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
