package location

type openMeteoSearchResponseDTO struct {
	Results []openMeteoSearchResultDTO `json:"results"`
}

type openMeteoSearchResultDTO struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Elevation   float64 `json:"elevation"`
	FeatureCode string  `json:"feature_code"`
	CountryCode string  `json:"country_code"`
	Admin1ID    int     `json:"admin1_id"`
	Admin1      string  `json:"admin1"`
	CountryID   int     `json:"country_id"`
	Country     string  `json:"country"`
	Timezone    string  `json:"timezone"`
	Population  int     `json:"population"`
}
