package transport

import (
	"fmt"
	"net/http"

	"github.com/FlashpointProject/CommunityWebsite/types"
	"github.com/FlashpointProject/CommunityWebsite/utils"
	"github.com/gorilla/schema"
)

func (a *App) SearchGotdSuggestions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		writeError(ctx, w, perr("failed to parse form", http.StatusBadRequest))
		return
	}

	var query types.GotdSuggestionsSearchQuery
	err = schema.NewDecoder().Decode(&query, r.Form)
	if err != nil {
		writeError(ctx, w, perr(fmt.Sprintf("failed to decode form: %s", err.Error()), http.StatusBadRequest))
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	suggestions, total, err := a.Service.SearchGotdSuggestions(ctx, &query, a.Fpfss)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, err)
		return
	}

	res := &types.GotdSuggestionsSearchResponse{
		Suggestions: suggestions,
		Total:       total,
	}

	writeResponse(ctx, w, res, http.StatusOK)
}
