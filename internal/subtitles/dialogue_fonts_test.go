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

	subDir := filepath.Join("sub")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join("sub", "nested.ass"), []byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fnCourier}{\fnCourier New}{\fnCalibri}`), 0o644); err != nil {
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
		{FileName: "root.ass", Fonts: []string{"Arial", "Times New Roman"}},
		{FileName: "sub/nested.ass", Fonts: []string{"Calibri", "Courier", "Courier New"}},
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

	if err := os.WriteFile("style-tags.ass", []byte(`Dialogue: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,{\fn微软雅黑\fs18\b0\bord1\shad1\3C&h2F2F2F&,Hello}`), 0o644); err != nil {
		t.Fatalf("write style-tags.ass failed: %v", err)
	}

	got, err := ListDialogueFontsByAssFiles()
	if err != nil {
		t.Fatalf("ListDialogueFontsByAssFiles() error = %v", err)
	}

	want := []AssDialogueFonts{
		{FileName: "style-tags.ass", Fonts: []string{"微软雅黑"}},
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
