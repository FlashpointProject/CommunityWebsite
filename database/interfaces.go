package database

import (
	"context"

	"github.com/FlashpointProject/CommunityWebsite/types"
	"github.com/jackc/pgx/v5"
)

type PGDAL interface {
	NewSession(ctx context.Context) (PGDBSession, error)
	StoreSession(dbs PGDBSession, secret string, uid string, durationSeconds int64, ipAddr string) error
	GetSessionAuthInfo(dbs PGDBSession, secret string) (*types.SessionInfo, bool, error)

	GetRoles(dbs PGDBSession) ([]*types.DiscordRole, error)
	SaveRoles(dbs PGDBSession, roles []*types.DiscordRole) error
	SaveUser(dbs PGDBSession, uid string, name string, avatarURL string, roles []string) error
	GetUser(dbs PGDBSession, uid string) (*types.UserProfile, error)

	SearchPlaylists(dbs PGDBSession, query *types.PlaylistSearchQuery) ([]*types.Playlist, int64, error)
	GetPlaylist(dbs PGDBSession, id int64) (*types.Playlist, error)
	SavePlaylist(dbs PGDBSession, uid string, playlist *types.Playlist, fpfss types.IFpfss) error
	DeletePlaylist(dbs PGDBSession, id int64) error

	GetGames(dbs PGDBSession, ids []string, fpfss types.IFpfss) ([]*types.CachedGame, error)
	GetGame(dbs PGDBSession, id string, fpfss types.IFpfss) (*types.CachedGame, error)

	SearchNewsPosts(dbs PGDBSession, query *types.NewsPostSearchQuery) ([]*types.NewsPost, int64, error)
	GetNewsPost(dbs PGDBSession, id int64) (*types.NewsPost, error)
	SaveNewsPost(dbs PGDBSession, uid string, post *types.NewsPost) error

	SearchContentReports(dbs PGDBSession, query *types.ContentReportSearchQuery) ([]*types.ContentReport, int64, error)
	SaveContentReport(dbs PGDBSession, report *types.ContentReport) error

	SearchGotdSuggestions(dbs PGDBSession, query *types.GotdSuggestionsSearchQuery, fpfss types.IFpfss) ([]*types.GotdSuggestionInternal, int64, error)
	GetGotdSuggestion(dbs PGDBSession, sugId int64, fpfss types.IFpfss) (*types.GotdSuggestionInternal, error)
	DeleteGotdSuggestion(dbs PGDBSession, uid string, sugId int64) error
	SaveGotdSuggestion(dbs PGDBSession, uid string, suggestion *types.GotdSuggestionInternal) error
	GetGotdCurrent(dbs PGDBSession, query *types.GetGotdCurrentQuery) ([]*types.GotdGame, error)
	AssignGotd(dbs PGDBSession, uid string, sugId int64, date string, fpfss types.IFpfss) error
	UnassignGotd(dbs PGDBSession, uid string, date string) error
}

type PGDBSession interface {
	Commit() error
	Rollback() error
	Tx() pgx.Tx
	Ctx() context.Context
}
