package ios

import (
	"os"
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
