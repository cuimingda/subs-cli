package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestStyleFontListCommand_Success(t *testing.T) {
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

	if err := os.WriteFile(
		"fonts.ass",
		[]byte("[V4+ Styles]\nFormat: Name, Fontname, Fontsize, PrimaryColour\nStyle: Default,Microsoft YaHei,22,&H00FFFFFF\nStyle: LOGO,SimHei,20,&H00FFFFFF\n"),
		0o644,
	); err != nil {
		t.Fatalf("write fonts.ass failed: %v", err)
	}

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"style", "font", "list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	want := "fonts.ass: Microsoft YaHei,SimHei"
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}
}

func TestStyleFontListCommand_HelpOnStyleAndFont(t *testing.T) {
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

	var styleOut bytes.Buffer
	rootCmd.SetOut(&styleOut)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"style"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	if !strings.Contains(styleOut.String(), "Available Commands:") {
		t.Fatalf("output = %q, want to contain Available Commands", styleOut.String())
	}

	var fontOut bytes.Buffer
	rootCmd.SetOut(&fontOut)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"style", "font"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	if !strings.Contains(fontOut.String(), "Available Commands:") {
		t.Fatalf("output = %q, want to contain Available Commands", fontOut.String())
	}
}

func TestStyleFontListCommand_NoFonts(t *testing.T) {
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

	if err := os.WriteFile("not-found-fonts.ass", []byte("[V4+ Styles]\nFormat: Name, Fontname, Fontsize\nStyle: Default,,22\n"), 0o644); err != nil {
		t.Fatalf("write not-found-fonts.ass failed: %v", err)
	}

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"style", "font", "list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	if output != "not-found-fonts.ass: None" {
		t.Fatalf("output = %q, want not-found-fonts.ass: None", output)
	}
}

func TestStyleFontListCommand_RejectsArgs(t *testing.T) {
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"style", "font", "list", "extra"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}
