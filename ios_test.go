package ios

import (
	"fmt"
	"os"
	"strings"
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

func TestNewACL(t *testing.T) {
	b, _ := os.ReadFile("./test_data/sample.ios")
	cfg := string(b)
	s, _ := GetACL("103", cfg, "standard")
	acl, err := NewACL(s)
	if err != nil {
		t.Error(err)
	}
	cfg_acl := strings.Join(s, "\n")
	s_acl := fmt.Sprintf("%s", acl)
	if s_acl != cfg_acl {
		t.Errorf("Expected:\n%s\nActual:\n%s", cfg_acl, s_acl)
	}
}
