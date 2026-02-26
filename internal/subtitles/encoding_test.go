package subtitles

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListCurrentDirSubtitleFileEncodings(t *testing.T) {
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

	writeFile(t, "a.srt")
	if err := os.WriteFile("b.ass", []byte("file-b"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	originalDetector := detectTextEncoding
	detectTextEncoding = func(content []byte) (string, error) {
		if strings.Contains(string(content), "content") {
			return "utf-8", nil
		}
		return "gb18030", nil
	}
	t.Cleanup(func() {
		detectTextEncoding = originalDetector
	})

	got, err := ListCurrentDirSubtitleFileEncodings()
	if err != nil {
		t.Fatalf("ListCurrentDirSubtitleFileEncodings() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("encoding count = %d, want 2", len(got))
	}
	if got[0].FileName != "a.srt" || got[0].Encoding != "UTF-8" {
		t.Fatalf("first entry = %+v, want {FileName:a.srt Encoding:UTF-8}", got[0])
	}
	if got[1].FileName != "b.ass" || got[1].Encoding != "GB18030" {
		t.Fatalf("second entry = %+v, want {FileName:b.ass Encoding:GB18030}", got[1])
	}
}

func TestListCurrentDirSubtitleFileEncodings_NoMatch(t *testing.T) {
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

	if err := os.WriteFile("a.txt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	_, err = ListCurrentDirSubtitleFileEncodings()
	if !errors.Is(err, ErrNoSubtitleFiles) {
		t.Fatalf("expected ErrNoSubtitleFiles, got %v", err)
	}
}

func TestListCurrentDirSubtitleFileEncodings_UnknownOnDetectorError(t *testing.T) {
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

	if err := os.WriteFile("a.srt", []byte("x"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	originalDetector := detectTextEncoding
	detectTextEncoding = func(content []byte) (string, error) {
		return "", errors.New("boom")
	}
	t.Cleanup(func() {
		detectTextEncoding = originalDetector
	})

	got, err := ListCurrentDirSubtitleFileEncodings()
	if err != nil {
		t.Fatalf("ListCurrentDirSubtitleFileEncodings() error = %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("encoding count = %d, want 1", len(got))
	}
	if got[0].Encoding != UnknownEncoding {
		t.Fatalf("encoding = %q, want %q", got[0].Encoding, UnknownEncoding)
	}
}

func TestEnsureCurrentDirAssFilesUTF8_AllowsUTF8(t *testing.T) {
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

	if err := os.WriteFile("a.ass", []byte("UTF-8 text"), 0o644); err != nil {
		t.Fatalf("write a.ass failed: %v", err)
	}

	if err := EnsureCurrentDirAssFilesUTF8([]string{"a.ass"}); err != nil {
		t.Fatalf("EnsureCurrentDirAssFilesUTF8() error = %v", err)
	}
}

func TestEnsureCurrentDirAssFilesUTF8_RejectsNonUTF8(t *testing.T) {
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

	sample, err := os.ReadFile(filepath.Join(originalDir, "..", "..", "data", "foobar.ass"))
	if err != nil {
		t.Fatalf("read sample file failed: %v", err)
	}
	if err := os.WriteFile("foobar.ass", sample, 0o644); err != nil {
		t.Fatalf("write foobar.ass failed: %v", err)
	}

	err = EnsureCurrentDirAssFilesUTF8([]string{"foobar.ass"})
	if err == nil {
		t.Fatalf("expected non-UTF-8 error, got nil")
	}
	if !strings.Contains(err.Error(), "Please run `subs encoding reset`") {
		t.Fatalf("error = %q, want contains suggestion", err)
	}
}
