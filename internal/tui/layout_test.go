package tui

import "testing"

func TestIsTerminalTooSmall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		width    int
		height   int
		expected layoutMode
	}{
		{
			name:     "unknown size is not considered too small",
			width:    0,
			height:   0,
			expected: layoutModeCompactScrollable,
		},
		{
			name:     "normal terminal",
			width:    100,
			height:   30,
			expected: layoutModeGrid,
		},
		{
			name:     "too narrow",
			width:    minTerminalWidth - 1,
			height:   minTerminalHeight,
			expected: layoutModeTooSmall,
		},
		{
			name:     "too short",
			width:    minTerminalWidth,
			height:   minTerminalHeight - 1,
			expected: layoutModeTooSmall,
		},
		{
			name:     "minimum size",
			width:    minTerminalWidth,
			height:   minTerminalHeight,
			expected: layoutModeCompactScrollable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := Model{
				width:  tt.width,
				height: tt.height,
			}

			got := model.layoutMode()

			if got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestContentWidthUsesDefaultWhenUnknown(t *testing.T) {
	t.Parallel()

	model := Model{}

	got := model.contentWidth()

	if got != defaultTerminalWidth-appHorizontalPadding {
		t.Fatalf("unexpected content width: %d", got)
	}
}

func TestContentWidthUsesTerminalWidth(t *testing.T) {
	t.Parallel()

	model := Model{
		width: 100,
	}

	got := model.contentWidth()

	if got != 96 {
		t.Fatalf("expected content width 96, got %d", got)
	}
}
