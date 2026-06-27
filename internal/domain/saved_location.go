package domain

import "strings"

type SavedLocation struct {
	Name        string  `json:"name"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Admin1      string  `json:"admin1"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
}

func SavedLocationFromSearchResult(result LocationSearchResult) SavedLocation {
	return NormalizeSavedLocation(SavedLocation{
		Name:        result.Name,
		Country:     result.Country,
		CountryCode: result.CountryCode,
		Admin1:      result.Admin1,
		Latitude:    result.Latitude,
		Longitude:   result.Longitude,
		Timezone:    result.Timezone,
	})
}

func SavedLocationFromLocation(location Location) SavedLocation {
	return NormalizeSavedLocation(SavedLocation{
		Name:        location.City,
		CountryCode: location.Country,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		Timezone:    location.Timezone,
	})
}

func NormalizeSavedLocation(location SavedLocation) SavedLocation {
	return SavedLocation{
		Name:        strings.TrimSpace(location.Name),
		Country:     strings.TrimSpace(location.Country),
		CountryCode: strings.ToUpper(strings.TrimSpace(location.CountryCode)),
		Admin1:      strings.TrimSpace(location.Admin1),
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		Timezone:    strings.TrimSpace(location.Timezone),
	}
}

func SavedLocationKey(location SavedLocation) string {
	location = NormalizeSavedLocation(location)

	name := normalizeSavedLocationKeyPart(location.Name)
	countryCode := normalizeSavedLocationKeyPart(location.CountryCode)

	if name == "" || countryCode == "" {
		return ""
	}

	return name + "_" + countryCode
}

func SameSavedLocation(loc1 SavedLocation, loc2 SavedLocation) bool {
	return SavedLocationKey(loc1) != "" && SavedLocationKey(loc1) == SavedLocationKey(loc2)
}

func normalizeSavedLocationKeyPart(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))

	var builder strings.Builder
	previousWasSeparator := false

	for _, char := range value {
		if (char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			(char >= 'à' && char <= 'ÿ') {
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
