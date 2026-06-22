package cache

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestForecastStoreWriteAndReadForecast(t *testing.T) {
	t.Parallel()

	cacheDir := t.TempDir()
	store := newTestForecastStore(t, cacheDir)

	fixedTime := time.Date(2026, 6, 10, 12, 30, 0, 0, time.UTC)
	store.now = func() time.Time {
		return fixedTime
	}

	report := testWeatherReport()

	err := store.WriteForecast("Copenhagen", "DK", report)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	entry, err := store.ReadForecast("Copenhagen", "DK")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if entry.Key != "copenhagen_dk" {
		t.Fatalf("expected key copenhagen_dk, got %q", entry.Key)
	}

	if entry.City != "Copenhagen" {
		t.Fatalf("expected city Copenhagen, got %q", entry.City)
	}

	if entry.Country != "DK" {
		t.Fatalf("expected country DK, got %q", entry.Country)
	}

	if !entry.CachedAt.Equal(fixedTime) {
		t.Fatalf("expected cached_at %v, got %v", fixedTime, entry.CachedAt)
	}

	if entry.Report.Source.Provider != domain.WeatherProviderOpenMeteo {
		t.Fatalf("expected provider Open-Meteo, got %q", entry.Report.Source.Provider)
	}

	if !entry.Report.Source.Cached {
		t.Fatal("expected cached report source")
	}

	if !entry.Report.Source.CachedAt.Equal(fixedTime) {
		t.Fatalf("expected source cached_at %v, got %v", fixedTime, entry.Report.Source.CachedAt)
	}

	if entry.Report.Location.City != "Copenhagen" {
		t.Fatalf("expected report city Copenhagen, got %q", entry.Report.Location.City)
	}

	if entry.Report.Current.TemperatureC != 18.5 {
		t.Fatalf("expected temperature 18.5, got %.1f", entry.Report.Current.TemperatureC)
	}
}

func TestForecastStoreWriteCreatesParentDirectories(t *testing.T) {
	t.Parallel()

	cacheDir := filepath.Join(t.TempDir(), "nested", "cache", "forecasts")
	store := newTestForecastStore(t, cacheDir)

	err := store.WriteForecast("Copenhagen", "DK", testWeatherReport())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedPath := filepath.Join(cacheDir, "copenhagen_dk.json")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Fatalf("expected cache file to exist at %q, got %v", expectedPath, err)
	}
}

func TestForecastStoreReadReturnsCacheNotFound(t *testing.T) {
	t.Parallel()

	store := newTestForecastStore(t, t.TempDir())

	_, err := store.ReadForecast("Copenhagen", "DK")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrCacheNotFound) {
		t.Fatalf("expected ErrCacheNotFound, got %v", err)
	}
}

func TestForecastStoreReadRejectsMalformedJSON(t *testing.T) {
	t.Parallel()

	cacheDir := t.TempDir()
	path := filepath.Join(cacheDir, "copenhagen_dk.json")

	if err := os.WriteFile(path, []byte(`{ invalid json`), 0o664); err != nil {
		t.Fatalf("failed to write malformed cache file: %v", err)
	}

	store := newTestForecastStore(t, cacheDir)
	_, err := store.ReadForecast("Copenhagen", "DK")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "decode forecast cache file") {
		t.Fatalf("expected decode error, got %v", err)
	}
}

func TestForecastStoreReadRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	cacheDir := t.TempDir()
	path := filepath.Join(cacheDir, "copenhagen_dk.json")

	content := `{
		"key": "copenhagen_dk",
		"city": "Copenhagen",
		"country": "DK",
		"cached_at": "2026-06-10T12:30:00Z",
		"unknown_field": true,
		"report": {
			"Location": {
				"city": "Copenhagen",
				"country": "DK",
				"latitude": 57.048,
				"longitude": 9.9187,
				"timezone": "Europe/Copenhagen"
			},
			"updated_at": "2026-06-10T12:00:00Z",
			"source": {},
			"current": {},
			"metrics": {},
			"daily": [],
			"hourly": []
		}
	}`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write cache file: %v", err)
	}

	store := newTestForecastStore(t, cacheDir)
	_, err := store.ReadForecast("Copenhagen", "DK")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("expected unknown field error, got %v", err)
	}
}

func TestForecastStoreReadRejectsMismatchedKey(t *testing.T) {
	t.Parallel()

	cacheDir := t.TempDir()
	store := newTestForecastStore(t, cacheDir)
	err := store.WriteForecast("Copenhagen", "DK", testWeatherReport())
	if err != nil {
		t.Fatalf("expected no error writing cache, got %v", err)
	}

	path := filepath.Join(cacheDir, "copenhagen_dk.json")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("faild to read cache file: %v", err)
	}

	modified := strings.Replace(string(content), `"key": "copenhagen_dk"`, `"key": "madrid_es"`, 1)
	if err := os.WriteFile(path, []byte(modified), 0o664); err != nil {
		t.Fatalf("failed to write modified cache file: %v", err)
	}

	_, err = store.ReadForecast("Copenhagen", "DK")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidCacheKey) {
		t.Fatalf("expected ErrInvalidCacheEntry, got %v", err)
	}
}

func TestForecastStoreWriteRejectsInvalidLocation(t *testing.T) {
	t.Parallel()

	store := newTestForecastStore(t, t.TempDir())
	err := store.WriteForecast("Copenhagen", "DNK", testWeatherReport())
	if err == nil {
		t.Fatal("expected err, got nil")
	}
	if !errors.Is(err, ErrInvalidCacheCountry) {
		t.Fatalf("expected ErrInvalidCacheCountry, got %v", err)
	}
}

func TestNewForecastStoreRejectsEmptyCacheDir(t *testing.T) {
	t.Parallel()

	_, err := NewForecastStore("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrCacheDirRequired) {
		t.Fatalf("expected ErrCacheDirRequired, got %v", err)
	}
}

func newTestForecastStore(t *testing.T, cacheDir string) ForecastStore {
	t.Helper()

	store, err := NewForecastStore(cacheDir)
	if err != nil {
		t.Fatalf("expected no error creating forecast store, got %v", err)
	}

	return store
}

func testWeatherReport() domain.WeatherReport {
	location := domain.Location{
		City:      "Copenhagen",
		Country:   "DK",
		Latitude:  55.6761,
		Longitude: 12.5683,
		Timezone:  "Europe/Copenhagen",
	}

	return domain.WeatherReport{
		Location:  location,
		UpdatedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		Source:    domain.NewFreshWeatherSource(domain.WeatherProviderOpenMeteo),
		Current: domain.CurrentWeather{
			TemperatureC:     18.5,
			FeelsLikeC:       17.9,
			Condition:        "Partly Cloudy",
			WeatherCode:      2,
			WindSpeedKmh:     18.0,
			WindDirectionDeg: 240,
		},
		Metrics: domain.WeatherMetrics{
			HumidityPercent:   65,
			PressureHPa:       1015.2,
			PrecipitationMM:   0.0,
			CloudCoverPercent: 40,
			WindSpeedKmh:      18.0,
			WindDirectionDeg:  240,
		},
		Daily: []domain.DailyForecast{
			{
				Date:            time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
				MinTemperatureC: 12.5,
				MaxTemperatureC: 20.1,
				Condition:       "Partly cloudy",
				WeatherCode:     2,
				PrecipitationMM: 0.0,
				MaxWindKmh:      22.0,
			},
		},
		Hourly: []domain.HourlyForecast{
			{
				Time:              time.Date(2026, 6, 10, 8, 0, 0, 0, time.UTC),
				TemperatureC:      14.2,
				FeelsLikeC:        13.8,
				Condition:         "Partly cloudy",
				WeatherCode:       2,
				HumidityPercent:   70,
				PrecipitationMM:   0.0,
				CloudCoverPercent: 40,
				WindSpeedKmh:      16.0,
			},
		},
	}
}
