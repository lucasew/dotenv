package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEnvTerm_Args(t *testing.T) {
	env := map[string]string{}

	// Test simple key value
	err := ParseEnvTerm("--FOO=bar", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %v", env["FOO"])
	}

	// Test value with equal sign
	err = ParseEnvTerm("--API_KEY=123=456", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["API_KEY"] != "123=456" {
		t.Errorf("expected API_KEY=123=456, got %v", env["API_KEY"])
	}

	// Test invalid arg
	err = ParseEnvTerm("--INVALID", env)
	if err == nil {
		t.Errorf("expected syntax error, got nil")
	}
}

func TestParseEnvTerm_File(t *testing.T) {
	env := map[string]string{}
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	// Write dummy .env
	content := []byte("DB_HOST=localhost\nDB_PASS=secret=pass")
	if err := os.WriteFile(envFile, content, 0644); err != nil {
		t.Fatalf("failed to write tmp env file: %v", err)
	}

	err := ParseEnvTerm("@"+envFile, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", env["DB_HOST"])
	}

	if env["DB_PASS"] != "secret=pass" {
		t.Errorf("expected DB_PASS=secret=pass, got %v", env["DB_PASS"])
	}
}
