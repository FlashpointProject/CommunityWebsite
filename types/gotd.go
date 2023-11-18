package types

import "time"

type GotdGame struct {
	ID           string    `json:"id"`
	Author       string    `json:"author"`
	Description  string    `json:"description"`
	AssignedDate time.Time `json:"date"`
}

type GotdFileGame struct {
	ID          string `json:"id"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type GotdFile struct {
	Games []GotdFileGame `json:"games"`
}

type GotdSuggestion struct {
	ID            int64       `json:"id"`
	Game          *CachedGame `json:"game"`
	Author        string      `json:"author"`
	Description   string      `json:"description"`
	SuggestedDate *time.Time  `json:"suggested_date"`
	CreatedAt     time.Time   `json:"created_at"`
}

type GotdSuggestionInternal struct {
	ID            int64        `json:"id"`
	Game          *CachedGame  `json:"game"`
	Author        *UserProfile `json:"author"`
	Anonymous     bool         `json:"anonymous"`
	Description   string       `json:"description"`
	SuggestedDate *time.Time   `json:"suggested_date"`
	CreatedAt     time.Time    `json:"created_at"`
}

type GetGotdCurrentQuery struct {
	ShowFuture bool `json:"show_future"`
}

type GotdSuggestionsSearchQuery struct {
	Page           int64  `schema:"page"`
	PageSize       int64  `schema:"page_size"`
	OrderBy        string `schema:"order_by"`
	OrderDirection string `schema:"order_direction"`
	IncludeTotal   bool   `schema:"include_total"`
}

type GotdSuggestionsSearchResponse struct {
	Suggestions []*GotdSuggestion `json:"suggestions"`
	Total       int64             `json:"total"`
}

func (s *GotdSuggestionInternal) ToExternal() *GotdSuggestion {
	var author string
	if s.Anonymous {
		author = "Anonymous"
	} else {
		author = s.Author.Username
	}
	return &GotdSuggestion{
		ID:            s.ID,
		Game:          s.Game,
		Author:        author,
		Description:   s.Description,
		SuggestedDate: s.SuggestedDate,
		CreatedAt:     s.CreatedAt,
	}
}
