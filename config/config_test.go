package config

import (
	"os"
	"testing"
)

func TestReadAPIKey(t *testing.T) {
	tempDir := t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "testfile-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString("123xyz")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	_ = os.Setenv("APIKEY_PATH", tempFile.Name())
	result, _ := ReadAPIKey()
	expected := "123xyz"
	if result != expected {
		t.Errorf("ReadAPIKey() = %s; want %s", result, expected)
	}
}

func TestGetEnvData(t *testing.T) {
	_ = os.Setenv("SYMBOL_ENV", "GOOG")
	result, _ := GetEnvData(SYMBOL_ENV)
	expected := "GOOG"
	if result != expected {
		t.Errorf("GetEnvData(SYMBOL_ENV) = %s; want %s", result, expected)
	}

	_ = os.Setenv("NDAYS_ENV", "5")
	result, _ = GetEnvData(NDAYS_ENV)
	expected = "5"
	if result != expected {
		t.Errorf("GetEnvData(NDAYS_ENV) = %s; want %s", result, expected)
	}
}
