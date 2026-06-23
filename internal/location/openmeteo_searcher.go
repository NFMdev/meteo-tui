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

const defaultOpenMeteoSearchBaseURL = "https://geocoding-api.open-meteo.com/v1/search"

type OpenMeteoSearcher struct {
	client  *http.Client
	baseURL string
}

func NewOpenMeteoSearcher(client *http.Client) OpenMeteoSearcher {
	if client == nil {
		client = http.DefaultClient
	}

	return OpenMeteoSearcher{
		client:  client,
		baseURL: defaultOpenMeteoGeocodingBaseURL,
	}
}

func (s OpenMeteoSearcher) Search(
	ctx context.Context,
	query string,
) ([]domain.LocationSearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, ErrSearchQueryRequired
	}

	requestURL, err := s.buildSearchURL(query)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create Open-Meteo location search request: %w", err)
	}

	response, err := s.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("execute Open-Meteo location search request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"Open-Meteo location search returned status %d",
			response.StatusCode,
		)
	}

	var dto openMeteoSearchResponseDTO
	if err := json.NewDecoder(response.Body).Decode(&dto); err != nil {
		return nil, fmt.Errorf("decode Open-Meteo location search response: %w", err)
	}

	return mapOpenMeteoSearchResult(dto), nil
}

func (s OpenMeteoSearcher) buildSearchURL(query string) (string, error) {
	baseURL, err := url.Parse(s.baseURL)
	if err != nil {
		return "", fmt.Errorf("parse Open-Meteo location search result URL: %w", err)
	}

	params := baseURL.Query()
	params.Set("name", query)
	params.Set("count", "10")
	params.Set("language", "en")
	params.Set("format", "json")

	baseURL.RawQuery = params.Encode()

	return baseURL.String(), nil
}

func mapOpenMeteoSearchResult(dto openMeteoSearchResponseDTO) []domain.LocationSearchResult {
	results := make([]domain.LocationSearchResult, 0, len(dto.Results))

	for _, result := range dto.Results {
		name := strings.TrimSpace(result.Name)
		countryCode := strings.ToUpper(strings.TrimSpace(result.CountryCode))

		if name == "" || countryCode == "" {
			continue
		}

		results = append(results, domain.LocationSearchResult{
			Name:        name,
			Country:     strings.TrimSpace(result.Country),
			CountryCode: countryCode,
			Admin1:      result.Admin1,
			Latitude:    result.Latitude,
			Longitude:   result.Longitude,
			Timezone:    result.Timezone,
		})
	}

	return results
}
