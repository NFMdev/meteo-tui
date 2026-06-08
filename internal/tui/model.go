package tui

type Model struct {
	city    string
	country string

	width  int
	height int
}

func NewModel(city string, country string) Model {
	return Model{
		city:    city,
		country: country,
	}
}
