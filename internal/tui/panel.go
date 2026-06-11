package tui

import "strings"

func renderPanel(title string, lines []string, width int, height int) string {
	contentLines := []string{
		title,
		"",
	}

	contentLines = append(contentLines, lines...)

	content := joinTruncatedLines(contentLines, innerWidthFor(width))

	return panelStyle.
		Width(width).
		Height(height).
		Render(content)
}

func renderPanelWithRawContent(content string, width int, height int) string {
	return panelStyle.
		Width(width).
		Height(height).
		Render(content)
}

func innerWidthFor(panelWidth int) int {
	width := panelWidth - panelHorizontalFrame
	if width < 1 {
		return 1
	}

	return width
}

func limitLines(lines []string, maxLines int) []string {
	if maxLines <= 0 {
		return nil
	}

	if len(lines) <= maxLines {
		return lines
	}

	limited := make([]string, 0, maxLines)
	limited = append(limited, lines[:maxLines-1]...)
	limited = append(limited, "…")

	return limited
}

func blank(width int) string {
	if width <= 0 {
		return ""
	}

	return strings.Repeat(" ", width)
}
