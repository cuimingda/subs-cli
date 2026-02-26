package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDialogueFontsCommand_Success(t *testing.T) {
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

	if err := os.WriteFile("root.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\\fn Arial}{\\fnTimes New Roman}{\\fnArial}"), 0o644); err != nil {
		t.Fatalf("write root.ass failed: %v", err)
	}

	subDir := filepath.Join("sub")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join("sub", "child.ass"), []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\\fnCalibri}{\\fnHelvetica}"), 0o644); err != nil {
		t.Fatalf("write child.ass failed: %v", err)
	}

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"dialogue", "fonts"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("line count = %d, want 2, output=%q", len(lines), out.String())
	}

	if lines[0] != "root.ass: Arial,Times New Roman" {
		t.Fatalf("line1 = %q, want root.ass: Arial,Times New Roman", lines[0])
	}
	if lines[1] != "sub/child.ass: Calibri,Helvetica" {
		t.Fatalf("line2 = %q, want sub/child.ass: Calibri,Helvetica", lines[1])
	}
}

func TestDialogueFontsCommand_NoFonts(t *testing.T) {
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

	if err := os.WriteFile("no-fonts.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,hello"), 0o644); err != nil {
		t.Fatalf("write no-fonts.ass failed: %v", err)
	}

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"dialogue", "fonts"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	output := strings.TrimSuffix(out.String(), "\n")
	if output != "no-fonts.ass: None" {
		t.Fatalf("output = %q, want 'no-fonts.ass: None'", out.String())
	}
}

func TestDialogueFontsCommand_NoAssFiles(t *testing.T) {
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
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"dialogue", "fonts"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	if strings.TrimSpace(out.String()) != "" {
		t.Fatalf("output should be empty, got %q", out.String())
	}
}

func TestDialogueFontsCommand_RejectsArgs(t *testing.T) {
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"dialogue", "fonts", "extra"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}
