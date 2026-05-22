package registry

import (
	"strings"
	"testing"
)

// TestMergeRemoteFallsBackWhenSectionMissing reproduces the 2026-05 upstream
// schema drift where models.json dropped "iflow" entirely. The remote must be
// adopted for every section it provides while the missing section falls back
// to the embedded baseline rather than rejecting the whole catalog.
func TestMergeRemoteFallsBackWhenSectionMissing(t *testing.T) {
	base := &staticModelsJSON{
		Claude:      []*ModelInfo{{ID: "claude-base"}},
		Gemini:      []*ModelInfo{{ID: "gemini-base"}},
		Vertex:      []*ModelInfo{{ID: "vertex-base"}},
		GeminiCLI:   []*ModelInfo{{ID: "gemini-cli-base"}},
		AIStudio:    []*ModelInfo{{ID: "aistudio-base"}},
		CodexFree:   []*ModelInfo{{ID: "codex-free-base"}},
		CodexTeam:   []*ModelInfo{{ID: "codex-team-base"}},
		CodexPlus:   []*ModelInfo{{ID: "codex-plus-base"}},
		CodexPro:    []*ModelInfo{{ID: "codex-pro-base"}},
		IFlow:       []*ModelInfo{{ID: "iflow-base"}},
		Kimi:        []*ModelInfo{{ID: "kimi-base"}},
		Antigravity: []*ModelInfo{{ID: "antigravity-base"}},
	}

	remote := *base
	remote.Claude = []*ModelInfo{{ID: "claude-remote"}}
	remote.IFlow = nil

	merged, fallbacks := mergeRemoteIntoBase(base, &remote)

	if got := merged.Claude[0].ID; got != "claude-remote" {
		t.Fatalf("claude should adopt remote, got %q", got)
	}
	if got := merged.IFlow[0].ID; got != "iflow-base" {
		t.Fatalf("iflow should fall back to base, got %q", got)
	}

	if len(fallbacks) != 1 || !strings.HasPrefix(fallbacks[0], "iflow ") {
		t.Fatalf("expected exactly one iflow fallback entry, got %v", fallbacks)
	}
}

// TestMergeRejectsSectionWithDuplicateIDs ensures per-section validation kicks
// in: a malformed remote section must not poison the catalog.
func TestMergeRejectsSectionWithDuplicateIDs(t *testing.T) {
	base := &staticModelsJSON{
		Claude:      []*ModelInfo{{ID: "claude-base"}},
		Gemini:      []*ModelInfo{{ID: "gemini-base"}},
		Vertex:      []*ModelInfo{{ID: "vertex-base"}},
		GeminiCLI:   []*ModelInfo{{ID: "gemini-cli-base"}},
		AIStudio:    []*ModelInfo{{ID: "aistudio-base"}},
		CodexFree:   []*ModelInfo{{ID: "codex-free-base"}},
		CodexTeam:   []*ModelInfo{{ID: "codex-team-base"}},
		CodexPlus:   []*ModelInfo{{ID: "codex-plus-base"}},
		CodexPro:    []*ModelInfo{{ID: "codex-pro-base"}},
		IFlow:       []*ModelInfo{{ID: "iflow-base"}},
		Kimi:        []*ModelInfo{{ID: "kimi-base"}},
		Antigravity: []*ModelInfo{{ID: "antigravity-base"}},
	}

	remote := *base
	remote.Kimi = []*ModelInfo{{ID: "dup"}, {ID: "dup"}}

	merged, fallbacks := mergeRemoteIntoBase(base, &remote)

	if got := merged.Kimi[0].ID; got != "kimi-base" {
		t.Fatalf("kimi should fall back to base when remote has duplicate ids, got %q", got)
	}
	if len(fallbacks) != 1 || !strings.HasPrefix(fallbacks[0], "kimi ") {
		t.Fatalf("expected one kimi fallback entry, got %v", fallbacks)
	}
}

// TestEmbeddedCatalogValidates guards the invariant that the embedded
// baseline must always pass strict validation at process startup. If a future
// edit accidentally empties out a section in models.json, this test will fail
// loudly before users hit a runtime panic.
func TestEmbeddedCatalogValidates(t *testing.T) {
	if embeddedCatalog == nil {
		t.Fatal("embeddedCatalog must be populated by init()")
	}
	if err := validateModelsCatalog(embeddedCatalog); err != nil {
		t.Fatalf("embedded models.json must validate cleanly: %v", err)
	}
}
