package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FlashpointProject/CommunityWebsite/database"
	"github.com/FlashpointProject/CommunityWebsite/types"
	"github.com/FlashpointProject/CommunityWebsite/utils"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pgdal                    database.PGDAL
	authTokenProvider        AuthTokenizer
	sessionExpirationSeconds int64
	RoleCache                []*types.DiscordRole
}

type authTokenProvider struct {
}

func NewAuthTokenProvider() *authTokenProvider {
	return &authTokenProvider{}
}

type AuthTokenizer interface {
	CreateAuthToken(userID string) (*authToken, error)
}

func (a *authTokenProvider) CreateAuthToken(userID string) (*authToken, error) {
	s, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &authToken{
		Secret: s.String(),
		UserID: userID,
	}, nil
}

func NewService(pgdb *pgxpool.Pool, sessionExpirationSeconds int64) *Service {
	return &Service{
		pgdal:                    database.NewPostgresDAL(pgdb),
		authTokenProvider:        NewAuthTokenProvider(),
		sessionExpirationSeconds: sessionExpirationSeconds,
	}
}

func (s *Service) GetSessionAuthInfo(ctx context.Context, key string) (*types.SessionInfo, bool, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, false, dberr(err)
	}
	defer dbs.Rollback()

	info, ok, err := s.pgdal.GetSessionAuthInfo(dbs, key)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, false, dberr(err)
	}

	return info, ok, nil
}

func (s *Service) GetUser(ctx context.Context, uid string) (*types.UserProfile, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}
	defer dbs.Rollback()

	user, err := s.pgdal.GetUser(dbs, uid)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, nil
		}
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	return user, nil
}

func (s *Service) LoadRoles(ctx context.Context) error {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		return dberr(err)
	}
	defer dbs.Rollback()

	roles, err := s.pgdal.GetRoles(dbs)
	if err != nil {
		if err != sql.ErrNoRows {
			return dberr(err)
		}
		return nil
	}

	s.RoleCache = roles

	return nil
}

func (s *Service) CacheRoles(roles []*types.DiscordRole) {
	for _, role := range roles {
		found := false
		for _, cachedRole := range s.RoleCache {
			if role.ID == cachedRole.ID {
				cachedRole.Name = role.Name
				cachedRole.Color = role.Color
				found = true
				break
			}
		}
		if !found {
			s.RoleCache = append(s.RoleCache, role)
		}
	}
}

// ParseAuthToken parses map into token
func ParseAuthToken(value map[string]string) (*authToken, error) {
	secret, ok := value["Secret"]
	if !ok {
		return nil, fmt.Errorf("missing Secret")
	}
	userID, ok := value["userID"]
	if !ok {
		return nil, fmt.Errorf("missing userid")
	}
	return &authToken{
		Secret: secret,
		UserID: userID,
	}, nil
}
