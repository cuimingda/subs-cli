package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cuimingda/subs-cli/internal/subtitles"
)

func TestEncodingCommand_ShowsHelp(t *testing.T) {
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"encoding"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Usage:") {
		t.Fatalf("help output should contain Usage, got %q", output)
	}
	if !strings.Contains(output, "list") {
		t.Fatalf("help output should contain list subcommand, got %q", output)
	}
}

func TestEncodingCommand_RejectsArgs(t *testing.T) {
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"encoding", "extra"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}

func TestEncodingListCommand_Success(t *testing.T) {
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
	rootCmd.SetArgs([]string{"encoding", "list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("line count = %d, want 2, output=%q", len(lines), out.String())
	}
	if !strings.HasPrefix(lines[0], "a.srt - ") {
		t.Fatalf("first line format invalid: %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "b.ass - ") {
		t.Fatalf("second line format invalid: %q", lines[1])
	}
}

func TestEncodingListCommand_RejectsArgs(t *testing.T) {
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"encoding", "list", "extra"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}

func TestEncodingListCommand_NoSubtitleFiles(t *testing.T) {
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
	rootCmd.SetArgs([]string{"encoding", "list"})

	err = rootCmd.Execute()
	if !errors.Is(err, subtitles.ErrNoSubtitleFiles) {
		t.Fatalf("expected ErrNoSubtitleFiles, got %v", err)
	}
}
