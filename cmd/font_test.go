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

func TestFontCmd_Help(t *testing.T) {
	cmd := NewRootCmd()

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"font"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Available Commands:") {
		t.Fatalf("output = %q, want help content", output)
	}
	if !strings.Contains(output, "list") {
		t.Fatalf("output = %q, want include list", output)
	}
	if !strings.Contains(output, "url") {
		t.Fatalf("output = %q, want include url", output)
	}
	if !strings.Contains(output, "download") {
		t.Fatalf("output = %q, want include download", output)
	}
}

func TestFontListCommand(t *testing.T) {
	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"font", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	want := "yahei: " + fontURLs["yahei"]
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}
}

func TestFontURLCommand(t *testing.T) {
	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"font", "url", "yahei"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	if output != fontURLs["yahei"] {
		t.Fatalf("output = %q, want %q", output, fontURLs["yahei"])
	}
}

func TestFontURLCommand_UnsupportedFont(t *testing.T) {
	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"font", "url", "not-exist"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected error for unsupported font")
	}
	if !strings.Contains(err.Error(), "unsupported font") {
		t.Fatalf("error = %q, want contains unsupported font", err)
	}
}

func TestFontDownloadCommand(t *testing.T) {
	const content = "mock font bytes"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, content)
	}))
	defer server.Close()

	originFontURLs := make(map[string]string, len(fontURLs))
	for name, url := range fontURLs {
		originFontURLs[name] = url
	}
	t.Cleanup(func() {
		fontURLs = originFontURLs
	})
	fontURLs = map[string]string{"yahei": server.URL + "/mock-font.ttf"}

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
	cmd.SetArgs([]string{"font", "download", "yahei"})

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

func TestFontDownloadCommand_UnsupportedFont(t *testing.T) {
	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"font", "download", "not-exist"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected error for unsupported font")
	}
	if !strings.Contains(err.Error(), "unsupported font") {
		t.Fatalf("error = %q, want contains unsupported font", err)
	}
}
