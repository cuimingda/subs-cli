package subtitles

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListStyleFontsByAssFiles(t *testing.T) {
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

	if err := os.WriteFile("regular.ass", []byte("[Script Info]\n[V4+ Styles]\nFormat: Name, Fontname, Fontsize, PrimaryColour\nStyle: Default,Microsoft YaHei,22,&H00FFFFFF\nStyle: Logo,SimHei,20,&H00FFFFFF\nStyle: Fallback,Microsoft YaHei,20,&H00FFFFFF\n"), 0o644); err != nil {
		t.Fatalf("write regular.ass failed: %v", err)
	}
	if err := os.WriteFile("other.txt", []byte("not subtitle"), 0o644); err != nil {
		t.Fatalf("write other.txt failed: %v", err)
	}

	got, err := ListStyleFontsByAssFiles()
	if err != nil {
		t.Fatalf("ListStyleFontsByAssFiles() error = %v", err)
	}

	want := []AssStyleFonts{
		{
			FileName: "regular.ass",
			Fonts:    []string{"Microsoft YaHei", "SimHei"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListStyleFontsByAssFiles() = %+v, want %+v", got, want)
	}
}

func TestListStyleFontsByAssFiles_UsesFontnameColumnFromFormat(t *testing.T) {
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

	if err := os.WriteFile("reordered.ass", []byte("[V4+ Styles]\nFormat: Name, Fontsize, Fontname, Style, Colour, Bold\nStyle: Default,22,Microsoft YaHei,Default,&H00FFFFFF,0\nStyle: Title,20,Times New Roman,Default,&H00FFFFFF,0\n"), 0o644); err != nil {
		t.Fatalf("write reordered.ass failed: %v", err)
	}

	got, err := ListStyleFontsByAssFiles()
	if err != nil {
		t.Fatalf("ListStyleFontsByAssFiles() error = %v", err)
	}

	want := []AssStyleFonts{
		{
			FileName: "reordered.ass",
			Fonts:    []string{"Microsoft YaHei", "Times New Roman"},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListStyleFontsByAssFiles() = %+v, want %+v", got, want)
	}
}

func TestListStyleFontsByAssFiles_NoFontname(t *testing.T) {
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

	if err := os.WriteFile("unknown-format.ass", []byte("[V4+ Styles]\nFormat: Name, Fontsize, Style\nStyle: Default,22,Default\n"), 0o644); err != nil {
		t.Fatalf("write unknown-format.ass failed: %v", err)
	}

	got, err := ListStyleFontsByAssFiles()
	if err != nil {
		t.Fatalf("ListStyleFontsByAssFiles() error = %v", err)
	}

	want := []AssStyleFonts{
		{
			FileName: "unknown-format.ass",
			Fonts:    []string{},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListStyleFontsByAssFiles() = %+v, want %+v", got, want)
	}
}

func TestListStyleFontsByAssFiles_IgnoresSubdirectories(t *testing.T) {
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

	subDir := filepath.Join("sub")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	if err := os.WriteFile("root.ass", []byte("[V4+ Styles]\nFormat: Name, Fontname, Fontsize\nStyle: Root,Microsoft YaHei,22\n"), 0o644); err != nil {
		t.Fatalf("write root.ass failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "child.ass"), []byte("[V4+ Styles]\nFormat: Name, Fontname, Fontsize\nStyle: Child,Times New Roman,20\n"), 0o644); err != nil {
		t.Fatalf("write child.ass failed: %v", err)
	}

	got, err := ListStyleFontsByAssFiles()
	if err != nil {
		t.Fatalf("ListStyleFontsByAssFiles() error = %v", err)
	}

	want := []AssStyleFonts{
		{FileName: "root.ass", Fonts: []string{"Microsoft YaHei"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListStyleFontsByAssFiles() = %+v, want %+v", got, want)
	}
}

func TestResetCurrentDirAssStyleFontsToMicrosoftYaHei(t *testing.T) {
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

	result, err := ResetCurrentDirAssStyleFontsToMicrosoftYaHei()
	if err != nil {
		t.Fatalf("ResetCurrentDirAssStyleFontsToMicrosoftYaHei() error = %v", err)
	}

	if result.TotalAssFiles != 1 {
		t.Fatalf("total files = %d, want 1", result.TotalAssFiles)
	}
	if result.UpdatedFiles != 1 {
		t.Fatalf("updated files = %d, want 1", result.UpdatedFiles)
	}
	if result.UpdatedFonts != 2 {
		t.Fatalf("updated fonts = %d, want 2", result.UpdatedFonts)
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

func TestResetCurrentDirAssStyleFontsToMicrosoftYaHei_UsesFontnameColumnFromReorderedFormat(t *testing.T) {
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
		"reordered-ass.ass",
		[]byte("[V4+ Styles]\nFormat: Name, Fontsize, Fontname, Bold\nStyle: Default,22,Arial,0\nStyle: Logo,20,SimHei,0\n"),
		0o644,
	); err != nil {
		t.Fatalf("write reordered-ass.ass failed: %v", err)
	}

	result, err := ResetCurrentDirAssStyleFontsToMicrosoftYaHei()
	if err != nil {
		t.Fatalf("ResetCurrentDirAssStyleFontsToMicrosoftYaHei() error = %v", err)
	}

	if result.UpdatedFonts != 2 {
		t.Fatalf("updated fonts = %d, want 2", result.UpdatedFonts)
	}

	gotContent, err := os.ReadFile("reordered-ass.ass")
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	expected := "[V4+ Styles]\nFormat: Name, Fontsize, Fontname, Bold\nStyle: Default,22,Microsoft YaHei,0\nStyle: Logo,20,Microsoft YaHei,0\n"
	if string(gotContent) != expected {
		t.Fatalf("content = %q, want %q", string(gotContent), expected)
	}
}

func TestResetCurrentDirAssStyleFontsToMicrosoftYaHei_NoChangesWhenAlreadySet(t *testing.T) {
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
		"already.ass",
		[]byte("[V4+ Styles]\nFormat: Name, Fontname, Fontsize\nStyle: Default,Microsoft YaHei,22\n"),
		0o644,
	); err != nil {
		t.Fatalf("write already.ass failed: %v", err)
	}

	result, err := ResetCurrentDirAssStyleFontsToMicrosoftYaHei()
	if err != nil {
		t.Fatalf("ResetCurrentDirAssStyleFontsToMicrosoftYaHei() error = %v", err)
	}

	if result.UpdatedFonts != 0 {
		t.Fatalf("updated fonts = %d, want 0", result.UpdatedFonts)
	}
	if result.UpdatedFiles != 0 {
		t.Fatalf("updated files = %d, want 0", result.UpdatedFiles)
	}
}
