package validator

import "testing"

func TestNotBlank(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "not blank",
			value:    "foobar",
			expected: true,
		},
		{
			name:     "blank",
			value:    "",
			expected: false,
		},
		{
			name:     "whitspace",
			value:    "\t\n",
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := NotBlank(tt.value)
			if result != tt.expected {
				t.Fatalf("expected NotBlank(%q) to return %t but got %t", tt.value, tt.expected, result)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		n        int
		expected bool
	}{
		{
			name:     "equal",
			value:    "foo",
			n:        3,
			expected: true,
		},
		{
			name:     "greater",
			value:    "foo",
			n:        2,
			expected: true,
		},
		{
			name:     "less",
			value:    "foo",
			n:        4,
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := MinLength(tt.value, tt.n)
			if result != tt.expected {
				t.Fatalf("expected MinLength(%q, %d) to return %t but got %t", tt.value, tt.n, tt.expected, result)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		n        int
		expected bool
	}{
		{
			name:     "equal",
			value:    "foo",
			n:        3,
			expected: true,
		},
		{
			name:     "less",
			value:    "foo",
			n:        4,
			expected: true,
		},
		{
			name:     "greater",
			value:    "foo",
			n:        2,
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxLength(tt.value, tt.n)
			if result != tt.expected {
				t.Fatalf("expected MaxLength(%q, %d) to return %t but got %t", tt.value, tt.n, tt.expected, result)
			}
		})
	}
}

func TestValidEmail(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "valid",
			value:    "test@email.com",
			expected: true,
		},
		{
			name:     "invalid",
			value:    "email.com",
			expected: false,
		},
		{
			name:     "blank",
			value:    "email.com",
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidEmail(tt.value)
			if result != tt.expected {
				t.Fatalf("expected ValidEmail(%q) to return %t but got %t", tt.value, tt.expected, result)
			}
		})
	}
}

func TestHTTPUrl(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "valid http url",
			value:    "http://example.com",
			expected: true,
		},
		{
			name:     "valid https url with path and query",
			value:    "https://www.google.com/search?q=go+lang",
			expected: true,
		},
		{
			name:     "valid http localhost with port",
			value:    "http://localhost:8080",
			expected: true,
		},
		{
			name:     "valid https ip address",
			value:    "https://192.168.1.1/path",
			expected: true,
		},
		{
			name:     "invalid ftp scheme",
			value:    "ftp://example.com",
			expected: false,
		},
		{
			name:     "invalid missing scheme",
			value:    "example.com",
			expected: false,
		},
		{
			name:     "invalid malformed url",
			value:    "invalid-url",
			expected: false,
		},
		{
			name:     "invalid empty string",
			value:    "",
			expected: false,
		},
		{
			name:     "invalid http empty host",
			value:    "http://",
			expected: false,
		},
		{
			name:     "invalid https empty host",
			value:    "https://",
			expected: false,
		},
		{
			name:     "valid http with username password",
			value:    "http://user:pass@host.com/path",
			expected: true,
		},
		{
			name:     "valid https with fragment",
			value:    "https://example.com/path#section",
			expected: true,
		},
		{
			name:     "invalid just a scheme",
			value:    "http:",
			expected: false,
		},
		{
			name:     "invalid just a domain",
			value:    "www.example.com", // Missing scheme
			expected: false,
		},
		{
			name:     "valid complex path and query",
			value:    "https://sub.domain.co.uk/a/b/c?param1=value1&param2=value2#fragment",
			expected: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := HTTPUrl(tt.value)
			if result != tt.expected {
				t.Fatalf("expected HTTPUrl(%q) to return %t but got %t", tt.value, tt.expected, result)
			}
		})
	}
}
