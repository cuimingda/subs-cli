package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveCommand_DeletesSubtitleAfterConfirm(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip remove command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip remove command test: test mkv not found")
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

	beforeStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() before error = %v", err)
	}

	subtitleStreams := 0
	for _, stream := range beforeStreams {
		if stream.Type == "Subtitle" {
			subtitleStreams++
		}
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetIn(strings.NewReader("y\n"))
	cmd.SetArgs([]string{"remove", filepath.Base(target), "--id", "2"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Removed stream") {
		t.Fatalf("output = %q, want remove confirmation", output)
	}

	afterStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() after error = %v", err)
	}

	afterSubtitleStreams := 0
	for _, stream := range afterStreams {
		if stream.Type == "Subtitle" {
			afterSubtitleStreams++
		}
		if stream.ID == "0:2" {
			t.Fatalf("stream 0:2 should be removed")
		}
	}

	if afterSubtitleStreams != subtitleStreams-1 {
		t.Fatalf("subtitle stream count = %d, want %d", afterSubtitleStreams, subtitleStreams-1)
	}
}

func TestRemoveCommand_RequiresConfirm(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip remove command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip remove command test: test mkv not found")
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

	beforeStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() before error = %v", err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetIn(strings.NewReader("\n"))
	cmd.SetArgs([]string{"remove", filepath.Base(target), "--id", "2"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	afterStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() after error = %v", err)
	}

	if len(afterStreams) != len(beforeStreams) {
		t.Fatalf("stream count changed without confirm: before=%d after=%d", len(beforeStreams), len(afterStreams))
	}
	if !strings.Contains(errOut.String(), "This will remove stream id=2, type=Subtitle") {
		t.Fatalf("prompt missing, err=%q", errOut.String())
	}
}

func TestRemoveCommand_RejectsInvalidIDType(t *testing.T) {
	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip remove command test: test mkv not found")
	}

	target := filepath.Join(tmpDir, filepath.Base(samplePath))
	if err := copyFile(samplePath, target); err != nil {
		t.Fatalf("copy target failed: %v", err)
	}

	cmd.SetArgs([]string{"remove", target, "--id", "abc"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected invalid stream id error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid stream id") {
		t.Fatalf("error = %q, want invalid stream id", err)
	}
}

func TestRemoveCommand_RejectsNonSubtitleStream(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip remove command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip remove command test: test mkv not found")
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
	cmd.SetArgs([]string{"remove", filepath.Base(target), "--id", "1"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected non-subtitle stream error, got nil")
	}
	if !strings.Contains(err.Error(), "is not a subtitle stream") {
		t.Fatalf("error = %q, want not subtitle stream", err)
	}
}

func TestRemoveCommand_RejectsStreamNotFound(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip remove command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip remove command test: test mkv not found")
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
	cmd.SetArgs([]string{"remove", filepath.Base(target), "--id", "999"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected stream not found error, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("error = %q, want stream not found", err)
	}
}

func TestRemoveCommand_RejectsMissingIDFlag(t *testing.T) {
	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "sample.mkv")
	if err := os.WriteFile(target, []byte("x"), 0o644); err != nil {
		t.Fatalf("write sample.mkv failed: %v", err)
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"remove", target})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected missing id flag error, got nil")
	}
	if !strings.Contains(err.Error(), "required flag(s) \"id\" not set") {
		t.Fatalf("error = %q, want required flag(s) \"id\" not set", err)
	}
}
