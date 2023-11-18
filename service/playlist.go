package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FlashpointProject/CommunityWebsite/constants"
	"github.com/FlashpointProject/CommunityWebsite/database"
	"github.com/FlashpointProject/CommunityWebsite/types"
	"github.com/FlashpointProject/CommunityWebsite/utils"
)

func (s *Service) SearchPlaylists(ctx context.Context, searchOpts *types.PlaylistSearchQuery) ([]*types.Playlist, int64, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, 0, dberr(err)
	}
	defer dbs.Rollback()

	playlists, total, err := s.pgdal.SearchPlaylists(dbs, searchOpts)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, 0, dberr(err)
	}

	return playlists, total, nil
}

func (s *Service) GetDownloadablePlaylist(ctx context.Context, id int64) (*types.Playlist, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}
	defer dbs.Rollback()

	playlist, err := s.pgdal.GetPlaylist(dbs, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	return playlist, nil
}

func (s *Service) GetPlaylist(ctx context.Context, id int64, fpfss types.IFpfss) (*types.FullPlaylist, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}
	defer dbs.Rollback()

	playlist, err := s.pgdal.GetPlaylist(dbs, id)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	if playlist == nil {
		return nil, nil
	}

	populatedPlaylist, err := s.fillPlaylist(dbs, playlist, fpfss)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	return populatedPlaylist, nil
}

func (s *Service) SubmitPlaylist(ctx context.Context, uid string, playlist *types.Playlist, fpfss types.IFpfss) error {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	defer dbs.Rollback()

	err = s.pgdal.SavePlaylist(dbs, uid, playlist, fpfss)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}

	err = dbs.Commit()
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}

	return nil
}

func (s *Service) fillPlaylist(dbs database.PGDBSession, playlist *types.Playlist, fpfss types.IFpfss) (*types.FullPlaylist, error) {
	fullPlaylist := &types.FullPlaylist{
		ID:           playlist.ID,
		Name:         playlist.Name,
		Description:  playlist.Description,
		Author:       playlist.Author,
		Library:      playlist.Library,
		Icon:         playlist.Icon,
		Games:        make([]types.GameWithNotes, len(playlist.Games)),
		Public:       playlist.Public,
		Extreme:      playlist.Extreme,
		FilterGroups: playlist.FilterGroups,
		CreatedAt:    playlist.CreatedAt,
		UpdatedAt:    playlist.UpdatedAt,
		TotalGames:   playlist.TotalGames,
	}

	gameIDs := make([]string, len(playlist.Games))
	for i, game := range playlist.Games {
		gameIDs[i] = game.GameID
		fullPlaylist.Games[i].GameID = game.GameID
		fullPlaylist.Games[i].Notes = game.Notes
	}
	loadedGames, err := s.pgdal.GetGames(dbs, gameIDs, fpfss)
	if err != nil {
		return nil, err
	}
	for _, game := range loadedGames {
		for i, gameWithNotes := range fullPlaylist.Games {
			if gameWithNotes.GameID == game.ID {
				fullPlaylist.Games[i].Game = game
			}
		}
	}

	return fullPlaylist, nil
}

func (s *Service) UpdatePlaylist(ctx context.Context, uid string, playlist *types.Playlist, fpfss types.IFpfss) (*types.Playlist, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}
	defer dbs.Rollback()

	existingPlaylist, err := s.pgdal.GetPlaylist(dbs, playlist.ID)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}
	if existingPlaylist == nil {
		return nil, fmt.Errorf("playlist not found")
	}

	if existingPlaylist.Author.UserID != uid {
		return nil, fmt.Errorf("user is not the author of this playlist")
	}

	existingPlaylist.Games = playlist.Games
	existingPlaylist.Name = playlist.Name
	existingPlaylist.Description = playlist.Description
	existingPlaylist.Library = playlist.Library
	existingPlaylist.Icon = playlist.Icon
	existingPlaylist.Public = playlist.Public

	err = s.pgdal.SavePlaylist(dbs, uid, playlist, fpfss)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	err = dbs.Commit()
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	return existingPlaylist, nil
}

func (s *Service) DeletePlaylist(ctx context.Context, uid string, id int64) error {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	defer dbs.Rollback()

	playlist, err := s.pgdal.GetPlaylist(dbs, id)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	if playlist == nil {
		return fmt.Errorf("playlist not found")
	}

	roles, err := s.pgdal.GetRoles(dbs)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	if !constants.IsModerator(roles) {
		if playlist.Author.UserID != uid {
			return fmt.Errorf("user is not the author of this playlist")
		}
	}

	err = s.pgdal.DeletePlaylist(dbs, id)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}

	err = dbs.Commit()
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}

	return nil
}

func (s *Service) SearchNewsPosts(ctx context.Context, query *types.NewsPostSearchQuery) ([]*types.NewsPost, int64, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, 0, dberr(err)
	}
	defer dbs.Rollback()

	posts, total, err := s.pgdal.SearchNewsPosts(dbs, query)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, 0, dberr(err)
	}

	return posts, total, nil
}

func (s *Service) GetNewsPost(ctx context.Context, id int64) (*types.NewsPost, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}
	defer dbs.Rollback()

	post, err := s.pgdal.GetNewsPost(dbs, id)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	return post, nil
}

func (s *Service) SubmitNewsPost(ctx context.Context, uid string, post *types.NewsPost) error {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	defer dbs.Rollback()

	err = s.pgdal.SaveNewsPost(dbs, uid, post)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}

	err = dbs.Commit()
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}

	return nil
}

func (s *Service) SearchContentReports(ctx context.Context, query *types.ContentReportSearchQuery) ([]*types.ContentReport, int64, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, 0, dberr(err)
	}
	defer dbs.Rollback()

	reports, total, err := s.pgdal.SearchContentReports(dbs, query)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, 0, dberr(err)
	}

	return reports, total, nil
}

func (s *Service) SubmitContentReport(ctx context.Context, report *types.ContentReport) error {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}
	defer dbs.Rollback()

	err = s.pgdal.SaveContentReport(dbs, report)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}

	err = dbs.Commit()
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return dberr(err)
	}

	return nil
}

func (s *Service) GetGame(ctx context.Context, id string, fpfss types.IFpfss) (*types.CachedGame, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}
	defer dbs.Rollback()

	game, err := s.pgdal.GetGame(dbs, id, fpfss)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	return game, nil
}
