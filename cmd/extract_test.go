package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractCommand_Success_AllSubtitles(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip extract command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip extract command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	if err := copyFile(samplePath, filepath.Join(tmpDir, filepath.Base(samplePath))); err != nil {
		t.Fatalf("copy sample failed: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	mkvName := filepath.Base(samplePath)
	outDir := strings.TrimSuffix(mkvName, filepath.Ext(mkvName)) + "_subs"

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"extract", mkvName})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	if !strings.Contains(output, "Found 8 subtitle stream(s)") {
		t.Fatalf("output = %q, want contains found subtitle count", output)
	}

	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatalf("read output dir failed: %v", err)
	}
	if len(entries) != 8 {
		t.Fatalf("output subtitle file count = %d, want 8", len(entries))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			t.Fatalf("expected no subdirectory in output dir, got %q", entry.Name())
		}
		if !strings.HasPrefix(entry.Name(), strings.TrimSuffix(mkvName, filepath.Ext(mkvName))+"_") {
			t.Fatalf("unexpected filename: %q", entry.Name())
		}
		if !(strings.HasSuffix(entry.Name(), ".srt") || strings.HasSuffix(entry.Name(), ".ass")) {
			t.Fatalf("expected subtitle extension, got %q", entry.Name())
		}
	}
}

func TestExtractCommand_SingleID(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip extract command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip extract command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	if err := copyFile(samplePath, filepath.Join(tmpDir, filepath.Base(samplePath))); err != nil {
		t.Fatalf("copy sample failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	mkvName := filepath.Base(samplePath)
	outDir := strings.TrimSuffix(mkvName, filepath.Ext(mkvName)) + "_subs"

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"extract", "--id", "4", mkvName})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	if !strings.Contains(output, "Found 1 subtitle stream(s)") {
		t.Fatalf("output = %q, want single subtitle count", output)
	}

	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatalf("read output dir failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("output subtitle file count = %d, want 1", len(entries))
	}
}

func TestExtractCommand_Success_WithOutputDir(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip extract command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip extract command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	if err := copyFile(samplePath, filepath.Join(tmpDir, filepath.Base(samplePath))); err != nil {
		t.Fatalf("copy sample failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	outputRoot := filepath.Join(tmpDir, "custom_output")
	if err := os.MkdirAll(outputRoot, 0o755); err != nil {
		t.Fatalf("mkdir output root failed: %v", err)
	}

	mkvName := filepath.Base(samplePath)
	outDir := filepath.Join(outputRoot, strings.TrimSuffix(mkvName, filepath.Ext(mkvName))+"_subs")

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"extract", "--output", outputRoot, mkvName})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatalf("read output dir failed: %v", err)
	}
	if len(entries) != 8 {
		t.Fatalf("output subtitle file count = %d, want 8", len(entries))
	}
}

func TestExtractCommand_InvalidOutputDir(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip extract command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip extract command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	if err := copyFile(samplePath, filepath.Join(tmpDir, filepath.Base(samplePath))); err != nil {
		t.Fatalf("copy sample failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "not_a_directory.txt")
	if err := os.WriteFile(outputFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("write output file failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"extract", "--output", outputFile, filepath.Base(samplePath)})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected invalid output directory error, got nil")
	}
	if !strings.Contains(err.Error(), "output is not a directory") {
		t.Fatalf("error = %q, want output is not a directory", err)
	}
}

func TestExtractCommand_OutputSubDirExists(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip extract command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip extract command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	if err := copyFile(samplePath, filepath.Join(tmpDir, filepath.Base(samplePath))); err != nil {
		t.Fatalf("copy sample failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	mkvName := filepath.Base(samplePath)
	outDir := strings.TrimSuffix(mkvName, filepath.Ext(mkvName)) + "_subs"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("mkdir output subdir failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"extract", mkvName})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected output subdir exists error, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("error = %q, want contains 'already exists'", err)
	}
}

func TestExtractCommand_RejectedNonSubtitleID(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip extract command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip extract command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	if err := copyFile(samplePath, filepath.Join(tmpDir, filepath.Base(samplePath))); err != nil {
		t.Fatalf("copy sample failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"extract", "--id", "1", filepath.Base(samplePath)})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected non-subtitle stream error, got nil")
	}
	if !strings.Contains(err.Error(), "not a subtitle stream") {
		t.Fatalf("error = %q, want contains 'not a subtitle stream'", err)
	}
}

func TestExtractCommand_RejectsNonMkvFile(t *testing.T) {
	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	if err := os.WriteFile("sample.txt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write sample.txt failed: %v", err)
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"extract", "sample.txt"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected non-mkv validation error, got nil")
	}
	if err.Error() != "file must be an mkv file: sample.txt" {
		t.Fatalf("error = %q, want file must be an mkv file: sample.txt", err)
	}
}

func TestExtractCommand_RejectsArgs(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"extract", "a.mkv", "extra"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}

func TestExtractCommand_RejectsNoSuchID(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip extract command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip extract command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	if err := copyFile(samplePath, filepath.Join(tmpDir, filepath.Base(samplePath))); err != nil {
		t.Fatalf("copy sample failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"extract", "--id", "999", filepath.Base(samplePath)})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected stream not found error, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("error = %q, want contains stream not found", err)
	}
}

func resolveTestMkvPath(t *testing.T) string {
	t.Helper()

	candidates := []string{
		filepath.Join("resources", "low_quality_with_subtitles_5s.mkv"),
		filepath.Join("..", "resources", "low_quality_with_subtitles_5s.mkv"),
		filepath.Join("cmd", "resources", "low_quality_with_subtitles_5s.mkv"),
		filepath.Join("..", "cmd", "resources", "low_quality_with_subtitles_5s.mkv"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
