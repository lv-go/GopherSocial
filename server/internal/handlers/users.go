package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/internal/utils"
	"gorm.io/gorm"
)

type UsersHandlers struct {
	authMiddlewares     auth.Middlewares
	followersRepository *repositories.FollowersRepository
}

func NewUsersHandlers(
	authMiddlewares auth.Middlewares,
	followersRepository *repositories.FollowersRepository,
) UsersHandlers {
	return UsersHandlers{
		authMiddlewares:     authMiddlewares,
		followersRepository: followersRepository,
	}
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

	followedID := uuid.NewSHA1(uuid.NameSpaceURL, []byte(chi.URLParam(r, "userID")))

	ctx := r.Context()

	if err := h.followersRepository.Follow(ctx, followedID, followerUser.ID); err != nil {
		switch {
		case errors.Is(err, gorm.ErrCheckConstraintViolated):
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

	unfollowedID := uuid.NewSHA1(uuid.NameSpaceURL, []byte(chi.URLParam(r, "userID")))
	ctx := r.Context()

	if err := h.followersRepository.Unfollow(ctx, unfollowedID, followerUser.ID); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
