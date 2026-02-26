package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cuimingda/subs-cli/internal/subtitles"
)

func TestListCommand_Success(t *testing.T) {
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
	if err := os.WriteFile("b.ass", []byte("x"), 0o644); err != nil {
		t.Fatalf("write b.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	wantLines := []string{"a.srt", "b.ass"}
	if len(gotLines) != len(wantLines) {
		t.Fatalf("output line count = %d, want %d, output=%q", len(gotLines), len(wantLines), out.String())
	}
	for i := range wantLines {
		if gotLines[i] != wantLines[i] {
			t.Fatalf("line %d = %q, want %q", i, gotLines[i], wantLines[i])
		}
	}
}

func TestListCommand_NoSubtitleFiles(t *testing.T) {
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
	cmd.SetArgs([]string{"list"})

	err = cmd.Execute()
	if !errors.Is(err, subtitles.ErrNoSubtitleFiles) {
		t.Fatalf("expected ErrNoSubtitleFiles, got %v", err)
	}
}

func TestListCommand_RejectsArgs(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"list", "extra"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}
