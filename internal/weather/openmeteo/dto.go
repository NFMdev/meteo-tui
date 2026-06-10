package openmeteo

type forecastResponseDTO struct {
	Latitude             float64    `json:"latitude"`
	Longitude            float64    `json:"longitude"`
	GenerationTimeMS     float64    `json:"generationtime_ms"`
	UTCOffsetSeconds     int        `json:"utc_offset_seconds"`
	Timezone             string     `json:"timezone"`
	TimezoneAbbreviation string     `json:"timezone_abbreviation"`
	Current              currentDTO `json:"current"`
	Hourly               hourlyDTO  `json:"hourly"`
	Daily                dailyDTO   `json:"daily"`
}

type currentDTO struct {
	Time                string  `json:"time"`
	Interval            int     `json:"interval"`
	Temperature2m       float64 `json:"temperature_2m"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	RelativeHumidity2m  int     `json:"relative_humidity_2m"`
	Precipitation       float64 `json:"precipitation"`
	WeatherCode         int     `json:"weather_code"`
	CloudCover          int     `json:"cloud_cover"`
	PressureMSL         float64 `json:"pressure_msl"`
	WindSpeed10m        float64 `json:"wind_speed_10m"`
	WindDirection10m    int     `json:"wind_direction_10m"`
}

type hourlyDTO struct {
	Time                []string  `json:"time"`
	Interval            []int     `json:"interval"`
	Temperature2m       []float64 `json:"temperature_2m"`
	ApparentTemperature []float64 `json:"apparent_temperature"`
	RelativeHumidity2m  []int     `json:"relative_humidity_2m"`
	Precipitation       []float64 `json:"precipitation"`
	WeatherCode         []int     `json:"weather_code"`
	CloudCover          []int     `json:"cloud_cover"`
	PressureMSL         []float64 `json:"pressure_msl"`
	WindSpeed10m        []float64 `json:"wind_speed_10m"`
}

type dailyDTO struct {
	Time             []string  `json:"time"`
	WeatherCode      []int     `json:"weather_code"`
	Temperature2mMax []float64 `json:"temperature_2m_max"`
	Temperature2mMin []float64 `json:"temperature_2m_min"`
	PrecipitationSum []float64 `json:"precipitation_sum"`
	WindSpeed10mMax  []float64 `json:"wind_speed_10m_max"`
}
