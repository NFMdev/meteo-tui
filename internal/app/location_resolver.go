package app

import (
	"context"

	"github.com/nfmdev/meteo/internal/domain"
)

type LocationResolver interface {
	Resolve(ctx context.Context, city string, country string) (domain.Location, error)
}
