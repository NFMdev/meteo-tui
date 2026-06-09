package location

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/nfmdev/meteo/internal/domain"
)

/*										  https://geocoding-api.open-meteo.com/v1/search?name=Berlin&count=10&language=en&format=json*/
const defaultOpenMeteoGeocodingBaseURL = "https://geocoding-api.open-meteo.com/v1/search"

type OpenMeteoResolver struct {
	client  *http.Client
	baseURL string
}

func NewOpenMeteoResolver(client *http.Client) OpenMeteoResolver {
	if client == nil {
		client = http.DefaultClient
	}

	return OpenMeteoResolver{
		client:  client,
		baseURL: defaultOpenMeteoGeocodingBaseURL,
	}
}

func (r OpenMeteoResolver) Resolve(
	ctx context.Context,
	city string,
	country string,
) (domain.Location, error) {
	city = strings.TrimSpace(city)
	country = strings.ToUpper(strings.TrimSpace(country))

	if city == "" {
		return domain.Location{}, ErrCityRequired
	}

	if country == "" {
		return domain.Location{}, ErrCountryRequired
	}

	requestURL, err := buildGeocodingURL(r.baseURL, city)
	if err != nil {
		return domain.Location{}, fmt.Errorf("build geocoding URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return domain.Location{}, fmt.Errorf("create geocoding request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return domain.Location{}, fmt.Errorf("execute geocoding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return domain.Location{}, fmt.Errorf(
			"geocoding request failed with stats %d",
			resp.StatusCode,
		)
	}

	var dto geocodingResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return domain.Location{}, fmt.Errorf("decode geocoding response: %w", err)
	}

	for _, result := range dto.Results {
		if strings.EqualFold(result.CountryCode, country) {
			return mapGeocodingResult(result), nil
		}
	}

	return domain.Location{}, ErrLocationNotFound
}

func buildGeocodingURL(baseUrl string, city string) (string, error) {
	parsedURL, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	query := parsedURL.Query()
	query.Set("name", city)
	query.Set("count", "10")
	query.Set("language", "en")
	query.Set("format", "json")

	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

func mapGeocodingResult(result geocodingResultDTO) domain.Location {
	return domain.Location{
		City:      result.Name,
		Country:   result.CountryCode,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Timezone:  result.Timezone,
	}
}
