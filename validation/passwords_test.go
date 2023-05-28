package validation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/bcrypt"
)

type hashPasswordResult struct {
	val string
	err error
}

func TestHashPassword(t *testing.T) {
	tests := map[string]struct {
		input string
		want  hashPasswordResult
	}{
		"empty":  {input: "", want: hashPasswordResult{"", nil}},
		"simple": {input: "simple", want: hashPasswordResult{"", nil}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, got := HashPassword(tc.input)
			diff := cmp.Equal(got, tc.want)
			if diff {
				t.Fatalf("%t", diff)
			}
		})
	}
}

type checkPasswordParams struct {
	password string
	hash     string
}

func mockGetHash(input string) string {
  // We don't care about security here, makes it faster with lower cost.
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(input), 1)
	return string(hashedBytes)
}

func TestCheckPasswordHash(t *testing.T) {
	tests := map[string]struct {
		input checkPasswordParams
		want  bool
	}{
		"empty":               {input: checkPasswordParams{password: "", hash: mockGetHash("")}, want: true},
		"correct":             {input: checkPasswordParams{password: "correct", hash: mockGetHash("correct")}, want: true},
		"incorrect":           {input: checkPasswordParams{password: "correct", hash: mockGetHash("incorrect")}, want: false},
		"wrong_casing":        {input: checkPasswordParams{password: "Correct", hash: mockGetHash("correct")}, want: false},
		"correct_casing":      {input: checkPasswordParams{password: "Correct", hash: mockGetHash("Correct")}, want: true},
		"incorrect_dont_trim": {input: checkPasswordParams{password: "  password    ", hash: mockGetHash("password")}, want: false},
		"correct_dont_trim":   {input: checkPasswordParams{password: "  password ", hash: mockGetHash("  password ")}, want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := CheckPasswordHash(tc.input.password, tc.input.hash)
			if got != tc.want {
				t.Fatalf("%t", got)
			}
		})
	}
}

