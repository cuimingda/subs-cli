package subtitles

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func TestResetCurrentDirSubtitleFilesToUTF8_GBK(t *testing.T) {
	testResetFileToUTF8ByEncoding(t, "GBK")
}

func TestResetCurrentDirSubtitleFilesToUTF8_GB2312(t *testing.T) {
	testResetFileToUTF8ByEncoding(t, "GB2312")
}

func TestResetCurrentDirSubtitleFilesToUTF8_GB18030(t *testing.T) {
	testResetFileToUTF8ByEncoding(t, "GB18030")
}

func TestResetCurrentDirSubtitleFilesToUTF8_GB18030WithHyphen(t *testing.T) {
	testResetFileToUTF8ByEncoding(t, "GB18030", "GB-18030")
}

func TestResetCurrentDirSubtitleFilesToUTF8_ANSI(t *testing.T) {
	testResetFileToUTF8ByEncoding(t, "ANSI")
}

func TestResetCurrentDirSubtitleFilesToUTF8_ANSIWindowsAlias(t *testing.T) {
	testResetFileToUTF8ByEncodingAndDetected(t, "ANSI", "WINDOWS-1252")
}

func TestResetCurrentDirSubtitleFilesToUTF8_SkipUTF8(t *testing.T) {
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

	original := []byte("已经是 UTF-8 文件")
	if err := os.WriteFile("a.srt", original, 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	originalDetector := detectTextEncoding
	detectTextEncoding = func(content []byte) (string, error) {
		return "UTF-8", nil
	}
	t.Cleanup(func() {
		detectTextEncoding = originalDetector
	})

	result, err := ResetCurrentDirSubtitleFilesToUTF8()
	if err != nil {
		t.Fatalf("ResetCurrentDirSubtitleFilesToUTF8() error = %v", err)
	}

	if result.Total != 1 {
		t.Fatalf("total = %d, want 1", result.Total)
	}
	if result.Updated != 0 {
		t.Fatalf("updated = %d, want 0", result.Updated)
	}

	content, err := os.ReadFile("a.srt")
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	if !bytes.Equal(content, original) {
		t.Fatalf("content changed: %q", content)
	}
	if !utf8.Valid(content) {
		t.Fatalf("content is not valid UTF-8 after reset")
	}
}

func testResetFileToUTF8ByEncoding(t *testing.T, encodeCharset string) {
	testResetFileToUTF8ByEncodingAndDetected(t, encodeCharset, encodeCharset)
}

func testResetFileToUTF8ByEncodingAndDetected(t *testing.T, encodeCharset, detectedCharset string) {
	t.Helper()

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

	text := textByCharset(encodeCharset)
	encoded, err := encodeTextByCharset(text, encodeCharset)
	if err != nil {
		t.Fatalf("encodeTextByCharset() error = %v", err)
	}
	if err := os.WriteFile("a.srt", encoded, 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	originalDetector := detectTextEncoding
	detectTextEncoding = func(content []byte) (string, error) {
		return detectedCharset, nil
	}
	t.Cleanup(func() {
		detectTextEncoding = originalDetector
	})

	result, err := ResetCurrentDirSubtitleFilesToUTF8()
	if err != nil {
		t.Fatalf("ResetCurrentDirSubtitleFilesToUTF8() error = %v", err)
	}

	if result.Total != 1 {
		t.Fatalf("total = %d, want 1", result.Total)
	}
	if result.Updated != 1 {
		t.Fatalf("updated = %d, want 1", result.Updated)
	}

	content, err := os.ReadFile("a.srt")
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	if string(content) != text {
		t.Fatalf("content = %q, want %q", string(content), text)
	}
	if !utf8.Valid(content) {
		t.Fatalf("content is not valid UTF-8")
	}
}

func encodeTextByCharset(text, charset string) ([]byte, error) {
	enc := encodingByCharset(charset)
	if enc == nil {
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}

	encoded, _, err := transform.Bytes(enc.NewEncoder(), []byte(text))
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func textByCharset(charset string) string {
	if charset == "ANSI" {
		return "caf\u00e9"
	}

	return "测试字幕"
}

func encodingByCharset(charset string) encoding.Encoding {
	switch charset {
	case "GBK":
		return simplifiedchinese.GBK
	case "GB2312":
		return simplifiedchinese.HZGB2312
	case "GB18030":
		return simplifiedchinese.GB18030
	case "ANSI":
		return charmap.Windows1252
	}

	return nil
}
