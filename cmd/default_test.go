package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultCommand_TogglesNonDefaultToDefault(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip default command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip default command test: test mkv not found")
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

	var defaultBefore bool
	for _, stream := range beforeStreams {
		if stream.ID == "0:2" && stream.IsDefault {
			defaultBefore = true
		}
	}
	if !defaultBefore {
		t.Fatalf("expected stream 0:2 to be default before")
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"default", filepath.Base(target), "--id", "3"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	afterStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() after error = %v", err)
	}

	for _, stream := range afterStreams {
		if stream.ID == "0:3" && !stream.IsDefault {
			t.Fatalf("expected stream 0:3 to be default")
		}
		if stream.ID == "0:2" && stream.IsDefault {
			t.Fatalf("expected previous default stream 0:2 to be reset")
		}
	}

	if _, err := os.Stat(filepath.Base(target) + ".tmp_subs.mkv"); err == nil {
		t.Fatalf("temporary output should not remain")
	}

	if !strings.Contains(out.String(), "Toggled default for stream 0:3") {
		t.Fatalf("output = %q, want toggle success message", out.String())
	}
}

func TestDefaultCommand_TogglesDefaultToNonDefault(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip default command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip default command test: test mkv not found")
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

	var defaultBefore bool
	for _, stream := range beforeStreams {
		if stream.ID == "0:2" && stream.IsDefault {
			defaultBefore = true
		}
	}
	if !defaultBefore {
		t.Fatalf("expected stream 0:2 to be default before")
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"default", filepath.Base(target), "--id", "2"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	afterStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() after error = %v", err)
	}

	for _, stream := range afterStreams {
		if stream.ID == "0:2" && stream.IsDefault {
			t.Fatalf("expected stream 0:2 to be non-default after toggle off")
		}
	}
}

func TestDefaultCommand_ExecutesWithoutConfirmPrompt(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip default command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip default command test: test mkv not found")
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
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"default", filepath.Base(target), "--id", "3"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	afterStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() after error = %v", err)
	}

	if len(beforeStreams) != len(afterStreams) {
		t.Fatalf("stream count changed without confirm: before=%d after=%d", len(beforeStreams), len(afterStreams))
	}
	if !strings.Contains(out.String(), "Toggled default for stream 0:3") {
		t.Fatalf("output = %q, want toggle success message", out.String())
	}
}

func TestDefaultCommand_RejectsNonSubtitleStream(t *testing.T) {
	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip default command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, filepath.Base(samplePath))
	if err := copyFile(samplePath, target); err != nil {
		t.Fatalf("copy target failed: %v", err)
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"default", target, "--id", "1"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected non-subtitle stream error, got nil")
	}
	if !strings.Contains(err.Error(), "is not a subtitle stream") {
		t.Fatalf("error = %q, want not subtitle stream", err)
	}
}

func TestDefaultCommand_RejectsMissingIDFlag(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "target.mkv")
	if err := os.WriteFile(tmpFile, []byte("mock mkv"), 0o644); err != nil {
		t.Fatalf("write target.mkv failed: %v", err)
	}

	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"default", tmpFile})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected missing id flag error, got nil")
	}
	if !strings.Contains(err.Error(), "required flag(s) \"id\" not set") {
		t.Fatalf("error = %q, want required flag(s) \"id\" not set", err)
	}
}

func TestDefaultCommand_RejectsWhenFFmpegMissing(t *testing.T) {
	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip default command test: test mkv not found")
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
	cmd.SetArgs([]string{"default", filepath.Base(target), "--id", "2"})
	t.Setenv("PATH", "/tmp/no-path-for-test")

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected ffmpeg missing error, got nil")
	}
	if err.Error() != "ffmpeg is not installed or not in PATH, please install ffmpeg" {
		t.Fatalf("error = %q, want ffmpeg missing message", err)
	}
}
