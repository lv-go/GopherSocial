package handlers

import (
	"net/http"
	"strconv"

	"github.com/sikozonpc/social/internal/auth"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/internal/utils"
)

type FeedsHandlers struct {
	postsRepository *repositories.PostsRepository
}

func NewFeedHandlers(postsRepository *repositories.PostsRepository) FeedsHandlers {
	return FeedsHandlers{postsRepository: postsRepository}
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
	pageSizeStr := r.URL.Query().Get("size")
	pageNumberStr := r.URL.Query().Get("number")

	var pageSize, pageNumber int
	if pageSizeStr == "" || pageNumberStr == "" {
		pageSize = 20
		pageNumber = 1
	} else {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			utils.BadRequestResponse(w, r, err)
			return
		}

		pageNumber, err = strconv.Atoi(pageNumberStr)
		if err != nil {
			utils.BadRequestResponse(w, r, err)
			return
		}
	}

	pageQuery := repositories.PageQuery{
		Size:   pageSize,
		Number: pageNumber,
	}
	if err := utils.Validate.Struct(pageQuery); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	user := auth.GetUserFromContext(r)

	feed, err := fh.postsRepository.GetPageByUserID(ctx, user.ID, pageQuery)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, feed)
}
