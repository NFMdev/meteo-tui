package location

type geocodingResponseDTO struct {
	Results []geocodingResultDTO `json:"results"`
}

type geocodingResultDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Elevation   float64 `json:"elevation"`
	FeatureCode string  `json:"feature_code"`
	CountryCode string  `json:"country_code"`
	Admin1ID    int     `json:"admin1_id"`
	Admin2ID    int     `json:"admin2_id"`
	Timezone    string  `json:"timezone"`
	Population  int     `json:"population"`
	CountryID   int     `json:"country_id"`
	Country     string  `json:"country"`
	Admin1      string  `json:"admin1"`
	Admin2      string  `json:"admin2"`
}
