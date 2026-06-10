package openmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/nfmdev/meteo/internal/domain"
)

const defaultForecastBaseURL = "https://api.open-meteo.com/v1/forecast"

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient(client *http.Client) Client {
	if client == nil {
		client = http.DefaultClient
	}

	return Client{
		client:  client,
		baseURL: defaultForecastBaseURL,
	}
}

func (c Client) GetForecast(
	ctx context.Context,
	location domain.Location,
) (domain.WeatherReport, error) {
	requestURL, err := buildForecastURL(c.baseURL, location)
	if err != nil {
		return domain.WeatherReport{}, fmt.Errorf("build forecast URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return domain.WeatherReport{}, fmt.Errorf("create forecast request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return domain.WeatherReport{}, fmt.Errorf("execute forecast request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return domain.WeatherReport{}, fmt.Errorf(
			"%w: status %d",
			ErrForecastUnavailable,
			resp.StatusCode,
		)
	}

	var dto forecastResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return domain.WeatherReport{}, fmt.Errorf("decode forecast response: %w", err)
	}

	report, err := mapForecastResponse(location, dto)
	if err != nil {
		return domain.WeatherReport{}, err
	}

	return report, nil
}

func buildForecastURL(baseURL string, location domain.Location) (string, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	timezone := location.Timezone
	if timezone == "" {
		timezone = "auto"
	}

	query := parsedURL.Query()
	query.Set("latitude", strconv.FormatFloat(location.Latitude, 'f', 6, 64))
	query.Set("longitude", strconv.FormatFloat(location.Longitude, 'f', 6, 64))
	query.Set("forecast_days", "7")
	query.Set("timezone", timezone)

	query.Set("current", joinVariables([]string{
		"temperature_2m",
		"apparent_temperature",
		"relative_humidity_2m",
		"precipitation",
		"weather_code",
		"cloud_cover",
		"pressure_msl",
		"wind_speed_10m",
		"wind_direction_10m",
	}))

	query.Set("hourly", joinVariables([]string{
		"temperature_2m",
		"apparent_temperature",
		"relative_humidity_2m",
		"precipitation",
		"weather_code",
		"cloud_cover",
		"pressure_msl",
		"wind_speed_10m",
	}))

	query.Set("daily", joinVariables([]string{
		"weather_code",
		"temperature_2m_max",
		"temperature_2m_min",
		"precipitation_sum",
		"wind_speed_10m_max",
	}))

	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

func joinVariables(values []string) string {
	var result strings.Builder

	for index, value := range values {
		if index > 0 {
			result.WriteString(",")
		}
		result.WriteString(value)
	}

	return result.String()
}
