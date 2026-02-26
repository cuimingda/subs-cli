package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cuimingda/subs-cli/internal/subtitles"
)

func TestFileRenameCommand_RenameAndStatus(t *testing.T) {
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

	if err := os.WriteFile("ignore.srt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write ignore.srt failed: %v", err)
	}
	if err := os.WriteFile("missing.S03E99.srt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write missing.S03E99.srt failed: %v", err)
	}
	if err := os.WriteFile("same.S03E12.srt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write same.S03E12.srt failed: %v", err)
	}
	if err := os.WriteFile("same.S03E12.mkv", []byte("x"), 0o644); err != nil {
		t.Fatalf("write same.S03E12.mkv failed: %v", err)
	}
	if err := os.WriteFile("rename abcd S02E04 abc.srt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write rename abcd S02E04 abc.srt failed: %v", err)
	}
	if err := os.WriteFile("abcd rename S02E04.mp4", []byte("x"), 0o644); err != nil {
		t.Fatalf("write abcd rename S02E04.mp4 failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"file", "rename"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	ignoreText := "\x1b[31mignore\x1b[0m"
	notFoundText := "\x1b[31mnot found\x1b[0m"
	sameText := "\x1b[32m(same)\x1b[0m"

	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	wantLines := []string{
		"ignore.srt => " + ignoreText,
		"missing.S03E99.srt => " + notFoundText,
		"rename abcd S02E04 abc.srt => abcd rename S02E04.srt (renamed)",
		"same.S03E12.srt => same.S03E12.mkv " + sameText,
	}
	if len(gotLines) != len(wantLines) {
		t.Fatalf("output line count = %d, want %d, output=%q", len(gotLines), len(wantLines), out.String())
	}
	for i := range wantLines {
		if gotLines[i] != wantLines[i] {
			t.Fatalf("line %d = %q, want %q", i, gotLines[i], wantLines[i])
		}
	}

	if _, err := os.Stat("rename abcd S02E04 abc.srt"); err == nil {
		t.Fatalf("old subtitle file still exists after rename")
	}
	if _, err := os.Stat("abcd rename S02E04.srt"); err != nil {
		t.Fatalf("expected renamed subtitle file missing: %v", err)
	}
}

func TestFileRenameCommand_NoSubtitleFiles(t *testing.T) {
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
	cmd.SetArgs([]string{"file", "rename"})

	err = cmd.Execute()
	if !errors.Is(err, subtitles.ErrNoSubtitleFiles) {
		t.Fatalf("expected ErrNoSubtitleFiles, got %v", err)
	}
}
