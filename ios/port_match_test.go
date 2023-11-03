package ios

import (
	"testing"
)

func TestParsePort(t *testing.T) {
	type test struct {
		input []string
		want  string
	}

	tests := []test{
		{input: []string{"any", "any"}, want: ""},
		{input: []string{"eq", "bgp", "foo"}, want: "eq bgp"},
		{input: []string{"gt", "1024"}, want: "gt 1024"},
		{input: []string{"lt", "1024"}, want: "lt 1024"},
		{input: []string{"range", "ftp", "ftp-data"}, want: "range ftp ftp-data"},
		{input: []string{"eq", "179"}, want: "eq bgp"},
	}

	for i, tc := range tests {
		m, _, e := parsePort(tc.input)
		if e != nil {
			t.Errorf("failed to parse input: %s, got: %v", tc.input, e)
		}
		if m.String() != tc.want {
			t.Errorf("test %d: given: %s, expected: %s, got: %s", i+1, tc.input, tc.want, m)
		}
	}
}
