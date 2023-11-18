package types

type FpfssGame struct {
	ID                  string      `json:"id"`
	Title               string      `json:"title"`
	Series              string      `json:"series"`
	Developer           string      `json:"developer"`
	Publisher           string      `json:"publisher"`
	ReleaseDate         string      `json:"release_date"`
	PlayMode            string      `json:"play_mode"`
	Language            string      `json:"language"`
	OriginalDescription string      `json:"original_description"`
	Platform            string      `json:"platform_name"`
	Tags                []*FpfssTag `json:"tags"`
}

type FpfssTag struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type ResponseFpfssGamesFetch struct {
	Games []*FpfssGame `json:"games"`
}

type IFpfss interface {
	GetGame(id string) (*FpfssGame, error)
	GetGames(ids []string) ([]*FpfssGame, error)
}
