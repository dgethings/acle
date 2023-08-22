package ios

import (
	"testing"
)

func TestLoadEmptyFilenme(t *testing.T) {
	_, err := LoadConfig("")
	if err == nil {
		t.Errorf("Expected 'nil' file name to return an error")
	}
}

func TestLoadFilename(t *testing.T) {
	_, err := LoadConfig("../../test_data/sample.ios")
	if err != nil {
		t.Error(err)
	}
}

func TestGetAcl(t *testing.T) {
	cfg, _ := LoadConfig("../../test_data/sample.ios")
	_, err := GetACL("103", cfg)
	if err != nil {
		t.Error(err)
	}
}
