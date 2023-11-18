package types

import "time"

type CachedGame struct {
	ID                  string       `json:"id"`
	Title               string       `json:"title"`
	Series              string       `json:"series"`
	Developer           string       `json:"developer"`
	Publisher           string       `json:"publisher"`
	ReleaseDate         string       `json:"release_date"`
	PlayMode            []string     `json:"play_mode"`
	Language            []string     `json:"language"`
	OriginalDescription string       `json:"original_description"`
	Platform            string       `json:"platform"`
	Extreme             bool         `json:"extreme"`
	FilterGroups        []string     `json:"filter_groups"`
	Tags                []*CachedTag `json:"tags"`
	UpdatedAt           time.Time    `json:"updated_at"`
	Missing             bool         `json:"missing"`
}

type CachedTag struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}
