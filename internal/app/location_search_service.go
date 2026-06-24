package app

import (
	"context"
	"strings"

	"github.com/nfmdev/meteo/internal/domain"
)

type LocationSearchService interface {
	SearchLocations(
		ctx context.Context,
		query string,
	) ([]domain.LocationSearchResult, error)
}

type LocationSearcher interface {
	Search(
		ctx context.Context,
		query string,
	) ([]domain.LocationSearchResult, error)
}

type RealLocationSearchService struct {
	searcher LocationSearcher
}

func NewLocationSearchService(
	searcher LocationSearcher,
) RealLocationSearchService {
	return RealLocationSearchService{
		searcher: searcher,
	}
}

func (s RealLocationSearchService) SearchLocations(
	ctx context.Context,
	query string,
) ([]domain.LocationSearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, ErrLocationSearchQueryRequired
	}
	if s.searcher == nil {
		return nil, ErrLocationSearcherRequired
	}

	return s.searcher.Search(ctx, query)
}
