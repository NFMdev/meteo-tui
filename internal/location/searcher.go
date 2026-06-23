package location

import (
	"context"

	"github.com/nfmdev/meteo/internal/domain"
)

type Seacher interface {
	Search(ctx context.Context, query string) ([]domain.LocationSearchResult, error)
}
