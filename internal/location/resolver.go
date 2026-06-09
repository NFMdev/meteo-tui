package location

import (
	"context"

	"github.com/nfmdev/meteo/internal/domain"
)

type Resolver interface {
	Resolve(ctx context.Context, city string, country string) (domain.Location, error)
}
