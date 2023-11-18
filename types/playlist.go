package types

import "time"

type PlaylistInfo struct {
	ID           int64        `json:"id"`
	Name         string       `json:"name"`
	TotalGames   int          `json:"total_games"`
	Description  string       `json:"description"`
	Author       *UserProfile `json:"author"`
	Library      string       `json:"library"`
	Icon         string       `json:"icon"`
	Public       bool         `json:"public"`
	Extreme      bool         `json:"extreme"`
	FilterGroups []string     `json:"filter_groups"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type Playlist struct {
	ID           int64                  `json:"id"`
	Name         string                 `json:"name"`
	TotalGames   int                    `json:"total_games"`
	Description  string                 `json:"description"`
	Author       *UserProfile           `json:"author"`
	Library      string                 `json:"library"`
	Icon         string                 `json:"icon"`
	Games        []LauncherPlaylistGame `json:"games"`
	Public       bool                   `json:"public"`
	Extreme      bool                   `json:"extreme"`
	FilterGroups []string               `json:"filter_groups"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type FullPlaylist struct {
	ID           int64           `json:"id"`
	Name         string          `json:"name"`
	TotalGames   int             `json:"total_games"`
	Description  string          `json:"description"`
	Author       *UserProfile    `json:"author"`
	Library      string          `json:"library"`
	Icon         string          `json:"icon"`
	Games        []GameWithNotes `json:"games"`
	Public       bool            `json:"public"`
	Extreme      bool            `json:"extreme"`
	FilterGroups []string        `json:"filter_groups"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type GameWithNotes struct {
	Game   *CachedGame `json:"game"`
	GameID string      `json:"game_id"`
	Notes  string      `json:"notes"`
}

type LauncherPlaylistGame struct {
	GameID string `json:"gameId"`
	Notes  string `json:"notes"`
}

type LauncherPlaylist struct {
	ID          string                 `json:"id"`
	Author      string                 `json:"author"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Library     string                 `json:"library"`
	Icon        string                 `json:"icon"`
	Extreme     bool                   `json:"extreme"`
	Games       []LauncherPlaylistGame `json:"games"`
}

type PlaylistSearchQuery struct {
	Page           int64  `json:"page" schema:"page"`
	PageSize       int64  `json:"page_size" schema:"page_size"`
	UserID         string `json:"user_id" schema:"user_id"`
	Library        string `json:"library" schema:"library"`
	Title          string `json:"title" schema:"title"`
	Extreme        bool   `json:"extreme" schema:"extreme"`
	OrderBy        string `json:"order_by" schema:"order_by"`
	OrderDirection string `json:"order_direction" schema:"order_direction"`
	IncludeTotal   bool   `json:"include_total" schema:"include_total"`
}

type PlaylistSearchResponse struct {
	Playlists []*PlaylistInfo `json:"playlists"`
	Total     int64           `json:"total"`
}

type NewsPostSearchQuery struct {
	Page           int64  `json:"page" schema:"page"`
	PageSize       int64  `json:"page_size" schema:"page_size"`
	OrderBy        string `json:"order_by" schema:"order_by"`
	OrderDirection string `json:"order_direction" schema:"order_direction"`
	IncludeTotal   bool   `json:"include_total" schema:"include_total"`
	PostType       string `json:"post_type" schema:"post_type"`
	AuthorID       string `json:"author_id" schema:"author_id"`
	Title          string `json:"title" schema:"title"`
}

type NewsPostSearchResponse struct {
	Posts []*NewsPost `json:"posts"`
	Total int64       `json:"total"`
}
