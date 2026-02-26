package subtitles

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListDialogueFontsByAssFiles(t *testing.T) {
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

	if err := os.WriteFile("root.ass", []byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fn Arial}{\fnTimes New Roman}{\fnArial}`), 0o644); err != nil {
		t.Fatalf("write root.ass failed: %v", err)
	}
	if err := os.WriteFile("nested.ass", []byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fnCourier}{\fnCourier New}{\fnCalibri}`), 0o644); err != nil {
		t.Fatalf("write nested.ass failed: %v", err)
	}

	if err := os.WriteFile("other.txt", []byte("dummy"), 0o644); err != nil {
		t.Fatalf("write other.txt failed: %v", err)
	}

	got, err := ListDialogueFontsByAssFiles()
	if err != nil {
		t.Fatalf("ListDialogueFontsByAssFiles() error = %v", err)
	}

	want := []AssDialogueFonts{
		{FileName: "nested.ass", Fonts: []string{"Calibri", "Courier", "Courier New"}},
		{FileName: "root.ass", Fonts: []string{"Arial", "Times New Roman"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListDialogueFontsByAssFiles() = %+v, want %+v", got, want)
	}
}

func TestListDialogueFontsByAssFiles_WithStyleTagsAfterFont(t *testing.T) {
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

	if err := os.WriteFile("style-tags.ass", []byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fnSans Serif\fs18\b0\bord1\shad1\3C&h2F2F2F&,Hello}`), 0o644); err != nil {
		t.Fatalf("write style-tags.ass failed: %v", err)
	}

	got, err := ListDialogueFontsByAssFiles()
	if err != nil {
		t.Fatalf("ListDialogueFontsByAssFiles() error = %v", err)
	}

	want := []AssDialogueFonts{
		{FileName: "style-tags.ass", Fonts: []string{"Sans Serif"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListDialogueFontsByAssFiles() = %+v, want %+v", got, want)
	}
}

func TestListDialogueFontsByAssFiles_NoFonts(t *testing.T) {
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

	if err := os.WriteFile("emptyfont.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,no font here"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	got, err := ListDialogueFontsByAssFiles()
	if err != nil {
		t.Fatalf("ListDialogueFontsByAssFiles() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("result count = %d, want 1", len(got))
	}
	if len(got[0].Fonts) != 0 {
		t.Fatalf("fonts = %v, want empty", got[0].Fonts)
	}
}

func TestListDialogueFontsByAssFiles_NoAssFiles(t *testing.T) {
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

	if err := os.WriteFile("a.srt", []byte("just test"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	got, err := ListDialogueFontsByAssFiles()
	if err != nil {
		t.Fatalf("ListDialogueFontsByAssFiles() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("result count = %d, want 0", len(got))
	}
}

func TestListDialogueFontsByAssFiles_IgnoresSubdirectories(t *testing.T) {
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

	if err := os.WriteFile("root.ass", []byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fn Arial}`), 0o644); err != nil {
		t.Fatalf("write root.ass failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "child.ass"), []byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fn Courier}`), 0o644); err != nil {
		t.Fatalf("write child.ass failed: %v", err)
	}

	got, err := ListDialogueFontsByAssFiles()
	if err != nil {
		t.Fatalf("ListDialogueFontsByAssFiles() error = %v", err)
	}

	want := []AssDialogueFonts{{FileName: "root.ass", Fonts: []string{"Arial"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListDialogueFontsByAssFiles() = %+v, want %+v", got, want)
	}
}

func TestPruneDialogueFontTagsFromAssFiles(t *testing.T) {
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
		[]byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fnArial\fs18\b0\bord1\shad1\3C&h2F2F2F&,Hello}\N{\fnSans Serif\foobar}`),
		0o644,
	); err != nil {
		t.Fatalf("write fonts.ass failed: %v", err)
	}

	result, err := PruneDialogueFontTagsFromAssFiles()
	if err != nil {
		t.Fatalf("PruneDialogueFontTagsFromAssFiles() error = %v", err)
	}

	if result.TotalAssFiles != 1 {
		t.Fatalf("total files = %d, want 1", result.TotalAssFiles)
	}
	if result.RemovedTags != 2 {
		t.Fatalf("removed tags = %d, want 2", result.RemovedTags)
	}

	content, err := os.ReadFile("fonts.ass")
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}

	expected := `Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fs18\b0\bord1\shad1\3C&h2F2F2F&,Hello}\N{\foobar}`
	if string(content) != expected {
		t.Fatalf("content = %q, want %q", string(content), expected)
	}
}

func TestPruneDialogueFontTagsFromAssFiles_MixedFiles(t *testing.T) {
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

	if err := os.WriteFile("has-font.ass", []byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fnArial}hello`), 0o644); err != nil {
		t.Fatalf("write has-font.ass failed: %v", err)
	}
	if err := os.WriteFile("no-font.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,hello"), 0o644); err != nil {
		t.Fatalf("write no-font.ass failed: %v", err)
	}
	if err := os.WriteFile("not-subtitle.srt", []byte("hello"), 0o644); err != nil {
		t.Fatalf("write not-subtitle.srt failed: %v", err)
	}

	result, err := PruneDialogueFontTagsFromAssFiles()
	if err != nil {
		t.Fatalf("PruneDialogueFontTagsFromAssFiles() error = %v", err)
	}

	if result.TotalAssFiles != 2 {
		t.Fatalf("total files = %d, want 2", result.TotalAssFiles)
	}
	if result.RemovedTags != 1 {
		t.Fatalf("removed tags = %d, want 1", result.RemovedTags)
	}

	hasFontContent, err := os.ReadFile("has-font.ass")
	if err != nil {
		t.Fatalf("read has-font.ass failed: %v", err)
	}
	if string(hasFontContent) != `Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{}hello` {
		t.Fatalf("content = %q, want %q", string(hasFontContent), `Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{}hello`)
	}
}

func TestPruneDialogueFontTagsFromAssFiles_NoFontTags(t *testing.T) {
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

	if err := os.WriteFile("no-font.ass", []byte("Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,hello"), 0o644); err != nil {
		t.Fatalf("write no-font.ass failed: %v", err)
	}

	result, err := PruneDialogueFontTagsFromAssFiles()
	if err != nil {
		t.Fatalf("PruneDialogueFontTagsFromAssFiles() error = %v", err)
	}

	if result.TotalAssFiles != 1 {
		t.Fatalf("total files = %d, want 1", result.TotalAssFiles)
	}
	if result.RemovedTags != 0 {
		t.Fatalf("removed tags = %d, want 0", result.RemovedTags)
	}
}

func TestPruneDialogueFontTagsFromAssFiles_IgnoresSubdirectories(t *testing.T) {
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

	if err := os.WriteFile("root.ass", []byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fn Arial}hello`), 0o644); err != nil {
		t.Fatalf("write root.ass failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "child.ass"), []byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fn Courier}hello`), 0o644); err != nil {
		t.Fatalf("write child.ass failed: %v", err)
	}

	result, err := PruneDialogueFontTagsFromAssFiles()
	if err != nil {
		t.Fatalf("PruneDialogueFontTagsFromAssFiles() error = %v", err)
	}

	if result.TotalAssFiles != 1 {
		t.Fatalf("total files = %d, want 1", result.TotalAssFiles)
	}
	if result.RemovedTags != 1 {
		t.Fatalf("removed tags = %d, want 1", result.RemovedTags)
	}

	childContent, err := os.ReadFile(filepath.Join(subDir, "child.ass"))
	if err != nil {
		t.Fatalf("read child file failed: %v", err)
	}
	if string(childContent) != `Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fn Courier}hello` {
		t.Fatalf("child file changed unexpectedly: %q", string(childContent))
	}
}
