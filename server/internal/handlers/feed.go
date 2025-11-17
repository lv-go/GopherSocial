package handlers

import (
	"net/http"

	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/internal/store"
	"github.com/sikozonpc/social/internal/utils"
)

type FeedsHandlers struct {
	store           store.Storage
	postsRepository repositories.PostsRepository
}

func NewFeedHandlers(store store.Storage) FeedsHandlers {
	return FeedsHandlers{store: store}
}

// GetUserFeedHandler godoc
//
//	@Summary		Fetches the User feed
//	@Description	Fetches the User feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Size"
//	@Param			offset	query		int		false	"Number"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (fh *FeedsHandlers) GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
		Search: "",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.Validate.Struct(fq); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	user := auth.GetUserFromContext(r)

	feed, err := fh.postsRepository.GetPageByUserID(ctx, user.ID, repositories.PageQuery{
		Size:   20,
		Number: 1,
	})
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, feed)
}
