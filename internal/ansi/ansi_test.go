package ansi

import (
	"testing"
)

func TestLink(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		text     string
		expected string
	}{
		{
			name:     "basic link",
			url:      "https://example.com",
			text:     "Example",
			expected: "\033]8;;https://example.com\033\\Example\033]8;;\033\\",
		},
		{
			name:     "empty text",
			url:      "https://example.com",
			text:     "",
			expected: "\033]8;;https://example.com\033\\\033]8;;\033\\",
		},
		{
			name:     "empty url",
			url:      "",
			text:     "Example",
			expected: "\033]8;;\033\\Example\033]8;;\033\\",
		},
		{
			name:     "special characters in text",
			url:      "https://example.com",
			text:     "Test & Example",
			expected: "\033]8;;https://example.com\033\\Test & Example\033]8;;\033\\",
		},
		{
			name:     "url with query parameters",
			url:      "https://example.com?foo=bar&baz=qux",
			text:     "Query Link",
			expected: "\033]8;;https://example.com?foo=bar&baz=qux\033\\Query Link\033]8;;\033\\",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Link(tt.url, tt.text)
			if result != tt.expected {
				t.Errorf("Link(%q, %q) = %q, want %q", tt.url, tt.text, result, tt.expected)
			}
		})
	}
}