package ios

import (
	"fmt"
	"os"

	// "strings"
	"testing"
)

func TestLoadFilename(t *testing.T) {
	_, err := os.ReadFile("./test_data/sample.ios")
	if err != nil {
		t.Error(err)
	}
}

func TestGetAcl(t *testing.T) {
	b, _ := os.ReadFile("./test_data/sample.ios")
	cfg := string(b)
	_, err := GetACL("103", cfg, "standard")
	if err != nil {
		t.Error(err)
	}
}

func TestParseAction(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "permit", want: "permit"},
		{input: "deny", want: "deny"},
		{input: "", want: ""},
		{input: "bogus", want: ""},
	}

	for i, tc := range tests {
		got, err := parseAction(tc.input)
		if fmt.Sprintf("%s", got) != tc.want {
			t.Errorf("test %d: wanted: %s, got: %s, error: %v", i+1, tc.want, got, err)
		}
	}
}

func TestPrefixFromWildcard(t *testing.T) {
	type test struct {
		input string
		want  int
	}

	tests := []test{
		{input: "0.0.0.0", want: 32},
		{input: "0.0.0.255", want: 24},
		{input: "0.0.255.255", want: 16},
		{input: "0.255.255.255", want: 8},
		{input: "255.255.255.255", want: 0},
	}

	for i, tc := range tests {
		got, err := prefixFromWildcard(tc.input)
		if got != tc.want {
			t.Errorf("test %d: expected: %d, got: %d, error: %v", i+1, tc.want, got, err)
		}
	}
}

func TestWildcardFromPrefix(t *testing.T) {
	type test struct {
		input int
		want  string
	}
	tests := []test{
		{input: 32, want: "0.0.0.0"},
		{input: 24, want: "0.0.0.255"},
		{input: 16, want: "0.0.255.255"},
		{input: 8, want: "0.255.255.255"},
		{input: 0, want: "255.255.255.255"},
	}
	for i, tc := range tests {
		got := wildcardFromPrefix(tc.input)
		if got != tc.want {
			t.Errorf("test %d: given: %d, expected: %s, got: %s", i+1, tc.input, tc.want, got)
		}
	}
}

func TestParsePort(t *testing.T) {
	type test struct {
		input []string
		want  string
	}

	tests := []test{
		{input: []string{"eq", "bgp", "foo"}, want: "eq bgp"},
	}

	for i, tc := range tests {
		m, p, _, e := parsePort(tc.input)
		if e != nil {
			t.Errorf("failed to parse input: %s, got: %v", tc.input, e)
		}
		if fmt.Sprintf("%s %s", m, p) != tc.want {
			t.Errorf("test %d: given: %s, expected: %s, got: %s %s", i+1, tc.input, tc.want, m, p)
		}
	}
}

func TestParseAddr(t *testing.T) {
	type test struct {
		input []string
		want  string
	}

	tests := []test{
		{input: []string{"any", "foo"}, want: "any"},
		{input: []string{"host", "192.168.0.1", "foo"}, want: "host 192.168.0.1"},
		{input: []string{"192.168.0.0", "0.0.0.255", "foo"}, want: "192.168.0.0 0.0.0.255"},
		{input: []string{"10.0.0.0", "0.255.255.255", "foo"}, want: "10.0.0.0 0.255.255.255"},
	}

	for i, tc := range tests {
		got, _, err := parseAddr(tc.input)
		if fmt.Sprintf("%s", got) != tc.want {
			t.Errorf("test %d: expected: %s, got: %s, error: %v", i+1, tc.want, got, err)
		}
	}
}

// func TestNewACL(t *testing.T) {
// 	b, _ := os.ReadFile("./test_data/sample.ios")
// 	cfg := string(b)
// 	s, _ := GetACL("103", cfg, "standard")
// 	acl, err := NewACL(s)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	cfg_acl := strings.Join(s, "\n")
// 	s_acl := fmt.Sprintf("%s", acl)
// 	if s_acl != cfg_acl {
// 		t.Errorf("Expected:\n%s\nActual:\n%s", cfg_acl, s_acl)
// 	}
// }
