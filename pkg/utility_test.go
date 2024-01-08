package pkg

import (
	"strings"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	length := 10
	str := GenerateRandomString(length)
	if len(str) != length {
		t.Errorf("Expected string of length %d, got %d", length, len(str))
	}

	if !strings.ContainsAny(str, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") {
		t.Errorf("String contains unexpected characters")
	}
}

func TestIsValidURL(t *testing.T) {
	testCases := []struct {
		urlStr string
		valid  bool
	}{
		{"https://example1.com", true},
		{"example.com", false},
	}

	for _, tc := range testCases {
		if IsValidURL(tc.urlStr) != tc.valid {
			t.Errorf("IsValidURL(%s) expected to be %v", tc.urlStr, tc.valid)
		}
	}
}

func TestGetUniqueShortUrl(t *testing.T) {
	existingUrls := []string{"abc12", "xyz78"}
	shortUrl := GetUniqueShortUrl(existingUrls, 5)

	if shortUrl == "abc12" || shortUrl == "xyz78" {
		t.Errorf("Generated short URL should be unique")
	}
}
