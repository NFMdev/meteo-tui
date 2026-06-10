package tui

import "strings"

func truncateText(value string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	runes := []rune(value)
	if len(runes) <= maxWidth {
		return value
	}

	if maxWidth <= 1 {
		return "…"
	}

	return string(runes[:maxWidth-1]) + "…"
}

func truncateLines(lines []string, maxWidth int) []string {
	truncated := make([]string, 0, len(lines))

	for _, line := range lines {
		truncated = append(truncated, truncateText(line, maxWidth))
	}

	return truncated
}

func joinTruncatedLines(lines []string, maxWidth int) string {
	return strings.Join(truncateLines(lines, maxWidth), "\n")
}
