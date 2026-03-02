package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestForceCommand_TogglesNonForcedToForced(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip force command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip force command test: test mkv not found")
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
	for _, stream := range beforeStreams {
		if stream.ID == "0:3" && stream.IsForced {
			t.Fatalf("expected stream 0:3 to be non-forced before")
		}
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"force", filepath.Base(target), "--id", "3"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	afterStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() after error = %v", err)
	}

	for _, stream := range afterStreams {
		if stream.ID == "0:3" && !stream.IsForced {
			t.Fatalf("expected stream 0:3 to be forced")
		}
	}

	if _, err := os.Stat(filepath.Base(target) + ".tmp_subs.mkv"); err == nil {
		t.Fatalf("temporary output should not remain")
	}

	if !strings.Contains(out.String(), "Toggled forced for stream 0:3") {
		t.Fatalf("output = %q, want toggle success message", out.String())
	}
}

func TestForceCommand_TogglesForcedToNonForced(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("skip force command test: ffmpeg is not available")
	}

	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip force command test: test mkv not found")
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
	cmd.SetArgs([]string{"force", filepath.Base(target), "--id", "3"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() first toggle error = %v", err)
	}

	beforeStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() before error = %v", err)
	}
	var wasForced bool
	for _, stream := range beforeStreams {
		if stream.ID == "0:3" {
			wasForced = stream.IsForced
		}
	}
	if !wasForced {
		t.Fatalf("expected stream 0:3 to be forced before second toggle")
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"force", filepath.Base(target), "--id", "3"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() second toggle error = %v", err)
	}

	afterStreams, err := getMKVStreams(filepath.Base(target))
	if err != nil {
		t.Fatalf("getMKVStreams() after error = %v", err)
	}
	for _, stream := range afterStreams {
		if stream.ID == "0:3" && stream.IsForced {
			t.Fatalf("expected stream 0:3 to be non-forced after toggle off")
		}
	}
}

func TestForceCommand_RejectsNonSubtitleStream(t *testing.T) {
	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip force command test: test mkv not found")
	}

	cmd := NewRootCmd()
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, filepath.Base(samplePath))
	if err := copyFile(samplePath, target); err != nil {
		t.Fatalf("copy target failed: %v", err)
	}

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"force", target, "--id", "1"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected non-subtitle stream error, got nil")
	}
	if !strings.Contains(err.Error(), "is not a subtitle stream") {
		t.Fatalf("error = %q, want not subtitle stream", err)
	}
}

func TestForceCommand_RejectsMissingIDFlag(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "target.mkv")
	if err := os.WriteFile(tmpFile, []byte("mock mkv"), 0o644); err != nil {
		t.Fatalf("write target.mkv failed: %v", err)
	}

	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"force", tmpFile})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected missing id flag error, got nil")
	}
	if !strings.Contains(err.Error(), "required flag(s) \"id\" not set") {
		t.Fatalf("error = %q, want required flag(s) \"id\" not set", err)
	}
}

func TestForceCommand_RejectsWhenFFmpegMissing(t *testing.T) {
	samplePath := resolveTestMkvPath(t)
	if samplePath == "" {
		t.Skip("skip force command test: test mkv not found")
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
	cmd.SetArgs([]string{"force", filepath.Base(target), "--id", "3"})
	t.Setenv("PATH", "/tmp/no-path-for-test")

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected ffmpeg missing error, got nil")
	}
	if err.Error() != "ffmpeg is not installed or not in PATH, please install ffmpeg" {
		t.Fatalf("error = %q, want ffmpeg missing message", err)
	}
}
