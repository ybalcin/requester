package utility

import "testing"

func TestStrLength(t *testing.T) {
	t.Run("whitespace str", func(t *testing.T) {
		l := StrLength(" ")

		if l != 0 {
			t.Errorf("expected str length is 0, but got: %d", l)
		}
	})

	t.Run("empty str", func(t *testing.T) {
		l := StrLength("")

		if l != 0 {
			t.Errorf("expected str length is 0, but got: %d", l)
		}
	})
}

func TestIsStrEmpty(t *testing.T) {
	t.Run("whitespace str", func(t *testing.T) {
		empty := IsStrEmpty(" ")

		if !empty {
			t.Errorf("expected str isEmpty value is true, but got false")
		}
	})

	t.Run("empty str", func(t *testing.T) {
		empty := IsStrEmpty("")

		if !empty {
			t.Errorf("expected str isEmpty value is true, but got false")
		}
	})
}

func TestIsSchemeExistInURL(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected bool
	}{
		{
			name:     "with http",
			address:  "http://google.com",
			expected: true,
		},
		{
			name:     "with https",
			address:  "https://google.com",
			expected: true,
		},
		{
			name:     "without scheme",
			address:  "google.com",
			expected: false,
		},
		{
			name:     "wrong scheme",
			address:  "htt://google.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IsSchemeExistInURL(tt.address)
			if actual != tt.expected {
				t.Errorf("expected %v, but got: %v", tt.expected, actual)
			}
		})
	}
}
