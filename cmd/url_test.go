package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestURLCmd_Help(t *testing.T) {
	cmd := NewRootCmd()

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"url"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Available Commands:") {
		t.Fatalf("output = %q, want help content", output)
	}
	if !strings.Contains(output, "yahei") {
		t.Fatalf("output = %q, want include yahei", output)
	}
}

func TestURLYaheiCommand(t *testing.T) {
	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"url", "yahei"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	const expected = "https://raw.githubusercontent.com/chengda/popular-fonts/master/%E5%BE%AE%E8%BD%AF%E9%9B%85%E9%BB%91.ttf"
	if output != expected {
		t.Fatalf("output = %q, want %q", output, expected)
	}
}

func TestURLYaheiCommand_DownloadByFlag(t *testing.T) {
	const content = "mock font bytes"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, content)
	}))
	defer server.Close()

	originalURL := yaheiFontURL
	yaheiFontURL = server.URL + "/mock-font.ttf"
	t.Cleanup(func() {
		yaheiFontURL = originalURL
	})

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
	cmd.SetArgs([]string{"url", "yahei", "--download"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	if !strings.Contains(output, "Downloading") || !strings.Contains(output, "Download complete") {
		t.Fatalf("output = %q, want start and completion messages", output)
	}

	data, err := os.ReadFile("mock-font.ttf")
	if err != nil {
		t.Fatalf("read downloaded file failed: %v", err)
	}
	if string(data) != content {
		t.Fatalf("downloaded content = %q, want %q", string(data), content)
	}
}

func TestURLYaheiCommand_DownloadByShortFlag(t *testing.T) {
	const content = "font binary"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, content)
	}))
	defer server.Close()

	originalURL := yaheiFontURL
	yaheiFontURL = server.URL + "/foo.ttf"
	t.Cleanup(func() {
		yaheiFontURL = originalURL
	})

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
	cmd.SetArgs([]string{"url", "yahei", "-d"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	data, err := os.ReadFile("foo.ttf")
	if err != nil {
		t.Fatalf("read downloaded file failed: %v", err)
	}
	if string(data) != content {
		t.Fatalf("downloaded content = %q, want %q", string(data), content)
	}
}
