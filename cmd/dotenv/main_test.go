package main

import (
	"testing"
)

func TestParseEnvTerm(t *testing.T) {
	// Reset env
	env = make(map[string]string)

	// Test --key=value
	err := parseEnvTerm("--FOO=BAR")
	if err != nil {
		t.Fatalf("parseEnvTerm returned error: %v", err)
	}
	if val, ok := env["FOO"]; !ok || val != "BAR" {
		t.Errorf("expected FOO=BAR, got %v", val)
	}

	// Test invalid format
	err = parseEnvTerm("--INVALID")
	if err == nil {
		t.Error("expected error for --INVALID, got nil")
	}
}
