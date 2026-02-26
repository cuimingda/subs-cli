package subtitles

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListCurrentDirSubtitleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})

	writeFile(t, "a.srt")
	writeFile(t, "b.ass")
	writeFile(t, "c.txt")
	if err := os.Mkdir(filepath.Join(tmpDir, "dir.ass"), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	got, err := ListCurrentDirSubtitleFiles()
	if err != nil {
		t.Fatalf("ListCurrentDirSubtitleFiles() error = %v", err)
	}

	want := []string{"a.srt", "b.ass"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListCurrentDirSubtitleFiles() = %v, want %v", got, want)
	}
}

func TestListCurrentDirSubtitleFiles_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})

	writeFile(t, "a.txt")

	_, err = ListCurrentDirSubtitleFiles()
	if !errors.Is(err, ErrNoSubtitleFiles) {
		t.Fatalf("expected ErrNoSubtitleFiles, got %v", err)
	}
}

func writeFile(t *testing.T, name string) {
	t.Helper()

	if err := os.WriteFile(name, []byte("content"), 0o644); err != nil {
		t.Fatalf("write file %s failed: %v", name, err)
	}
}
