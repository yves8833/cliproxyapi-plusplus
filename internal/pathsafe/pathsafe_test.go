package pathsafe

import (
	"path/filepath"
	"testing"
)

func TestSafeJoinHappyPath(t *testing.T) {
	base := t.TempDir()
	got, err := SafeJoin(base, "file.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(base, "file.json")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestSafeJoinSubdir(t *testing.T) {
	base := t.TempDir()
	got, err := SafeJoin(base, "sub/file.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != filepath.Join(base, "sub", "file.json") {
		t.Fatalf("subdir join wrong: %q", got)
	}
}

func TestSafeJoinRejectsTraversal(t *testing.T) {
	base := t.TempDir()
	cases := []string{"../etc/passwd", "sub/../../escape", "..\\windows", "/abs/path"}
	for _, c := range cases {
		if _, err := SafeJoin(base, c); err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}

func TestSafeJoinRejectsEmpty(t *testing.T) {
	if _, err := SafeJoin("", "x"); err == nil {
		t.Error("expected error for empty base")
	}
	if _, err := SafeJoin("/tmp", ""); err == nil {
		t.Error("expected error for empty input")
	}
}

func TestSafeContainAccepts(t *testing.T) {
	base := t.TempDir()
	full := filepath.Join(base, "x", "y.json")
	got, err := SafeContain(base, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestSafeContainRejectsOutside(t *testing.T) {
	base := t.TempDir()
	if _, err := SafeContain(base, "/etc/passwd"); err == nil {
		t.Error("expected error for path outside base")
	}
}

func TestSafeContainRootDirectory(t *testing.T) {
	// Test the edge case where base is the filesystem root
	// This should accept any absolute path under root without error
	if _, err := SafeContain("/", "/tmp/test.txt"); err != nil {
		t.Errorf("SafeContain with root base should accept child path: %v", err)
	}
	if _, err := SafeContain("/", "/etc/passwd"); err != nil {
		t.Errorf("SafeContain with root base should accept child path: %v", err)
	}
}

func TestSafeJoinRootDirectory(t *testing.T) {
	// Test the edge case where base is the filesystem root
	if _, err := SafeJoin("/", "tmp"); err != nil {
		t.Errorf("SafeJoin with root base should accept simple child: %v", err)
	}
}
