package codebuddy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodeBuddyTokenStorage_SaveTokenToFile_RejectsPathInjection(t *testing.T) {
	// Test that the path safety check properly rejects path injection attempts
	// by validating against a trusted auth directory, not just the parent dir.

	storage := &CodeBuddyTokenStorage{
		AccessToken:  "test-token",
		RefreshToken: "test-refresh",
		TokenType:    "bearer",
		UserID:       "user123",
		Domain:       "example.com",
	}

	// Set a temporary auth directory for testing
	tempDir := t.TempDir()
	originalAuthDir := os.Getenv("CLIPROXY_AUTH_DIR")
	os.Setenv("CLIPROXY_AUTH_DIR", tempDir)
	defer func() {
		if originalAuthDir != "" {
			os.Setenv("CLIPROXY_AUTH_DIR", originalAuthDir)
		} else {
			os.Unsetenv("CLIPROXY_AUTH_DIR")
		}
	}()

	// Test 1: Valid path within auth directory should succeed
	validPath := filepath.Join(tempDir, "codebuddy-token.json")
	err := storage.SaveTokenToFile(validPath)
	if err != nil {
		t.Errorf("expected valid path to succeed, got error: %v", err)
	}

	// Test 2: Path with traversal should be rejected
	traversalPath := filepath.Join(tempDir, "..", "escape.json")
	err = storage.SaveTokenToFile(traversalPath)
	if err == nil {
		t.Error("expected traversal path to be rejected")
	}
	if !strings.Contains(err.Error(), "invalid token file path") {
		t.Errorf("expected path error, got: %v", err)
	}

	// Test 3: Path outside auth directory should be rejected
	outsidePath := "/tmp/outside.json"
	err = storage.SaveTokenToFile(outsidePath)
	if err == nil {
		t.Error("expected path outside auth directory to be rejected")
	}
	if !strings.Contains(err.Error(), "invalid token file path") {
		t.Errorf("expected path error, got: %v", err)
	}

	// Test 4: Verify the saved file exists at the expected location
	data, err := os.ReadFile(validPath)
	if err != nil {
		t.Errorf("expected token file to exist: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty token file")
	}
	if !strings.Contains(string(data), "test-token") {
		t.Error("expected token file to contain access token")
	}
}
