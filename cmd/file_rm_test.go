package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuimingda/subs-cli/internal/subtitles"
)

func TestFileRmCommand_DeclineDeletion(t *testing.T) {
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

	home := filepath.Join(tmpDir, "home")
	if err := os.Mkdir(home, 0o755); err != nil {
		t.Fatalf("mkdir home failed: %v", err)
	}
	t.Setenv("HOME", home)

	if err := os.WriteFile("a.srt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write a.srt failed: %v", err)
	}
	if err := os.WriteFile("b.ass", []byte("x"), 0o644); err != nil {
		t.Fatalf("write b.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"file", "rm"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "This will remove all subtitle files in current directory (srt/ass). Continue? [y/N]: ") {
		t.Fatalf("prompt not shown, output=%q", output)
	}

	if _, err := os.Stat("a.srt"); err != nil {
		t.Fatalf("a.srt should still exist, stat error: %v", err)
	}
	if _, err := os.Stat("b.ass"); err != nil {
		t.Fatalf("b.ass should still exist, stat error: %v", err)
	}
}

func TestFileRmCommand_ConfirmDeletion(t *testing.T) {
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

	home := filepath.Join(tmpDir, "home")
	if err := os.Mkdir(home, 0o755); err != nil {
		t.Fatalf("mkdir home failed: %v", err)
	}
	t.Setenv("HOME", home)

	if err := os.WriteFile("a.srt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write a.srt failed: %v", err)
	}
	if err := os.WriteFile("b.ass", []byte("x"), 0o644); err != nil {
		t.Fatalf("write b.ass failed: %v", err)
	}
	if err := os.WriteFile("c.txt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write c.txt failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetIn(strings.NewReader("y\n"))
	cmd.SetArgs([]string{"file", "rm"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	if _, err := os.Stat("a.srt"); err == nil {
		t.Fatalf("a.srt should be moved")
	}
	if _, err := os.Stat("b.ass"); err == nil {
		t.Fatalf("b.ass should be moved")
	}
	if _, err := os.Stat("c.txt"); err != nil {
		t.Fatalf("c.txt should remain: %v", err)
	}

	trashDir := filepath.Join(home, ".Trash")
	entries, err := os.ReadDir(trashDir)
	if err != nil {
		t.Fatalf("read trash dir failed: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("trash count = %d, want 2", len(entries))
	}
}

func TestFileRmCommand_NoSubtitleFiles(t *testing.T) {
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

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"file", "rm"})

	err = cmd.Execute()
	if !errors.Is(err, subtitles.ErrNoSubtitleFiles) {
		t.Fatalf("expected ErrNoSubtitleFiles, got %v", err)
	}
}
