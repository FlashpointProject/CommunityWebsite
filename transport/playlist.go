package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FlashpointProject/CommunityWebsite/constants"
	"github.com/FlashpointProject/CommunityWebsite/types"
	"github.com/FlashpointProject/CommunityWebsite/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func (a *App) SearchPlaylists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		writeError(ctx, w, perr("failed to parse form", http.StatusBadRequest))
		return
	}

	var query types.PlaylistSearchQuery
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

	playlists, total, err := a.Service.SearchPlaylists(ctx, &query)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, err)
		return
	}
	playlistInfos := make([]*types.PlaylistInfo, len(playlists))
	for i, playlist := range playlists {
		playlistInfos[i] = &types.PlaylistInfo{
			ID:           playlist.ID,
			Name:         playlist.Name,
			TotalGames:   playlist.TotalGames,
			Description:  playlist.Description,
			Author:       playlist.Author,
			Library:      playlist.Library,
			Icon:         playlist.Icon,
			Public:       playlist.Public,
			Extreme:      playlist.Extreme,
			FilterGroups: playlist.FilterGroups,
			CreatedAt:    playlist.CreatedAt,
			UpdatedAt:    playlist.UpdatedAt,
		}
	}

	res := &types.PlaylistSearchResponse{
		Playlists: playlistInfos,
		Total:     total,
	}

	writeResponse(ctx, w, res, http.StatusOK)
}

func (a *App) GetRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writeResponse(ctx, w, a.Service.RoleCache, http.StatusOK)
}

func (a *App) GetFilterGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writeResponse(ctx, w, &constants.ResponseFilterGroups{Groups: constants.GetFilterGroups()}, http.StatusOK)
}

func (a *App) GetUserSuggestions(w http.ResponseWriter, r *http.Request) {

}

func (a *App) GetPlaylistPreview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	idStr := params[constants.ResourceKeyPlaylistID]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(ctx, w, perr("invalid playlist id", http.StatusBadRequest))
		return
	}

	playlist, err := a.Service.GetDownloadablePlaylist(ctx, id)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to get playlist", http.StatusInternalServerError))
		return
	}

	if playlist == nil {
		writeError(ctx, w, perr("playlist not found", http.StatusNotFound))
		return
	}

	writeResponse(ctx, w, playlist, http.StatusOK)
}

func (a *App) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	idStr := params[constants.ResourceKeyPlaylistID]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(ctx, w, perr("invalid playlist id", http.StatusBadRequest))
		return
	}

	playlist, err := a.Service.GetPlaylist(ctx, id, a.Fpfss)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to get playlist", http.StatusInternalServerError))
		return
	}

	if playlist == nil {
		writeError(ctx, w, perr("playlist not found", http.StatusNotFound))
		return
	}

	writeResponse(ctx, w, playlist, http.StatusOK)
}

func (a *App) DownloadPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	idStr := params[constants.ResourceKeyPlaylistID]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(ctx, w, perr("invalid playlist id", http.StatusBadRequest))
		return
	}

	playlist, err := a.Service.GetDownloadablePlaylist(ctx, id)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to get playlist", http.StatusInternalServerError))
		return
	}

	if playlist == nil {
		writeError(ctx, w, perr("playlist not found", http.StatusNotFound))
		return
	}

	launcherPlaylist := &types.LauncherPlaylist{
		ID:          fmt.Sprintf("fpcommunity-%d", playlist.ID),
		Title:       playlist.Name,
		Author:      playlist.Author.Username,
		Description: playlist.Description,
		Library:     playlist.Library,
		Icon:        playlist.Icon,
		Extreme:     !playlist.Public,
		Games:       playlist.Games,
	}
	data, err := json.Marshal(launcherPlaylist)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to marshal playlist", http.StatusInternalServerError))
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", playlist.Name))
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

func (a *App) SubmitPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)
	var playlist types.LauncherPlaylist

	// Decode the JSON request body into the playlist struct
	err := json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		writeError(ctx, w, perr("failed to decode request body - "+err.Error(), http.StatusBadRequest))
		return
	}

	// Validate fields
	valid := true
	valid = valid && playlist.Title != ""
	valid = valid && playlist.Description != ""
	valid = valid && playlist.Library != ""
	valid = valid && playlist.Games != nil
	if !valid {
		http.Error(w, "title / description / library / games are required playlist fields", http.StatusBadRequest)
		return
	}
	valid = valid && len(playlist.Games) > 0
	if !valid {
		http.Error(w, "playlist must contain at least one game", http.StatusBadRequest)
		return
	}

	parsedPlaylist := &types.Playlist{
		ID:          0,
		Name:        playlist.Title,
		Description: playlist.Description,
		Library:     playlist.Library,
		Icon:        playlist.Icon,
		TotalGames:  len(playlist.Games),
		Games:       playlist.Games,
		Public:      true,
	}

	err = a.Service.SubmitPlaylist(ctx, uid, parsedPlaylist, a.Fpfss)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(parsedPlaylist)
}

func (a *App) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)
	params := mux.Vars(r)
	idStr := params[constants.ResourceKeyPlaylistID]
	var playlist types.LauncherPlaylist

	// Decode the JSON request body into the playlist struct
	err := json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		writeError(ctx, w, perr("failed to decode request body - "+err.Error(), http.StatusBadRequest))
		return
	}

	// Validate fields
	valid := true
	valid = valid && playlist.Title != ""
	valid = valid && playlist.Description != ""
	valid = valid && playlist.Library != ""
	valid = valid && playlist.Games != nil
	if !valid {
		http.Error(w, "title / description / library / games are required playlist fields", http.StatusBadRequest)
		return
	}
	valid = valid && len(playlist.Games) > 0
	if !valid {
		http.Error(w, "playlist must contain at least one game", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(ctx, w, perr("invalid post id", http.StatusBadRequest))
		return
	}

	parsedPlaylist := &types.Playlist{
		ID:          id,
		Name:        playlist.Title,
		Description: playlist.Description,
		Library:     playlist.Library,
		Icon:        playlist.Icon,
		TotalGames:  len(playlist.Games),
		Games:       playlist.Games,
		Public:      true,
	}

	newPlaylist, err := a.Service.UpdatePlaylist(ctx, uid, parsedPlaylist, a.Fpfss)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newPlaylist)
}

func (a *App) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)
	params := mux.Vars(r)
	idStr := params[constants.ResourceKeyPlaylistID]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(ctx, w, perr("invalid post id", http.StatusBadRequest))
		return
	}

	err = a.Service.DeletePlaylist(ctx, uid, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response back
	writeResponse(ctx, w, nil, http.StatusOK)
}

func (a *App) GetNewsPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	idStr := params[constants.ResourceKeyPostID]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(ctx, w, perr("invalid post id", http.StatusBadRequest))
		return
	}

	post, err := a.Service.GetNewsPost(ctx, id)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to get post", http.StatusInternalServerError))
		return
	}

	if post == nil {
		writeError(ctx, w, perr("post not found", http.StatusNotFound))
		return
	}

	writeResponse(ctx, w, post, http.StatusOK)
}

func (a *App) SearchNewsPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		writeError(ctx, w, perr("failed to parse form", http.StatusBadRequest))
		return
	}

	var query types.NewsPostSearchQuery
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

	posts, total, err := a.Service.SearchNewsPosts(ctx, &query)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, err)
		return
	}

	res := &types.NewsPostSearchResponse{
		Posts: posts,
		Total: total,
	}

	writeResponse(ctx, w, res, http.StatusOK)
}

func (a *App) SubmitNewsPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)
	var subNewsPost types.SubmittedNewsPost

	// Decode the JSON request body into the playlist struct
	err := json.NewDecoder(r.Body).Decode(&subNewsPost)
	if err != nil {
		writeError(ctx, w, perr("failed to decode request body - "+err.Error(), http.StatusBadRequest))
		return
	}

	// Validate fields
	if subNewsPost.Title == "" {
		http.Error(w, "title is a required field", http.StatusBadRequest)
		return
	}
	if subNewsPost.Content == "" {
		http.Error(w, "content is a required field", http.StatusBadRequest)
		return
	}
	if subNewsPost.PostType == "" {
		http.Error(w, "post_type is a required field", http.StatusBadRequest)
		return
	}

	newsPost := &types.NewsPost{
		PostType: subNewsPost.PostType,
		Title:    subNewsPost.Title,
		Content:  subNewsPost.Content,
	}

	err = a.Service.SubmitNewsPost(ctx, uid, newsPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, w, newsPost, http.StatusOK)
}

func (a *App) SearchContentReports(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		writeError(ctx, w, perr("failed to parse form", http.StatusBadRequest))
		return
	}

	var query types.ContentReportSearchQuery
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

	reports, total, err := a.Service.SearchContentReports(ctx, &query)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, err)
		return
	}

	res := &types.ContentReportSearchResponse{
		Reports: reports,
		Total:   total,
	}

	writeResponse(ctx, w, res, http.StatusOK)
}

func (a *App) SubmitContentReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := utils.UserID(ctx)
	var subReport types.SubmittedContentReport

	// Decode the JSON request body into the playlist struct
	err := json.NewDecoder(r.Body).Decode(&subReport)
	if err != nil {
		writeError(ctx, w, perr("failed to decode request body - "+err.Error(), http.StatusBadRequest))
		return
	}

	// Validate fields
	if subReport.ContentRef == "" {
		http.Error(w, "content_ref is a required field", http.StatusBadRequest)
		return
	}
	if subReport.ReportReason == "" {
		http.Error(w, "report_reason is a required field", http.StatusBadRequest)
		return
	}

	// Find referenced content
	contentRefs, err := utils.ParseContentRef(subReport.ContentRef)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reportedUser := ""
	switch contentRefs[0].ContentType {
	case "playlist":
		playlistId, err := strconv.Atoi(contentRefs[0].ContentID)
		if err != nil {
			writeError(ctx, w, perr("invalid playlist id", http.StatusBadRequest))
			return
		}
		playlist, err := a.Service.GetPlaylist(ctx, int64(playlistId), a.Fpfss)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if playlist == nil {
			writeError(ctx, w, perr("invalid playlist id", http.StatusBadRequest))
			return
		}
		reportedUser = playlist.Author.UserID
	}

	report := &types.ContentReport{
		ContentRef:   subReport.ContentRef,
		ReportReason: subReport.ReportReason,
		ReportedBy: &types.UserProfile{
			UserID: uid,
		},
		ReportedUser: &types.UserProfile{
			UserID: reportedUser,
		},
		ReportState: "reported",
		ActionTaken: "",
		ResolvedBy: &types.UserProfile{
			UserID: "",
		},
	}

	err = a.Service.SubmitContentReport(ctx, report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, w, report, http.StatusOK)
}

func (a *App) GetGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params[constants.ResourceKeyGameID]

	if id == "" {
		writeError(ctx, w, perr("invalid game id", http.StatusBadRequest))
		return
	}
	if len(id) != 36 {
		writeError(ctx, w, perr("invalid game id", http.StatusBadRequest))
		return
	}

	game, err := a.Service.GetGame(ctx, id, a.Fpfss)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		writeError(ctx, w, perr("failed to get game", http.StatusInternalServerError))
		return
	}

	if game == nil {
		writeError(ctx, w, perr("game not found", http.StatusNotFound))
		return
	}

	writeResponse(ctx, w, game, http.StatusOK)
}
