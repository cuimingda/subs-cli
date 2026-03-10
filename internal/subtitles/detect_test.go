package subtitles

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsLikelyChineseEnglishBilingual_SRT(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "bilingual.srt")
	content := "1\n00:00:01,000 --> 00:00:02,000\nHello.\n\n2\n00:00:03,000 --> 00:00:04,000\n你好。\n"
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("write bilingual.srt failed: %v", err)
	}

	if !IsLikelyChineseEnglishBilingual(file) {
		t.Fatalf("expected %s to be detected as bilingual", file)
	}
}

func TestIsLikelyChineseEnglishBilingual_ASS(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "bilingual.ass")
	content := "[Script Info]\nTitle: Test\n\n[Events]\nDialogue: 0,0:00:01.00,0:00:02.00,Default,,0,0,0,,你好\\NHello\n"
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("write bilingual.ass failed: %v", err)
	}

	if !IsLikelyChineseEnglishBilingual(file) {
		t.Fatalf("expected %s to be detected as bilingual", file)
	}
}

func TestIsLikelyChineseEnglishBilingual_ASSHeaderDoesNotTriggerFalsePositive(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "chinese-only.ass")
	content := "[Script Info]\nTitle: Chinese Only\nScriptType: v4.00+\n\n[Events]\nDialogue: 0,0:00:01.00,0:00:02.00,Default,,0,0,0,,你好\n"
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("write chinese-only.ass failed: %v", err)
	}

	if IsLikelyChineseEnglishBilingual(file) {
		t.Fatalf("expected %s to be detected as non-bilingual", file)
	}
}
