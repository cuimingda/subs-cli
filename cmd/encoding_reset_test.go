package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cuimingda/subs-cli/internal/subtitles"
)

func TestEncodingResetCommand_Success(t *testing.T) {
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

	if err := os.WriteFile("a.srt", []byte("hello"), 0o644); err != nil {
		t.Fatalf("write a.srt failed: %v", err)
	}
	if err := os.WriteFile("b.ass", []byte("world"), 0o644); err != nil {
		t.Fatalf("write b.ass failed: %v", err)
	}

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"encoding", "reset"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Total 2 file(s), updated 0 file(s)") {
		t.Fatalf("output should contain total count, got %q", output)
	}
	if !strings.Contains(output, "updated 0") {
		t.Fatalf("output should contain unchanged count, got %q", output)
	}
}

func TestEncodingResetCommand_NoSubtitleFiles(t *testing.T) {
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

	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"encoding", "reset"})

	err = rootCmd.Execute()
	if !errors.Is(err, subtitles.ErrNoSubtitleFiles) {
		t.Fatalf("expected ErrNoSubtitleFiles, got %v", err)
	}
}

func TestEncodingResetCommand_RejectsArgs(t *testing.T) {
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"encoding", "reset", "extra"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}
