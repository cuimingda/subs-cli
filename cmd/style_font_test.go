package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStyleFontListCommand_Success(t *testing.T) {
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

	if err := os.WriteFile(
		"fonts.ass",
		[]byte("[V4+ Styles]\nFormat: Name, Fontname, Fontsize, PrimaryColour\nStyle: Default,Microsoft YaHei,22,&H00FFFFFF\nStyle: LOGO,SimHei,20,&H00FFFFFF\n"),
		0o644,
	); err != nil {
		t.Fatalf("write fonts.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"style", "font", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	want := "fonts.ass: Microsoft YaHei,SimHei"
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}
}

func TestStyleFontListCommand_UTF16File(t *testing.T) {
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

	sample, err := os.ReadFile(filepath.Join(originalDir, "..", "data", "foobar.ass"))
	if err != nil {
		t.Fatalf("read sample file failed: %v", err)
	}
	if err := os.WriteFile("foobar.ass", sample, 0o644); err != nil {
		t.Fatalf("write foobar.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"style", "font", "list"})

	execErr := cmd.Execute()
	if execErr == nil {
		t.Fatalf("expected error for non UTF-8 file, got nil")
	}
	if !strings.Contains(execErr.Error(), "Please run `subs encoding reset`") {
		t.Fatalf("error = %q, want contains `Please run `subs encoding reset`", execErr)
	}
}

func TestStyleFontListCommand_HelpOnStyleAndFont(t *testing.T) {
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

	var styleOut bytes.Buffer
	cmd.SetOut(&styleOut)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"style"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	if !strings.Contains(styleOut.String(), "Available Commands:") {
		t.Fatalf("output = %q, want to contain Available Commands", styleOut.String())
	}

	var fontOut bytes.Buffer
	cmd.SetOut(&fontOut)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"style", "font"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	if !strings.Contains(fontOut.String(), "Available Commands:") {
		t.Fatalf("output = %q, want to contain Available Commands", fontOut.String())
	}
}

func TestStyleFontListCommand_NoFonts(t *testing.T) {
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

	if err := os.WriteFile("not-found-fonts.ass", []byte("[V4+ Styles]\nFormat: Name, Fontname, Fontsize\nStyle: Default,,22\n"), 0o644); err != nil {
		t.Fatalf("write not-found-fonts.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"style", "font", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	if output != "not-found-fonts.ass: None" {
		t.Fatalf("output = %q, want not-found-fonts.ass: None", output)
	}
}

func TestStyleFontListCommand_BOMAndSpacedHeaders(t *testing.T) {
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

	if err := os.WriteFile(
		"bom.ass",
		[]byte("\ufeff[V4+ Styles]\nFormat : Name , Fontname , Fontsize\nStyle : Default ,微軟雅黑,22\nStyle : LOGO,方正黑體_GBK,20\n"),
		0o644,
	); err != nil {
		t.Fatalf("write bom.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"style", "font", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := strings.TrimSpace(out.String())
	if output != "bom.ass: 微軟雅黑,方正黑體_GBK" {
		t.Fatalf("output = %q, want bom.ass: 微軟雅黑,方正黑體_GBK", output)
	}
}

func TestStyleFontListCommand_RejectsArgs(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"style", "font", "list", "extra"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}

func TestStyleFontResetCommand_Success(t *testing.T) {
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

	if err := os.WriteFile(
		"font-reset.ass",
		[]byte("[V4+ Styles]\nFormat: Name, Fontname, Fontsize\nStyle: Default,Arial,22\nStyle: Logo,SimHei,20\n"),
		0o644,
	); err != nil {
		t.Fatalf("write font-reset.ass failed: %v", err)
	}

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"style", "font", "reset"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Reset 2 font names in 1 file(s).") {
		t.Fatalf("output = %q, want contains summary", output)
	}

	gotContent, err := os.ReadFile("font-reset.ass")
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	expected := "[V4+ Styles]\nFormat: Name, Fontname, Fontsize\nStyle: Default,Microsoft YaHei,22\nStyle: Logo,Microsoft YaHei,20\n"
	if string(gotContent) != expected {
		t.Fatalf("content = %q, want %q", string(gotContent), expected)
	}
}

func TestStyleFontResetCommand_NoArgsForParent(t *testing.T) {
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
	cmd.SetArgs([]string{"style", "font", "reset", "extra"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected args validation error, got nil")
	}
}
