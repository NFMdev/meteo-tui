package cache

import (
	"fmt"
	"strings"
	"unicode"
)

func ForecastCacheKey(city string, country string) (string, error) {
	city = strings.TrimSpace(city)
	country = strings.ToLower(strings.TrimSpace(country))

	if city == "" {
		return "", ErrCacheCityRequired
	}
	if country == "" {
		return "", ErrCacheCountryRequired
	}
	if !isValidCountryCode(country) {
		return "", ErrInvalidCacheCountry
	}

	citySegment := sanitizeCacheSegment(city)
	if citySegment == "" {
		return "", ErrInvalidCacheKey
	}

	return fmt.Sprintf("%s_%s", citySegment, country), nil
}

func sanitizeCacheSegment(value string) string {
	value = strings.TrimSpace(value)

	var builder strings.Builder
	previousWasSeparator := false

	for _, char := range value {
		char = unicode.ToLower(char)

		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			builder.WriteRune(char)
			previousWasSeparator = false
			continue
		}

		if !previousWasSeparator {
			builder.WriteRune('_')
			previousWasSeparator = true
		}
	}

	return strings.Trim(builder.String(), "_")
}

func isValidCountryCode(country string) bool {
	if len(country) != 2 {
		return false
	}

	for _, char := range country {
		if char < 'a' || char > 'z' {
			return false
		}
	}

	return true
}
