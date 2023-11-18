package service

import (
	"context"

	"github.com/FlashpointProject/CommunityWebsite/types"
	"github.com/FlashpointProject/CommunityWebsite/utils"
)

func (s *Service) SearchGotdSuggestions(ctx context.Context, query *types.GotdSuggestionsSearchQuery, fpfss types.IFpfss) ([]*types.GotdSuggestion, int64, error) {
	dbs, err := s.pgdal.NewSession(ctx)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, 0, dberr(err)
	}
	defer dbs.Rollback()

	reports, total, err := s.pgdal.SearchGotdSuggestions(dbs, query, fpfss)
	if err != nil {
		utils.LogCtx(ctx).Error(err)
		return nil, 0, dberr(err)
	}

	externalReports := make([]*types.GotdSuggestion, 0)
	for _, report := range reports {
		externalReports = append(externalReports, report.ToExternal())
	}

	return externalReports, total, nil
}
