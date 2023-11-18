package types

import "time"

type SubmittedNewsPost struct {
	PostType string `json:"post_type"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type NewsPost struct {
	ID        int64        `json:"id"`
	PostType  string       `json:"post_type"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Author    *UserProfile `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
