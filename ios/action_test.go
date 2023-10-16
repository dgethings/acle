package ios

import (
	"fmt"
	"testing"
)

func TestParseAction_happy_path(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "permit", want: "permit"},
		{input: "deny", want: "deny"},
	}

	for i, tc := range tests {
		got, err := parseAction(tc.input)
		if got.String() != tc.want {
			t.Errorf("test %d: wanted: '%s', got: '%s', error: '%v'", i+1, tc.want, got, err)
		}
	}
}

func TestParseAction_sad_path(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "", want: "no action field found"},
		{input: "bogus", want: "unrecognised action bogus"},
	}

	for i, tc := range tests {
		got, err := parseAction(tc.input)
		if fmt.Sprintf("%v", err) != tc.want {
			t.Errorf("test %d: wanted: '%s', got: '%s', error: '%v'", i+1, tc.want, got, err)
		}
	}
}
