package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDialogueFontsCommand_Success(t *testing.T) {
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

	if err := os.WriteFile("root.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\\fn Arial}{\\fnTimes New Roman}{\\fnArial}"), 0o644); err != nil {
		t.Fatalf("write root.ass failed: %v", err)
	}
	if err := os.WriteFile("child.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\\fnCalibri}{\\fnHelvetica}"), 0o644); err != nil {
		t.Fatalf("write child.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"dialogue", "font", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("line count = %d, want 2, output=%q", len(lines), out.String())
	}

	if lines[0] != "child.ass: Calibri,Helvetica" {
		t.Fatalf("line1 = %q, want child.ass: Calibri,Helvetica", lines[0])
	}
	if lines[1] != "root.ass: Arial,Times New Roman" {
		t.Fatalf("line2 = %q, want root.ass: Arial,Times New Roman", lines[1])
	}
}

func TestDialogueFontCommand_Help(t *testing.T) {
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
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"dialogue", "font"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Available Commands:") {
		t.Fatalf("output = %q, want to contain Available Commands", output)
	}
}

func TestDialogueFontsCommand_NoFonts(t *testing.T) {
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

	if err := os.WriteFile("no-fonts.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,hello"), 0o644); err != nil {
		t.Fatalf("write no-fonts.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"dialogue", "font", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := strings.TrimSuffix(out.String(), "\n")
	if output != "no-fonts.ass: None" {
		t.Fatalf("output = %q, want 'no-fonts.ass: None'", out.String())
	}
}

func TestDialogueFontsCommand_NoAssFiles(t *testing.T) {
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
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"dialogue", "font", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	if strings.TrimSpace(out.String()) != "" {
		t.Fatalf("output should be empty, got %q", out.String())
	}
}

func TestDialogueFontsCommand_RejectsArgs(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"dialogue", "font", "list", "extra"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}

func TestDialogueFontPruneCommand_Success(t *testing.T) {
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

	if err := os.WriteFile("fonts.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\\fnArial\\fs18\\b0\\bord1\\shad1\\3C&h2F2F2F&,Hello}{\\fnTimes New Roman\\foobar}"), 0o644); err != nil {
		t.Fatalf("write fonts.ass failed: %v", err)
	}
	if err := os.WriteFile("nofont.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,hello"), 0o644); err != nil {
		t.Fatalf("write nofont.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"dialogue", "font", "prune"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Pruned 2 font tags in 2 files.") {
		t.Fatalf("output = %q, want include removal summary", output)
	}

	content, err := os.ReadFile("fonts.ass")
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	if content == nil {
		t.Fatalf("read content is nil")
	}
	expected := "Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\\fs18\\b0\\bord1\\shad1\\3C&h2F2F2F&,Hello}{\\foobar}"
	if string(content) != expected {
		t.Fatalf("content = %q, want %q", string(content), expected)
	}
}

func TestDialogueFontPruneCommand_NoFontTags(t *testing.T) {
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

	if err := os.WriteFile("nofont.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,hello"), 0o644); err != nil {
		t.Fatalf("write nofont.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"dialogue", "font", "prune"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Pruned 0 font tags in 1 files.") {
		t.Fatalf("output = %q, want include zero update summary", output)
	}
}

func TestDialogueFontPruneCommand_RejectsArgs(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"dialogue", "font", "prune", "extra"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}

func TestDialogueFontPruneCommand_RepeatedRun(t *testing.T) {
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

	if err := os.WriteFile("fonts.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\\fnArial\\fs18\\b0\\bord1\\shad1\\3C&h2F2F2F&,Hello}{\\fnTimes New Roman\\foobar}"), 0o644); err != nil {
		t.Fatalf("write fonts.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"dialogue", "font", "prune"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), "Pruned 2 font tags in 1 files.") {
		t.Fatalf("first run output = %q, want removed summary", out.String())
	}

	out.Reset()
	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), "Pruned 0 font tags in 1 files.") {
		t.Fatalf("second run output = %q, want no-tag summary", out.String())
	}
}
