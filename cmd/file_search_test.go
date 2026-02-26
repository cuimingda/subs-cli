package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cuimingda/subs-cli/internal/subtitles"
)

func TestFileSearchCommand_Success(t *testing.T) {
	cmd := NewRootCmd()
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

	if err := os.WriteFile("a.srt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write a.srt failed: %v", err)
	}
	if err := os.WriteFile("found.S03E10.ass", []byte("x"), 0o644); err != nil {
		t.Fatalf("write found.S03E10.ass failed: %v", err)
	}
	if err := os.WriteFile("same.S03E12.ass", []byte("x"), 0o644); err != nil {
		t.Fatalf("write same.S03E12.ass failed: %v", err)
	}
	if err := os.WriteFile("c.S03E99.srt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write c.S03E99.srt failed: %v", err)
	}
	if err := os.WriteFile("other_show.S03E10.mkv", []byte("x"), 0o644); err != nil {
		t.Fatalf("write other_show.S03E10.mkv failed: %v", err)
	}
	if err := os.WriteFile("same.S03E12.mkv", []byte("x"), 0o644); err != nil {
		t.Fatalf("write same.S03E12.mkv failed: %v", err)
	}
	if err := os.WriteFile("other_show.S03E10.mp4", []byte("x"), 0o644); err != nil {
		t.Fatalf("write other_show.S03E10.mp4 failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"file", "search"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	ignoreText := "\x1b[31mignore\x1b[0m"
	notFoundText := "\x1b[31mnot found\x1b[0m"
	sameText := "\x1b[32m(same)\x1b[0m"
	wantLines := []string{
		"a.srt => " + ignoreText,
		"c.S03E99.srt => " + notFoundText,
		"found.S03E10.ass => other_show.S03E10.mkv (found)",
		"same.S03E12.ass => same.S03E12.mkv " + sameText,
	}
	if len(gotLines) != len(wantLines) {
		t.Fatalf("output line count = %d, want %d, output=%q", len(gotLines), len(wantLines), out.String())
	}
	for i := range wantLines {
		if gotLines[i] != wantLines[i] {
			t.Fatalf("line %d = %q, want %q", i, gotLines[i], wantLines[i])
		}
	}
}

func TestFileSearchCommand_NoSubtitleFiles(t *testing.T) {
	cmd := NewRootCmd()
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

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"file", "search"})

	err = cmd.Execute()
	if !errors.Is(err, subtitles.ErrNoSubtitleFiles) {
		t.Fatalf("expected ErrNoSubtitleFiles, got %v", err)
	}
}
