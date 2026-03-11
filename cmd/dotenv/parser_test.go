package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEnvTerm(t *testing.T) {
	t.Run("parse valid --key=value", func(t *testing.T) {
		env := map[string]string{}
		err := ParseEnvTerm(env, "--FOO=bar")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if env["FOO"] != "bar" {
			t.Errorf("expected env[\"FOO\"] to be \"bar\", got %q", env["FOO"])
		}
	})

	t.Run("parse invalid --key=value", func(t *testing.T) {
		env := map[string]string{}
		err := ParseEnvTerm(env, "--FOO")
		if err == nil {
			t.Fatal("expected error for malformed input, got nil")
		}
	})

	t.Run("parse valid @file.txt", func(t *testing.T) {
		env := map[string]string{}

		tempDir := t.TempDir()
		envFile := filepath.Join(tempDir, ".env.test")
		err := os.WriteFile(envFile, []byte("TEST_VAR=123\nANOTHER_VAR=456"), 0644)
		if err != nil {
			t.Fatalf("failed to create temp env file: %v", err)
		}

		err = ParseEnvTerm(env, "@"+envFile)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if env["TEST_VAR"] != "123" {
			t.Errorf("expected env[\"TEST_VAR\"] to be \"123\", got %q", env["TEST_VAR"])
		}
		if env["ANOTHER_VAR"] != "456" {
			t.Errorf("expected env[\"ANOTHER_VAR\"] to be \"456\", got %q", env["ANOTHER_VAR"])
		}
	})

	t.Run("parse missing @file.txt", func(t *testing.T) {
		env := map[string]string{}
		err := ParseEnvTerm(env, "@/path/to/nonexistent/file.env")
		if err == nil {
			t.Fatal("expected error for nonexistent file, got nil")
		}
	})

	t.Run("ignore arbitrary arguments", func(t *testing.T) {
		env := map[string]string{}
		err := ParseEnvTerm(env, "some_random_arg")
		if err != nil {
			t.Fatalf("expected no error for ignored arg, got %v", err)
		}
		if len(env) != 0 {
			t.Errorf("expected env to be empty, got size %d", len(env))
		}
	})
}
