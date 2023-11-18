package service

import (
	"context"

	"github.com/FlashpointProject/CommunityWebsite/types"
	"github.com/FlashpointProject/CommunityWebsite/utils"
)

type authToken struct {
	Secret string
	UserID string
}

func (s *Service) SaveUser(ctx context.Context, fpfssUser *types.FPFSSProfile, ipAddr string) (*authToken, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}
	defer dbs.Rollback()

	authToken, err := s.authTokenProvider.CreateAuthToken(fpfssUser.ID)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, err
	}

	if err = s.pgdal.SaveRoles(dbs, fpfssUser.Roles); err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	// Update cached roles
	s.CacheRoles(fpfssUser.Roles)

	roles := make([]string, len(fpfssUser.Roles))
	for i, role := range fpfssUser.Roles {
		roles[i] = role.ID
	}
	if err = s.pgdal.SaveUser(dbs, fpfssUser.ID, fpfssUser.Username, fpfssUser.AvatarURL, roles); err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	if err = s.pgdal.StoreSession(dbs, authToken.Secret, fpfssUser.ID, s.sessionExpirationSeconds, ipAddr); err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	if err := dbs.Commit(); err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, dberr(err)
	}

	return authToken, nil
}

func MapAuthToken(token *authToken) map[string]string {
	return map[string]string{"Secret": token.Secret, "userID": token.UserID}
}
