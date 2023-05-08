package validation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type result struct {
	val string
	err error
}

func TestHashPassword(t *testing.T) {
	tests := map[string]struct {
		input string
		want  result
	}{
		"empty":  {input: "", want: result{"", nil}},
		"simple": {input: "simple", want: result{"", nil}},
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
