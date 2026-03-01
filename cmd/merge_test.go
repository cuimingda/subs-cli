package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeCommand_SrtSuccess(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip merge command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip merge command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, filepath.Base(samplePath))
	subtitle := filepath.Join(tmpDir, "foobar.srt")
	if err := copyFile(samplePath, target); err != nil {
		t.Fatalf("copy target failed: %v", err)
	}
	if err := copyFile(filepath.Join(filepath.Dir(samplePath), "foobar.srt"), subtitle); err != nil {
		t.Fatalf("copy subtitle failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	beforeStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() before error = %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"merge", "--target", filepath.Base(target), "--language", "eng", "--title", "Subtitle From srt", "foobar.srt"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	afterStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() after error = %v", err)
	}

	if len(afterStreams)-len(beforeStreams) != 1 {
		t.Fatalf("stream count diff = %d, want 1", len(afterStreams)-len(beforeStreams))
	}
	if afterStreams[len(afterStreams)-1].Type != "Subtitle" {
		t.Fatalf("last stream type = %q, want Subtitle", afterStreams[len(afterStreams)-1].Type)
	}
	if afterStreams[len(afterStreams)-1].Language != "eng" {
		t.Fatalf("last stream language = %q, want eng", afterStreams[len(afterStreams)-1].Language)
	}
	if afterStreams[len(afterStreams)-1].Title != "Subtitle From srt" {
		t.Fatalf("last stream title = %q, want Subtitle From srt", afterStreams[len(afterStreams)-1].Title)
	}

	output := strings.TrimSpace(out.String())
	if !strings.Contains(output, "Found") {
		t.Fatalf("output = %q, want found streams message", output)
	}
}

func TestMergeCommand_AssSuccessWithoutOptions(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip merge command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip merge command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, filepath.Base(samplePath))
	subtitle := filepath.Join(tmpDir, "foobar.ass")
	if err := copyFile(samplePath, target); err != nil {
		t.Fatalf("copy target failed: %v", err)
	}
	if err := copyFile(filepath.Join(filepath.Dir(samplePath), "foobar.ass"), subtitle); err != nil {
		t.Fatalf("copy subtitle failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	beforeStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() before error = %v", err)
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"merge", "--target", filepath.Base(target), "foobar.ass"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	afterStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() after error = %v", err)
	}
	if len(afterStreams)-len(beforeStreams) != 1 {
		t.Fatalf("stream count diff = %d, want 1", len(afterStreams)-len(beforeStreams))
	}
	if afterStreams[len(afterStreams)-1].Type != "Subtitle" {
		t.Fatalf("last stream type = %q, want Subtitle", afterStreams[len(afterStreams)-1].Type)
	}
}

func TestMergeCommand_InvalidLanguage(t *testing.T) {
	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "target.mkv")
	subtitlePath := filepath.Join(tmpDir, "subtitle.srt")
	if err := os.WriteFile(targetPath, []byte("not real mkv"), 0o644); err != nil {
		t.Fatalf("write target.mkv failed: %v", err)
	}
	if err := os.WriteFile(subtitlePath, []byte("1\\n00:00:00,000 --> 00:00:01,000\\nhello"), 0o644); err != nil {
		t.Fatalf("write subtitle.srt failed: %v", err)
	}

	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"merge", "--target", targetPath, "--language", "en-US", subtitlePath})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected invalid language error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid language tag") {
		t.Fatalf("error = %q, want contains invalid language tag", err)
	}
}

func TestMergeCommand_InvalidTargetFileType(t *testing.T) {
	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "target.txt")
	subtitlePath := filepath.Join(tmpDir, "subtitle.srt")
	if err := os.WriteFile(subtitlePath, []byte("1\n00:00:00,000 --> 00:00:01,000\nhello\n"), 0o644); err != nil {
		t.Fatalf("write subtitle.srt failed: %v", err)
	}
	if err := os.WriteFile(targetPath, []byte("not a mkv"), 0o644); err != nil {
		t.Fatalf("write target.txt failed: %v", err)
	}

	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"merge", "--target", targetPath, subtitlePath})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected invalid target type error, got nil")
	}
	if !strings.Contains(err.Error(), "target must be an mkv file") {
		t.Fatalf("error = %q, want target must be an mkv file", err)
	}
}

func TestMergeCommand_InvalidSubtitleFileType(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := NewRootCmd()
	if err := os.WriteFile(filepath.Join(tmpDir, "foo.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write foo.txt failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "target.mkv"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write target.mkv failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"merge", "--target", "target.mkv", "foo.txt"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected unsupported subtitle format error, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported subtitle format") {
		t.Fatalf("error = %q, want unsupported subtitle format", err)
	}
}
