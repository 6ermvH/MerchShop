package hasher

import (
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	testCases := []struct {
		name    string
		hasPwd  string
		isError bool
	}{
		{
			name:    "Simple Password",
			hasPwd:  "a",
			isError: false,
		},
		{
			name:    "Empty Password",
			hasPwd:  "",
			isError: false,
		},
		{
			name:    "Ascii Password",
			hasPwd:  "isGoodPassword",
			isError: false,
		},
		{
			name:    "UTF-8 Password",
			hasPwd:  "isUtf😜",
			isError: false,
		},
		{
			name:    "Long Password",
			hasPwd:  strings.Repeat("a", 73),
			isError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if hash, err := HashPassword(test.hasPwd); (err == nil) == test.isError {
				t.Fatalf("password: %q, hashed to %q, error: %v", test.hasPwd, hash, err)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	testCases := []struct {
		name    string
		pwd     string
		hash    string
		isError bool
	}{
		{
			name:    "Simple",
			pwd:     "isGoodPassword",
			hash:    hash("isGoodPassword"),
			isError: false,
		},
		{
			name:    "Bad hash",
			pwd:     "isGoodPassword",
			hash:    hash("isGoodPassword1"),
			isError: true,
		},
		{
			name:    "Bad hash",
			pwd:     "isGoodPassword",
			hash:    hash("IsGoodPassword"),
			isError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if err := CheckPassword(test.hash, test.pwd); (err == nil) == test.isError {
				t.Fatalf("password: %q, hash: %q, error: %v", test.pwd, test.hash, err)
			}
		})
	}
}

func hash(pwd string) string {
	if hash, err := HashPassword(pwd); err == nil {
		return hash
	}
	return ""
}
