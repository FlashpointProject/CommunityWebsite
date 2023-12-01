package database

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/FlashpointProject/CommunityWebsite/config"
	"github.com/FlashpointProject/CommunityWebsite/constants"
	"github.com/FlashpointProject/CommunityWebsite/types"
	"github.com/FlashpointProject/CommunityWebsite/utils"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type postgresDAL struct {
	db *pgxpool.Pool
}

func NewPostgresDAL(conn *pgxpool.Pool) *postgresDAL {
	return &postgresDAL{
		db: conn,
	}
}

// OpenPostgresDB opens DAL or panics
func OpenPostgresDB(l *logrus.Entry, conf *config.AppConfig) *pgxpool.Pool {
	l.Infoln("connecting to the postgres database")

	user := conf.PostgresUser
	pass := conf.PostgresPassword
	ip := conf.PostgresHost
	port := conf.PostgresPort

	conn, err := pgxpool.New(context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, pass, ip, port, user))

	if err != nil {
		l.Fatal(err)
	}

	l.Infoln("postgres database connected")
	return conn
}

type PostgresSession struct {
	context     context.Context
	transaction pgx.Tx
}

// NewSession begins a transaction
func (d *postgresDAL) NewSession(ctx context.Context) (PGDBSession, error) {
	tx, err := d.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	return &PostgresSession{
		context:     ctx,
		transaction: tx,
	}, nil
}

func (dbs *PostgresSession) Commit() error {
	return dbs.transaction.Commit(dbs.context)
}

func (dbs *PostgresSession) Rollback() error {
	err := dbs.Tx().Rollback(dbs.context)
	if err != nil && err.Error() == "sql: transaction has already been committed or rolled back" {
		err = nil
	}
	if err != nil {
		utils.LogCtx(dbs.Ctx()).Error(err)
	}
	return err
}

func (dbs *PostgresSession) Tx() pgx.Tx {
	return dbs.transaction
}

func (dbs *PostgresSession) Ctx() context.Context {
	return dbs.context
}

func (dal *postgresDAL) StoreSession(dbs PGDBSession, secret string, uid string, durationSeconds int64, ipAddr string) error {
	_, err := dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO session (secret, uid, expires_at, ip_addr) VALUES ($1, $2, NOW() + ($3 * INTERVAL '1 second'), $4)", secret, uid, durationSeconds, ipAddr)
	if err != nil {
		return err
	}
	return nil
}

func (dal *postgresDAL) GetRoles(dbs PGDBSession) ([]*types.DiscordRole, error) {
	var roles []*types.DiscordRole
	rows, err := dbs.Tx().Query(dbs.Ctx(), "SELECT id, name, color FROM fpcomm_role")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name string
		var color string
		err := rows.Scan(&id, &name, &color)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &types.DiscordRole{
			ID:    id,
			Name:  name,
			Color: color,
		})
	}
	return roles, nil
}

func (dal *postgresDAL) SaveRoles(dbs PGDBSession, roles []*types.DiscordRole) error {
	for _, role := range roles {
		_, err := dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO fpcomm_role (id, name, color) VALUES ($1, $2, $3) ON CONFLICT(id) DO UPDATE SET name = $2, color = $3", role.ID, role.Name, role.Color)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dal *postgresDAL) SaveUser(dbs PGDBSession, uid string, name string, avatarURL string, roles []string) error {
	_, err := dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO fpcomm_user (id, name, avatar, roles, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) ON CONFLICT(id) DO UPDATE SET name = $2, avatar = $3, roles = $4, updated_at = CURRENT_TIMESTAMP", uid, name, avatarURL, roles)
	if err != nil {
		return err
	}
	return nil
}

func (dal *postgresDAL) GetUser(dbs PGDBSession, uid string) (*types.UserProfile, error) {
	row := dbs.Tx().QueryRow(dbs.Ctx(), "SELECT name, avatar, roles, updated_at FROM fpcomm_user WHERE id=$1", uid)

	var name string
	var avatar string
	var roles []string
	var updatedAt time.Time
	err := row.Scan(&name, &avatar, &roles, &updatedAt)
	if err != nil {
		utils.LogCtx(dbs.Ctx()).Error(err)
		return nil, err
	}
	return &types.UserProfile{
		UserID:    uid,
		Username:  name,
		AvatarURL: avatar,
		Roles:     roles,
		UpdatedAt: updatedAt,
	}, nil
}

// GetSessionAuthInfo returns user ID + scope and/or expiration state
func (d *postgresDAL) GetSessionAuthInfo(dbs PGDBSession, secret string) (*types.SessionInfo, bool, error) {
	row := dbs.Tx().QueryRow(dbs.Ctx(), `SELECT id, uid, expires_at, ip_addr FROM session WHERE secret=$1`, secret)

	var id int64
	var uid string
	var expiration time.Time
	var ipAddr string
	err := row.Scan(&id, &uid, &expiration, &ipAddr)
	if err != nil {
		return nil, false, err
	}

	if expiration.Unix() <= time.Now().Unix() {
		return nil, false, nil
	}

	return &types.SessionInfo{
		ID:        id,
		UID:       uid,
		IpAddr:    ipAddr,
		ExpiresAt: expiration,
	}, true, nil
}

func (d *postgresDAL) SearchPlaylists(dbs PGDBSession, query *types.PlaylistSearchQuery) ([]*types.Playlist, int64, error) {
	total := 0

	playlists := make([]*types.Playlist, 0)

	builder := NewSqlBuilder("SELECT id, name, total_games, description, author_id, icon, library, public, extreme, filter_groups, created_at, updated_at FROM playlist")
	if query.UserID != "" {
		builder.Where("author_id=$1", query.UserID)
	}
	if query.Library != "" {
		builder.Where("library=$1", query.Library)
	}
	if query.Title != "" {
		builder.Where("name ILIKE $1", "%"+query.Title+"%")
	}
	if !query.Extreme {
		builder.Where("extreme=false")
	}
	builder.Limit(query.PageSize)
	builder.Offset((query.Page - 1) * query.PageSize)
	builder.OrderBy(query.OrderBy, query.OrderDirection, []string{"name", "created_at", "updated_at", "total_games"})

	sqlQuery := builder.Build(0)
	args := builder.Arguments()
	rows, err := dbs.Tx().Query(dbs.Ctx(), sqlQuery, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return playlists, 0, nil
		}
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		var totalGames int
		var description string
		var authorID string
		var icon string
		var library string
		var public bool
		var extreme bool
		var filterGroups []string
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(&id, &name, &totalGames, &description, &authorID, &icon, &library, &public, &extreme, &filterGroups, &createdAt, &updatedAt)
		if err != nil {
			return nil, 0, err
		}

		author := &types.UserProfile{
			UserID:    authorID,
			Username:  "Deleted User",
			AvatarURL: "",
			Roles:     []string{},
			UpdatedAt: time.Now(),
		} // Fill in after fetching all rows

		playlists = append(playlists, &types.Playlist{
			ID:           id,
			Name:         name,
			TotalGames:   totalGames,
			Description:  description,
			Author:       author,
			Library:      library,
			Icon:         icon,
			Public:       public,
			Extreme:      extreme,
			FilterGroups: filterGroups,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		})
	}

	for _, playlist := range playlists {
		err = d.db.QueryRow(dbs.Ctx(), "SELECT name, avatar, roles, updated_at FROM fpcomm_user WHERE id=$1", playlist.Author.UserID).Scan(&playlist.Author.Username, &playlist.Author.AvatarURL, &playlist.Author.Roles, &playlist.Author.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
	}

	if query.IncludeTotal {
		builder.SetBase("SELECT COUNT(*) FROM playlist")
		err = dbs.Tx().QueryRow(dbs.Ctx(), builder.Count(0), builder.ArgumentsCount()...).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}

	return playlists, int64(total), nil
}

func (d *postgresDAL) GetPlaylist(dbs PGDBSession, id int64) (*types.Playlist, error) {
	row := dbs.Tx().QueryRow(dbs.Ctx(), "SELECT id, name, total_games, description, author_id, icon, library, public, extreme, filter_groups, created_at, updated_at FROM playlist WHERE id=$1", id)

	var name string
	var totalGames int
	var description string
	var authorID string
	var icon string
	var library string
	var public bool
	var extreme bool
	var filterGroups []string
	var createdAt time.Time
	var updatedAt time.Time
	err := row.Scan(&id, &name, &totalGames, &description, &authorID, &icon, &library, &public, &extreme, &filterGroups, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	author, err := d.GetUser(dbs, authorID)
	if err != nil {
		if err == pgx.ErrNoRows {
			author = &types.UserProfile{
				UserID:    authorID,
				Username:  "Deleted User",
				AvatarURL: "",
				Roles:     []string{},
				UpdatedAt: time.Now(),
			}
		}
		return nil, err
	}

	games, err := d.GetPlaylistGames(dbs, id)
	if err != nil {
		return nil, err
	}

	playlist := &types.Playlist{
		ID:           id,
		Name:         name,
		TotalGames:   totalGames,
		Description:  description,
		Author:       author,
		Library:      library,
		Icon:         icon,
		Public:       public,
		Extreme:      extreme,
		FilterGroups: filterGroups,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Games:        games,
	}

	return playlist, nil
}

func (d *postgresDAL) GetPlaylistGames(dbs PGDBSession, id int64) ([]types.LauncherPlaylistGame, error) {
	games := make([]types.LauncherPlaylistGame, 0)

	rows, err := dbs.Tx().Query(dbs.Ctx(), "SELECT game_id, notes FROM playlist_game WHERE playlist_id=$1", id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return games, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var gameID string
		var notes string
		err := rows.Scan(&gameID, &notes)
		if err != nil {
			return nil, err
		}
		games = append(games, types.LauncherPlaylistGame{
			GameID: gameID,
			Notes:  notes,
		})
	}

	return games, nil
}

func (d *postgresDAL) GetGames(dbs PGDBSession, ids []string, fpfss types.IFpfss) ([]*types.CachedGame, error) {
	games := make([]*types.CachedGame, len(ids))
	for i, id := range ids {
		games[i] = &types.CachedGame{
			ID:      id,
			Missing: true,
		}
	}
	rows, err := dbs.Tx().Query(dbs.Ctx(), `SELECT id, title, series, developer, publisher, release_date, play_mode,
	 	language, original_description, platform_name, extreme, filter_groups, updated_at FROM game_cache WHERE id=ANY($1)`, ids)
	if err != nil {
		if err == pgx.ErrNoRows {
			return games, nil
		}
		utils.LogCtx(dbs.Ctx()).Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var title string
		var series string
		var developer string
		var publisher string
		var releaseDate string
		var playMode []string
		var language []string
		var originalDescription string
		var platformName string
		var extreme bool
		var filterGroups []string
		var updatedAt time.Time
		err := rows.Scan(&id, &title, &series, &developer, &publisher, &releaseDate, &playMode, &language, &originalDescription, &platformName, &extreme, &filterGroups, &updatedAt)
		if err != nil {
			utils.LogCtx(dbs.Ctx()).Error(err)
			return nil, err
		}
		for _, game := range games {
			if game.ID == id {
				game.Title = title
				game.Series = series
				game.Developer = developer
				game.Publisher = publisher
				game.ReleaseDate = releaseDate
				game.PlayMode = playMode
				game.Language = language
				game.OriginalDescription = originalDescription
				game.Platform = platformName
				game.UpdatedAt = updatedAt
				game.Extreme = extreme
				game.FilterGroups = filterGroups
				game.Missing = false
				break
			}
		}
	}

	// Check for outdated games
	yesterday := time.Now().Add(-time.Hour * 24)
	outdatedIds := make([]string, 0)
	for _, game := range games {
		if !game.Missing && game.UpdatedAt.Before(yesterday) {
			outdatedIds = append(outdatedIds, game.ID)
		}
	}
	if len(outdatedIds) != 0 {
		fpfssGames, err := fpfss.GetGames(outdatedIds)
		if err != nil {
			utils.LogCtx(dbs.Ctx()).Error(err)
			return nil, err
		}
		for _, fpfssGame := range fpfssGames {
			for i, game := range games {
				if game.ID == fpfssGame.ID {
					playModesArr := strings.Split(fpfssGame.PlayMode, ";")
					for i, pm := range playModesArr {
						playModesArr[i] = strings.TrimSpace(pm)
					}
					playModes := &pgtype.TextArray{}
					err := playModes.Set(playModesArr)
					if err != nil {
						utils.LogCtx(dbs.Ctx()).Error(err)
						return nil, err
					}
					languagesArr := strings.Split(fpfssGame.Language, ";")
					for i, lang := range languagesArr {
						languagesArr[i] = strings.TrimSpace(lang)
					}
					languages := &pgtype.TextArray{}
					err = languages.Set(languagesArr)
					if err != nil {
						utils.LogCtx(dbs.Ctx()).Error(err)
						return nil, err
					}

					// Build tag cache
					extreme := false
					var filterGroups = make([]string, 0)
					for _, tag := range fpfssGame.Tags {
						for _, filterGroup := range constants.GetFilterGroups() {
							if utils.StringInSlice(tag.Name, filterGroup.Tags) {
								filterGroups = append(filterGroups, filterGroup.Name)
								if filterGroup.Extreme {
									extreme = true
								}
								break
							}
						}
						_, err := dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO tag_cache (id, name, description, category) VALUES ($1, $2, $3, $4) ON CONFLICT(id) DO UPDATE SET name = $2, description = $3, category = $4, updated_at = CURRENT_TIMESTAMP", tag.ID, tag.Name, tag.Description, tag.Category)
						if err != nil {
							utils.LogCtx(dbs.Ctx()).Error(err)
							return nil, err
						}
					}
					filterGroups = utils.RemoveSliceDuplicates(filterGroups)
					sort.Strings(filterGroups)

					games[i] = &types.CachedGame{
						ID:                  fpfssGame.ID,
						Title:               fpfssGame.Title,
						Series:              fpfssGame.Series,
						Developer:           fpfssGame.Developer,
						Publisher:           fpfssGame.Publisher,
						ReleaseDate:         fpfssGame.ReleaseDate,
						PlayMode:            playModesArr,
						Language:            languagesArr,
						OriginalDescription: fpfssGame.OriginalDescription,
						Platform:            fpfssGame.Platform,
						Extreme:             extreme,
						FilterGroups:        filterGroups,
						UpdatedAt:           time.Now(),
					}
					// Save to cache
					_, err = dbs.Tx().Exec(dbs.Ctx(), "UPDATE game_cache SET title=$2, series=$3, developer=$4, publisher=$5, release_date=$6, play_mode=$7, language=$8, original_description=$9, platform_name=$10, updated_at=CURRENT_TIMESTAMP WHERE id=$1",
						games[i].ID, games[i].Title, games[i].Series, games[i].Developer, games[i].Publisher, games[i].ReleaseDate, playModes, languages, games[i].OriginalDescription, games[i].Platform)
					if err != nil {
						utils.LogCtx(dbs.Ctx()).Error(err)
						return nil, err
					}

					// Save tags as well
				}
			}
		}
	}

	idsMissing := make([]string, 0)
	for _, game := range games {
		if game.Missing {
			idsMissing = append(idsMissing, game.ID)
		}
	}
	if len(idsMissing) != 0 {
		fpfssGames, err := fpfss.GetGames(idsMissing)
		if err != nil {
			utils.LogCtx(dbs.Ctx()).Error(err)
			return nil, err
		}
		for _, fpfssGame := range fpfssGames {
			for i, game := range games {
				if game.ID == fpfssGame.ID {
					// Safely convert to postgres arrays
					playModesArr := strings.Split(fpfssGame.PlayMode, ";")
					for i, pm := range playModesArr {
						playModesArr[i] = strings.TrimSpace(pm)
					}
					playModes := &pgtype.TextArray{}
					err := playModes.Set(playModesArr)
					if err != nil {
						utils.LogCtx(dbs.Ctx()).Error(err)
						return nil, err
					}
					languagesArr := strings.Split(fpfssGame.Language, ";")
					for i, lang := range languagesArr {
						languagesArr[i] = strings.TrimSpace(lang)
					}
					languages := &pgtype.TextArray{}
					err = languages.Set(languagesArr)
					if err != nil {
						utils.LogCtx(dbs.Ctx()).Error(err)
						return nil, err
					}

					// Build tag cache
					extreme := false
					var filterGroups = make([]string, 0)
					for _, tag := range fpfssGame.Tags {
						for _, filterGroup := range constants.GetFilterGroups() {
							if utils.StringInSlice(tag.Name, filterGroup.Tags) {
								filterGroups = append(filterGroups, filterGroup.Name)
								if filterGroup.Extreme {
									extreme = true
								}
								break
							}
						}
						_, err := dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO tag_cache (id, name, description, category) VALUES ($1, $2, $3, $4) ON CONFLICT(id) DO UPDATE SET name = $2, description = $3, category = $4, updated_at = CURRENT_TIMESTAMP", tag.ID, tag.Name, tag.Description, tag.Category)
						if err != nil {
							utils.LogCtx(dbs.Ctx()).Error(err)
							return nil, err
						}
					}
					filterGroups = utils.RemoveSliceDuplicates(filterGroups)
					sort.Strings(filterGroups)

					games[i] = &types.CachedGame{
						ID:                  fpfssGame.ID,
						Title:               fpfssGame.Title,
						Series:              fpfssGame.Series,
						Developer:           fpfssGame.Developer,
						Publisher:           fpfssGame.Publisher,
						ReleaseDate:         fpfssGame.ReleaseDate,
						PlayMode:            playModesArr,
						Language:            languagesArr,
						OriginalDescription: fpfssGame.OriginalDescription,
						Platform:            fpfssGame.Platform,
						Extreme:             extreme,
						FilterGroups:        filterGroups,
						UpdatedAt:           time.Now(),
					}
					// Save to cache
					_, err = dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO game_cache (id, title, series, developer, publisher, release_date, play_mode, language, original_description, platform_name, extreme, filter_groups, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP)",
						games[i].ID, games[i].Title, games[i].Series, games[i].Developer, games[i].Publisher, games[i].ReleaseDate, playModes, languages, games[i].OriginalDescription, games[i].Platform, games[i].Extreme, games[i].FilterGroups)
					if err != nil {
						utils.LogCtx(dbs.Ctx()).Error(err)
						return nil, err
					}
				}
			}
		}
	}

	return games, nil
}

func (d *postgresDAL) GetGame(dbs PGDBSession, gameId string, fpfss types.IFpfss) (*types.CachedGame, error) {
	row := dbs.Tx().QueryRow(dbs.Ctx(), "SELECT id, title, series, developer, publisher, release_date, play_mode, language, original_description, platform_name, updated_at FROM game_cache WHERE id=$1", gameId)
	game, err := ReadGame(row)
	if err != nil {
		return nil, err
	}

	if game == nil {
		// No game, fetch from fpfss
		fpfssGame, err := fpfss.GetGame(gameId)
		if err != nil {
			return nil, err
		}
		if fpfssGame == nil {
			return nil, nil
		}
		playModesArr := strings.Split(fpfssGame.PlayMode, ";")
		for i, pm := range playModesArr {
			playModesArr[i] = strings.TrimSpace(pm)
		}
		playModes := &pgtype.TextArray{}
		err = playModes.Set(playModesArr)
		if err != nil {
			return nil, err
		}
		languagesArr := strings.Split(fpfssGame.Language, ";")
		for i, lang := range languagesArr {
			languagesArr[i] = strings.TrimSpace(lang)
		}
		languages := &pgtype.TextArray{}
		err = languages.Set(languagesArr)
		if err != nil {
			return nil, err
		}

		// Build tag cache
		extreme := false
		var filterGroups = make([]string, 0)
		for _, tag := range fpfssGame.Tags {
			for _, filterGroup := range constants.GetFilterGroups() {
				if utils.StringInSlice(tag.Name, filterGroup.Tags) {
					if filterGroup.Extreme {
						extreme = true
					}
					filterGroups = append(filterGroups, filterGroup.Name)
					break
				}
			}
			_, err := dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO tag_cache (id, name, description, category) VALUES ($1, $2, $3, $4) ON CONFLICT(id) DO UPDATE SET name = $2, description = $3, category = $4, updated_at = CURRENT_TIMESTAMP", tag.ID, tag.Name, tag.Description, tag.Category)
			if err != nil {
				utils.LogCtx(dbs.Ctx()).Error(err)
				return nil, err
			}
		}
		filterGroups = utils.RemoveSliceDuplicates(filterGroups)
		sort.Strings(filterGroups)

		game := &types.CachedGame{
			ID:                  fpfssGame.ID,
			Title:               fpfssGame.Title,
			Series:              fpfssGame.Series,
			Developer:           fpfssGame.Developer,
			Publisher:           fpfssGame.Publisher,
			ReleaseDate:         fpfssGame.ReleaseDate,
			PlayMode:            playModesArr,
			Language:            languagesArr,
			OriginalDescription: fpfssGame.OriginalDescription,
			Platform:            fpfssGame.Platform,
			Extreme:             extreme,
			FilterGroups:        filterGroups,
			UpdatedAt:           time.Now(),
		}
		_, err = dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO game_cache (id, title, series, developer, publisher, release_date, play_mode, language, original_description, platform_name, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP)",
			game.ID, game.Title, game.Series, game.Developer, game.Publisher, game.ReleaseDate, playModes, languages, game.OriginalDescription, game.Platform)
		if err != nil {
			return nil, err
		}
		return game, nil
	} else {
		// Found game, check if it needs updated
		// Update cache daily
		if game.UpdatedAt.Before(time.Now().Add(-time.Hour * 24)) {
			fpfssGame, err := fpfss.GetGame(gameId)
			if err != nil {
				return nil, err
			}
			if fpfssGame == nil {
				// Game deleted from FPFSS
				_, err = dbs.Tx().Exec(dbs.Ctx(), "DELETE FROM game_cache WHERE id=$1", gameId)
				if err != nil {
					return nil, err
				}
				return nil, nil
			}
			// Update from fpfss
			playModesArr := strings.Split(fpfssGame.PlayMode, ";")
			for i, pm := range playModesArr {
				playModesArr[i] = strings.TrimSpace(pm)
			}
			playModes := &pgtype.TextArray{}
			err = playModes.Set(playModesArr)
			if err != nil {
				return nil, err
			}
			languagesArr := strings.Split(fpfssGame.Language, ";")
			for i, lang := range languagesArr {
				languagesArr[i] = strings.TrimSpace(lang)
			}
			languages := &pgtype.TextArray{}
			err = languages.Set(languagesArr)
			if err != nil {
				return nil, err
			}

			// Build tag cache
			extreme := false
			var filterGroups = make([]string, 0)
			for _, tag := range fpfssGame.Tags {
				for _, filterGroup := range constants.GetFilterGroups() {
					if utils.StringInSlice(tag.Name, filterGroup.Tags) {
						if filterGroup.Extreme {
							extreme = true
						}
						filterGroups = append(filterGroups, filterGroup.Name)
						break
					}
				}
				_, err := dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO tag_cache (id, name, description, category) VALUES ($1, $2, $3, $4) ON CONFLICT(id) DO UPDATE SET name = $2, description = $3, category = $4, updated_at = CURRENT_TIMESTAMP", tag.ID, tag.Name, tag.Description, tag.Category)
				if err != nil {
					utils.LogCtx(dbs.Ctx()).Error(err)
					return nil, err
				}
			}
			filterGroups = utils.RemoveSliceDuplicates(filterGroups)
			sort.Strings(filterGroups)

			game := &types.CachedGame{
				ID:                  fpfssGame.ID,
				Title:               fpfssGame.Title,
				Series:              fpfssGame.Series,
				Developer:           fpfssGame.Developer,
				Publisher:           fpfssGame.Publisher,
				ReleaseDate:         fpfssGame.ReleaseDate,
				PlayMode:            playModesArr,
				Language:            languagesArr,
				OriginalDescription: fpfssGame.OriginalDescription,
				Platform:            fpfssGame.Platform,
				Extreme:             extreme,
				FilterGroups:        filterGroups,
				UpdatedAt:           time.Now(),
			}
			_, err = dbs.Tx().Exec(dbs.Ctx(), "UPDATE game_cache SET title=$2, series=$3, developer=$4, publisher=$5, release_date=$6, play_mode=$7, language=$8, original_description=$9, platform_name=$10, updated_at=CURRENT_TIMESTAMP WHERE id=$1",
				game.ID, game.Title, game.Series, game.Developer, game.Publisher, game.ReleaseDate, playModes, languages, game.OriginalDescription, game.Platform)
			if err != nil {
				return nil, err
			}
			return game, nil
		} else {
			return game, nil
		}
	}
}

func (d *postgresDAL) SavePlaylist(dbs PGDBSession, uid string, playlist *types.Playlist, fpfss types.IFpfss) error {
	// force cache games
	var filterGroups = make([]string, 0)
	extreme := false
	gameIds := make([]string, 0)
	for _, game := range playlist.Games {
		gameIds = append(gameIds, game.GameID)
	}
	games, err := d.GetGames(dbs, gameIds, fpfss)
	if err != nil {
		return err
	}

	for _, game := range games {
		if game.Extreme {
			extreme = true
		}
		filterGroups = append(filterGroups, game.FilterGroups...)
	}
	filterGroups = utils.RemoveSliceDuplicates(filterGroups)
	sort.Strings(filterGroups)

	if playlist.ID == 0 {
		err = dbs.Tx().QueryRow(dbs.Ctx(), "INSERT INTO playlist (name, total_games, description, author_id, icon, public, extreme, filter_groups, library) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
			playlist.Name, playlist.TotalGames, playlist.Description, uid, playlist.Icon, playlist.Public, extreme, filterGroups, playlist.Library).Scan(&playlist.ID)
		if err != nil {
			return err
		}
	} else {
		_, err = dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO playlist (id, name, total_games, description, author_id, icon, public, extreme, filter_groups, library) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (id) DO UPDATE SET name = $2, total_games = $3, description = $4, author_id = $5, icon = $6, public = $7, extreme = $8, filter_groups = $9, library = $10",
			playlist.ID, playlist.Name, playlist.TotalGames, playlist.Description, uid, playlist.Icon, playlist.Public, extreme, filterGroups, playlist.Library)
		if err != nil {
			return err
		}
	}

	_, err = dbs.Tx().Exec(dbs.Ctx(), "DELETE FROM playlist_game WHERE playlist_id=$1", playlist.ID)
	if err != nil {
		return err
	}

	for _, game := range playlist.Games {
		_, err = dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO playlist_game (playlist_id, game_id, notes) VALUES ($1, $2, $3)", playlist.ID, game.GameID, game.Notes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *postgresDAL) DeletePlaylist(dbs PGDBSession, id int64) error {
	_, err := dbs.Tx().Exec(dbs.Ctx(), "DELETE FROM playlist_game WHERE playlist_id=$1", id)
	if err != nil {
		return err
	}
	_, err = dbs.Tx().Exec(dbs.Ctx(), "DELETE FROM playlist WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (d *postgresDAL) SearchNewsPosts(dbs PGDBSession, query *types.NewsPostSearchQuery) ([]*types.NewsPost, int64, error) {
	total := int64(0)

	posts := make([]*types.NewsPost, 0)

	builder := NewSqlBuilder("SELECT id, title, content, post_type, author_id, created_at, updated_at FROM post")
	if query.AuthorID != "" {
		builder.Where("author_id=$1", query.AuthorID)
	}
	if query.PostType != "" {
		builder.Where("post_type=$1", query.PostType)
	}
	if query.Title != "" {
		builder.Where("title ILIKE $1", "%"+query.Title+"%")
	}
	builder.Limit(query.PageSize)
	builder.Offset((query.Page - 1) * query.PageSize)
	builder.OrderBy(query.OrderBy, query.OrderDirection, []string{"title", "created_at", "updated_at"})

	sqlQuery := builder.Build(0)
	args := builder.Arguments()
	rows, err := dbs.Tx().Query(dbs.Ctx(), sqlQuery, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return posts, 0, nil
		}
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var title string
		var content string
		var postType string
		var authorID string
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(&id, &title, &content, &postType, &authorID, &createdAt, &updatedAt)
		if err != nil {
			return nil, 0, err
		}

		author := &types.UserProfile{
			UserID:    authorID,
			Username:  "Deleted User",
			AvatarURL: "",
			Roles:     []string{},
			UpdatedAt: time.Now(),
		} // Fill in after fetching all rows

		posts = append(posts, &types.NewsPost{
			ID:        id,
			Title:     title,
			Content:   content,
			PostType:  postType,
			Author:    author,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	for _, post := range posts {
		err = d.db.QueryRow(dbs.Ctx(), "SELECT name, avatar, roles, updated_at FROM fpcomm_user WHERE id=$1", post.Author.UserID).Scan(&post.Author.Username, &post.Author.AvatarURL, &post.Author.Roles, &post.Author.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
	}

	if query.IncludeTotal {
		builder.SetBase("SELECT COUNT(*) FROM post")
		err = dbs.Tx().QueryRow(dbs.Ctx(), builder.Count(0), builder.ArgumentsCount()...).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}

	return posts, total, nil
}

func (d *postgresDAL) GetNewsPost(dbs PGDBSession, id int64) (*types.NewsPost, error) {
	row := dbs.Tx().QueryRow(dbs.Ctx(), "SELECT id, title, content, post_type, author_id, created_at, updated_at FROM post WHERE id=$1", id)

	var title string
	var content string
	var postType string
	var authorID string
	var createdAt time.Time
	var updatedAt time.Time
	err := row.Scan(&id, &title, &content, &postType, &authorID, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	author, err := d.GetUser(dbs, authorID)
	if err != nil {
		if err == pgx.ErrNoRows {
			author = &types.UserProfile{
				UserID:    authorID,
				Username:  "Deleted User",
				AvatarURL: "",
				Roles:     []string{},
				UpdatedAt: time.Now(),
			}
		}
		return nil, err
	}

	post := &types.NewsPost{
		ID:        id,
		Title:     title,
		Content:   content,
		PostType:  postType,
		Author:    author,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return post, nil
}

func (d *postgresDAL) SaveNewsPost(dbs PGDBSession, uid string, post *types.NewsPost) error {
	var id int64
	err := dbs.Tx().QueryRow(dbs.Ctx(), "INSERT INTO post (title, content, post_type, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id",
		post.Title, post.Content, post.PostType, uid).Scan(&id)
	if err != nil {
		return err
	}
	post.ID = id
	return nil
}

func (d *postgresDAL) SearchContentReports(dbs PGDBSession, query *types.ContentReportSearchQuery) ([]*types.ContentReport, int64, error) {
	total := int64(0)

	reports := make([]*types.ContentReport, 0)

	builder := NewSqlBuilder(`SELECT id, content_ref, report_state, reported_by,
		report_reason, context, reported_user, resolved_by, resolved_at, action_taken, created_at, updated_at FROM content_report`)

	if query.ContentType != "" {
		builder.Where("content_ref ILIKE $1", "%"+query.ContentType+"%")
	}
	if query.ReportState != "" {
		builder.Where("report_state=$1", query.ReportState)
	}
	if query.ReportedBy != "" {
		builder.Where("reported_by=$1", query.ReportedBy)
	}
	if query.ReportedUser != "" {
		builder.Where("reported_user=$1", query.ReportedUser)
	}
	if query.ResolvedBy != "" {
		builder.Where("resolved_by=$1", query.ResolvedBy)
	}
	builder.Limit(query.PageSize)
	builder.Offset((query.Page - 1) * query.PageSize)
	builder.OrderBy(query.OrderBy, query.OrderDirection, []string{"created_at", "updated_at", "resolved_at"})

	sqlQuery := builder.Build(0)
	args := builder.Arguments()
	rows, err := dbs.Tx().Query(dbs.Ctx(), sqlQuery, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return reports, 0, nil
		}
		return nil, 0, err
	}

	for rows.Next() {
		var id int64
		var contentRef string
		var reportState string
		var reportedBy string
		var reportReason string
		var context string
		var reportedUser string
		var resolvedBy string
		var resolvedAt sql.NullTime
		var actionTaken string
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(&id, &contentRef, &reportState, &reportedBy, &reportReason, &context, &reportedUser, &resolvedBy, &resolvedAt, &actionTaken, &createdAt, &updatedAt)
		if err != nil {
			return nil, 0, err
		}
		var resolvedAtTime *time.Time
		if resolvedAt.Valid {
			resolvedAtTime = &resolvedAt.Time
		} else {
			resolvedAtTime = nil
		}
		reports = append(reports, &types.ContentReport{
			ID:          id,
			ContentRef:  contentRef,
			ReportState: reportState,
			ReportedBy: &types.UserProfile{
				UserID: reportedBy,
			},
			ReportReason: reportReason,
			ReportedUser: &types.UserProfile{
				UserID: reportedUser,
			},
			ResolvedBy: &types.UserProfile{
				UserID: resolvedBy,
			},
			ResolvedAt:  resolvedAtTime,
			ActionTaken: actionTaken,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	for _, report := range reports {
		if report.ReportedBy.UserID != "" {
			reportedByUser, err := d.GetUser(dbs, report.ReportedBy.UserID)
			if err != nil {
				if err == pgx.ErrNoRows {
					report.ReportedBy = &types.UserProfile{
						UserID:    report.ReportedBy.UserID,
						Username:  "Deleted User",
						AvatarURL: "",
						Roles:     []string{},
						UpdatedAt: time.Now(),
					}
				} else {
					return nil, 0, err
				}
			} else {
				report.ReportedBy = reportedByUser
			}
		} else {
			report.ReportedBy = &types.UserProfile{
				UserID:    report.ReportedUser.UserID,
				Username:  "",
				AvatarURL: "",
				Roles:     []string{},
				UpdatedAt: time.Now(),
			}
		}

		if report.ResolvedBy.UserID != "" {
			resolvedByUser, err := d.GetUser(dbs, report.ResolvedBy.UserID)
			if err != nil {
				if err == pgx.ErrNoRows {
					report.ReportedBy = &types.UserProfile{
						UserID:    report.ReportedBy.UserID,
						Username:  "Deleted User",
						AvatarURL: "",
						Roles:     []string{},
						UpdatedAt: time.Now(),
					}
				} else {
					return nil, 0, err
				}
			} else {
				report.ResolvedBy = resolvedByUser
			}
		} else {
			report.ResolvedBy = &types.UserProfile{
				UserID:    report.ResolvedBy.UserID,
				Username:  "",
				AvatarURL: "",
				Roles:     []string{},
				UpdatedAt: time.Now(),
			}
		}

		if report.ReportedUser.UserID != "" {
			reportedUserUser, err := d.GetUser(dbs, report.ReportedUser.UserID)
			if err != nil {
				if err == pgx.ErrNoRows {
					report.ReportedUser = &types.UserProfile{
						UserID:    report.ReportedUser.UserID,
						Username:  "Deleted User",
						AvatarURL: "",
						Roles:     []string{},
						UpdatedAt: time.Now(),
					}
				} else {
					return nil, 0, err
				}
			} else {
				report.ReportedUser = reportedUserUser
			}
		} else {
			report.ReportedUser = &types.UserProfile{
				UserID:    report.ReportedBy.UserID,
				Username:  "",
				AvatarURL: "",
				Roles:     []string{},
				UpdatedAt: time.Now(),
			}
		}
	}

	if query.IncludeTotal {
		builder.SetBase("SELECT COUNT(*) FROM content_report")
		err = dbs.Tx().QueryRow(dbs.Ctx(), builder.Count(0), builder.ArgumentsCount()...).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}

	return reports, total, nil
}

func (d *postgresDAL) SaveContentReport(dbs PGDBSession, report *types.ContentReport) error {
	var id int64
	err := dbs.Tx().QueryRow(dbs.Ctx(), "INSERT INTO content_report (content_ref, report_state, reported_by, report_reason, context, resolved_by, action_taken, reported_user, resolved_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id",
		report.ContentRef, report.ReportState, report.ReportedBy.UserID, report.ReportReason, report.AdditionalContext, report.ResolvedBy.UserID, report.ActionTaken, report.ReportedUser.UserID).Scan(&id)
	if err != nil {
		return err
	}
	report.ID = id
	return nil
}

func (d *postgresDAL) SearchGotdSuggestions(dbs PGDBSession, query *types.GotdSuggestionsSearchQuery, fpfss types.IFpfss) ([]*types.GotdSuggestionInternal, int64, error) {
	results := make([]*types.GotdSuggestionInternal, 0)
	var total int64

	builder := NewSqlBuilder("SELECT id, game_id, author_id, anonymous, description, suggested_date, created_at FROM gotd_suggestion")
	builder.Limit(query.PageSize)
	builder.Offset((query.Page - 1) * query.PageSize)
	builder.OrderBy(query.OrderBy, query.OrderDirection, []string{"created_at", "suggested_date"})

	sqlQuery := builder.Build(0)
	args := builder.Arguments()
	rows, err := dbs.Tx().Query(dbs.Ctx(), sqlQuery, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return results, 0, nil
		}
		return nil, 0, err
	}

	for rows.Next() {
		var id int64
		var gameID string
		var authorID string
		var anonymous bool
		var description string
		var suggestedDateInternal sql.NullTime
		var createdAt time.Time
		err := rows.Scan(&id, &gameID, &authorID, &anonymous, &description, &suggestedDateInternal, &createdAt)
		if err != nil {
			return nil, 0, err
		}

		game, err := d.GetGame(dbs, gameID, fpfss)
		if err != nil {
			return nil, 0, err
		}

		author, err := d.GetUser(dbs, authorID)
		if err != nil {
			if err == pgx.ErrNoRows {
				author = &types.UserProfile{
					UserID:    authorID,
					Username:  "Deleted User",
					AvatarURL: "",
					Roles:     []string{},
					UpdatedAt: time.Now(),
				}
			}
			return nil, 0, err
		}
		var suggestedDate *time.Time
		if suggestedDateInternal.Valid {
			suggestedDate = &suggestedDateInternal.Time
		} else {
			suggestedDate = nil
		}

		results = append(results, &types.GotdSuggestionInternal{
			ID:            id,
			Game:          game,
			Author:        author,
			Anonymous:     anonymous,
			Description:   description,
			SuggestedDate: suggestedDate,
			CreatedAt:     createdAt,
		})
	}

	if query.IncludeTotal {
		builder.SetBase("SELECT COUNT(*) FROM gotd_suggestion")
		err = dbs.Tx().QueryRow(dbs.Ctx(), builder.Count(0), builder.ArgumentsCount()...).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}

	return results, total, nil
}

func (d *postgresDAL) GetGotdSuggestion(dbs PGDBSession, sugId int64, fpfss types.IFpfss) (*types.GotdSuggestionInternal, error) {
	row := dbs.Tx().QueryRow(dbs.Ctx(), "SELECT id, game_id, author_id, anonymous, description, suggested_date FROM gotd_suggestion WHERE id=$1", sugId)
	var suggestion *types.GotdSuggestionInternal
	var authorID string
	var gameID string
	err := row.Scan(&suggestion.ID, gameID, &authorID, &suggestion.Anonymous, &suggestion.Description, &suggestion.SuggestedDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	game, err := d.GetGame(dbs, gameID, fpfss)
	if err != nil {
		return nil, err
	}
	suggestion.Game = game

	author, err := d.GetUser(dbs, authorID)
	if err != nil {
		if err == pgx.ErrNoRows {
			author = &types.UserProfile{
				UserID:    authorID,
				Username:  "Deleted User",
				AvatarURL: "",
				Roles:     []string{},
				UpdatedAt: time.Now(),
			}
		}
		return nil, err
	}
	suggestion.Author = author

	return suggestion, nil
}

func (d *postgresDAL) DeleteGotdSuggestion(dbs PGDBSession, uid string, sugId int64) error {
	_, err := dbs.Tx().Exec(dbs.Ctx(), "DELETE FROM gotd_suggestion WHERE id=$1", sugId)
	if err != nil {
		return err
	}

	return nil
}

func (d *postgresDAL) SaveGotdSuggestion(dbs PGDBSession, uid string, suggestion *types.GotdSuggestionInternal) error {
	row := dbs.Tx().QueryRow(dbs.Ctx(), "INSERT INTO gotd_suggestion (game_id, author_id, anonymous, description, suggested_date) VALUES ($1, $2, $3, $4) RETURNING id",
		suggestion.Game.ID, uid, suggestion.Anonymous, suggestion.Description, suggestion.SuggestedDate)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	suggestion.ID = id

	return nil
}

func (d *postgresDAL) GetGotdCurrent(dbs PGDBSession, query *types.GetGotdCurrentQuery) ([]*types.GotdGame, error) {
	base := "SELECT game_id, author, description, assigned_date FROM gotd"
	if !query.ShowFuture {
		base += " WHERE assigned_date < CURRENT_TIMESTAMP"
	}
	base += " ORDER BY assigned_date ASC"
	games := make([]*types.GotdGame, 0)
	rows, err := dbs.Tx().Query(dbs.Ctx(), base)
	if err != nil {
		if err == pgx.ErrNoRows {
			return games, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var gameID string
		var author string
		var description string
		var assignedDate time.Time
		err := rows.Scan(&gameID, &author, &description, &assignedDate)
		if err != nil {
			return nil, err
		}
		games = append(games, &types.GotdGame{
			ID:           gameID,
			Author:       author,
			Description:  description,
			AssignedDate: assignedDate,
		})
	}

	return games, nil
}

func (d *postgresDAL) AssignGotd(dbs PGDBSession, uid string, sugId int64, date string, fpfss types.IFpfss) error {
	suggestion, err := d.GetGotdSuggestion(dbs, sugId, fpfss)
	if err != nil {
		return err
	}
	if suggestion == nil {
		return fmt.Errorf("suggestion not found")
	}

	authorName := "Anonymous"
	if !suggestion.Anonymous && suggestion.Author.Username == "Deleted User" {
		authorName = suggestion.Author.Username
	}

	_, err = dbs.Tx().Exec(dbs.Ctx(), "INSERT INTO gotd (game_id, author, description, assigned_date) VALUES ($1, $2, $3, $4)",
		suggestion.Game.ID, authorName, suggestion.Description, date)
	if err != nil {
		return err
	}

	return nil
}

func (d *postgresDAL) UnassignGotd(dbs PGDBSession, uid string, date string) error {
	_, err := dbs.Tx().Exec(dbs.Ctx(), "DELETE FROM gotd WHERE assigned_date=$1", date)
	if err != nil {
		return err
	}
	return nil
}
