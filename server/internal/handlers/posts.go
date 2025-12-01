package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/models"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/internal/utils"
	"gorm.io/gorm"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type PostsHandlers struct {
	postsRepository    *repositories.PostsRepository
	commentsRepository *repositories.CommentsRepository
	rolesRepository    *repositories.RolesRepository
}

func NewPostsHandlers(
	postsRepository *repositories.PostsRepository,
	commentsRepository *repositories.CommentsRepository,
	rolesRepository *repositories.RolesRepository,
) PostsHandlers {
	return PostsHandlers{
		postsRepository:    postsRepository,
		commentsRepository: commentsRepository,
		rolesRepository:    rolesRepository,
	}
}

// CreatePostHandler godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (h *PostsHandlers) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	user := auth.GetUserFromContext(r)

	post := &models.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}

	ctx := r.Context()

	if err := h.postsRepository.Create(ctx, post); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, post)
}

// GetPostHandler godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	store.Post
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (h *PostsHandlers) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := h.commentsRepository.GetByPostID(r.Context(), post.ID)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	post.Comments = comments

	utils.WriteJSON(w, http.StatusOK, post)
}

// DeletePostHandler godoc
//
//	@Summary		Deletes a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object} string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (h *PostsHandlers) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := h.postsRepository.DeleteByID(ctx, uint(id)); err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.NotFoundResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

// UpdatePostHandler godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (h *PostsHandlers) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	ctx := r.Context()

	if err := h.postsRepository.UpdateByID(ctx, uint(id), post); err != nil {
		utils.InternalServerError(w, r, err)
	}

	utils.WriteJSON(w, http.StatusOK, post)
}

func (h *PostsHandlers) PostsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			utils.InternalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := h.postsRepository.GetByID(ctx, uint(id))
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				utils.NotFoundResponse(w, r, err)
			default:
				utils.InternalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *models.Post {
	post, _ := r.Context().Value(postCtx).(*models.Post)
	return post
}

func (h *PostsHandlers) CheckPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r)
		post := getPostFromCtx(r)

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := h.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			utils.InternalServerError(w, r, err)
			return
		}

		if !allowed {
			utils.ForbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *PostsHandlers) checkRolePrecedence(ctx context.Context, user *models.User, roleName string) (bool, error) {
	role, err := h.rolesRepository.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
