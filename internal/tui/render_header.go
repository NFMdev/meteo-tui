package tui

import (
	"fmt"
	"strings"
)

func (m Model) renderHeader() string {
	location := m.report.Location

	title := titleStyle.Render(
		truncateText(
			fmt.Sprintf("Meteo — %s, %s", location.City, location.Country),
			m.contentWidth(),
		),
	)

	metadata := subtitleStyle.Render(
		truncateText(
			fmt.Sprintf(
				"Updated: %s • %.4f, %.4f • %s",
				m.report.UpdatedAt.Format("15:04:05"),
				location.Latitude,
				location.Longitude,
				location.Timezone,
			),
			m.contentWidth(),
		),
	)

	return strings.Join([]string{
		title,
		metadata,
	}, "\n")
}
