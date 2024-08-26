package utilities

import (
	"reflect"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove trailing slash",
			inputURL: "http://example.com/path/",
			expected: "example.com/path",
		},
		{
			name:     "process relative path",
			inputURL: "sub.example.com/path/to/resource/",
			expected: "sub.example.com/path/to/resource",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NormalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
    <body>
        <a href="/path/one">
            <span>Boot.dev</span>
        </a>
        <a href="https://other.com/path/one">
            <span>Boot.dev</span>
        </a>
    </body>
</html>
                `,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "no URLs",
			inputURL: "https://example.com",
			inputBody: `
<html>
    <body>
        <p>No links here!</p>
    </body>
</html>
            `,
			expected: []string{},
		},
		{
			name:     "relative URL with trailing slash",
			inputURL: "https://example.com",
			inputBody: `
<html>
    <body>
        <a href="/path/one/">
            <span>Boot.dev</span>
        </a>
    </body>
</html>
            `,
			expected: []string{"https://example.com/path/one/"},
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := GetURLsFromHTML(test.inputBody, test.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %s", i, test.name, err)
			}

			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected slice: %v, actual: %v", i, test.name, test.expected, actual)
			}
		})
	}
}
