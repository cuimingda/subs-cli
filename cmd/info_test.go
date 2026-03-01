package cmd

import (
	"bytes"
	"path/filepath"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestInfoCommand_Success(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip info command test: ffmpeg is not available")
	}

	candidates := []string{
		filepath.Join("resources", "low_quality_with_subtitles_5s.mkv"),
		filepath.Join("..", "resources", "low_quality_with_subtitles_5s.mkv"),
		filepath.Join("cmd", "resources", "low_quality_with_subtitles_5s.mkv"),
	}

	mkvPath := ""
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			mkvPath = candidate
			break
		}
	}
	if mkvPath == "" {
		t.Skip("skip info command test: mkv sample file not found")
	}

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"info", mkvPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 11 {
		t.Fatalf("line count = %d, want 11, output=%q", len(lines), out.String())
	}

	want := []string{
		"0:0 Video (EMPTY)",
		"0:1 Audio (EMPTY)",
		"0:2 Subtitle eng srt (EMPTY)",
		"0:3 Subtitle hun srt (EMPTY)",
		"0:4 Subtitle ger srt (EMPTY)",
		"0:5 Subtitle fre srt (EMPTY)",
		"0:6 Subtitle spa srt (EMPTY)",
		"0:7 Subtitle ita srt (EMPTY)",
		"0:8 Audio Commentary",
		"0:9 Subtitle jpn srt (EMPTY)",
		"0:10 Subtitle (EMPTY) srt (EMPTY)",
	}
	for i, expected := range want {
		if lines[i] != expected {
			t.Fatalf("line %d = %q, want %q", i, lines[i], expected)
		}
	}
}

func TestInfoCommand_RejectsNonMkvFile(t *testing.T) {
	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	if err := os.WriteFile("sample.txt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write sample.txt failed: %v", err)
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"info", "sample.txt"})

	err = cmd.Execute()
	if err == nil {
		t.Fatalf("expected non-mkv validation error, got nil")
	}
	if err.Error() != "file must be an mkv file: sample.txt" {
		t.Fatalf("error = %q, want file must be an mkv file: sample.txt", err)
	}
}

func TestInfoCommand_RejectsNoSuchFile(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"info", "does-not-exist.mkv"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected file missing error, got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected not exist error, got %v", err)
	}
}

func TestInfoCommand_RejectsArgs(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"info", "a.mkv", "extra"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}

func TestInfoCommand_RejectsWhenFFmpegMissing(t *testing.T) {
	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip info command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, filepath.Base(samplePath))
	if err := copyFile(samplePath, target); err != nil {
		t.Fatalf("copy target failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"info", filepath.Base(target)})
	t.Setenv("PATH", "/tmp/no-path-for-test")

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected ffmpeg missing error, got nil")
	}
	if err.Error() != "ffmpeg is not installed or not in PATH, please install ffmpeg" {
		t.Fatalf("error = %q, want ffmpeg missing message", err)
	}
}

func TestInfoCommand_RejectsInvalidMKVContent(t *testing.T) {
	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "bad.mkv")
	if err := os.WriteFile(target, []byte("not mkv"), 0o644); err != nil {
		t.Fatalf("write bad mkv failed: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"info", filepath.Base(target)})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected info parse error, got nil")
	}

	if !strings.Contains(err.Error(), "no stream lines found") && !strings.Contains(err.Error(), "failed") {
		t.Fatalf("error = %q, want parse failure", err)
	}
}
