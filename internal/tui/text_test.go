package tui

import "testing"

func TestTruncateTextKeepsShortText(t *testing.T) {
	t.Parallel()

	got := truncateText("Copenhagen", 12)

	if got != "Copenhagen" {
		t.Fatalf("expected Copenhagen, got %q", got)
	}
}

func TestTruncateTextTruncatesLongText(t *testing.T) {
	t.Parallel()

	got := truncateText("Europe/Copenhagen", 8)

	if got != "Europe/…" {
		t.Fatalf("expected Europe/…, got %q", got)
	}
}

func TestTruncateTextHandlesZeroWidth(t *testing.T) {
	t.Parallel()

	got := truncateText("Copenhagen", 0)

	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestTruncateTextHandlesWidthOne(t *testing.T) {
	t.Parallel()

	got := truncateText("Copenhagen", 1)

	if got != "…" {
		t.Fatalf("expected ellipsis, got %q", got)
	}
}

func TestTruncateLines(t *testing.T) {
	t.Parallel()

	lines := []string{
		"short",
		"very long line",
	}

	got := truncateLines(lines, 6)

	if got[0] != "short" {
		t.Fatalf("expected first line short, got %q", got[0])
	}

	if got[1] != "very …" {
		t.Fatalf("expected second line very …, got %q", got[1])
	}
}
