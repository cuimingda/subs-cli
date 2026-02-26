package subtitles

import (
	"os"
	"testing"
)

func TestExtractEpisodeTag(t *testing.T) {
	got, ok := ExtractEpisodeTag("Show.S03E12.srt")
	if !ok || got != "S03E12" {
		t.Fatalf("ExtractEpisodeTag() = %q, %v; want %q, true", got, ok, "S03E12")
	}

	if _, ok := ExtractEpisodeTag("Show.s03e12.srt"); ok {
		t.Fatalf("expected lowercase tags to be not matched")
	}
}

func TestFindVideoFileByEpisodeTag(t *testing.T) {
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

	writeTestFile(t, "video.S03E12.mp4")
	writeTestFile(t, "video.S03E12.mkv")
	writeTestFile(t, "video.S03E99.mkv")

	got, err := FindVideoFileByEpisodeTag("S03E12")
	if err != nil {
		t.Fatalf("FindVideoFileByEpisodeTag() error = %v", err)
	}
	if got != "video.S03E12.mkv" {
		t.Fatalf("FindVideoFileByEpisodeTag() = %q, want %q", got, "video.S03E12.mkv")
	}
}

func TestFindVideoFileByEpisodeTag_NotFound(t *testing.T) {
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

	writeTestFile(t, "video.S03E99.mkv")
	writeTestFile(t, "video.S04E12.mp4")

	got, err := FindVideoFileByEpisodeTag("S03E12")
	if err != nil {
		t.Fatalf("FindVideoFileByEpisodeTag() error = %v", err)
	}
	if got != "" {
		t.Fatalf("FindVideoFileByEpisodeTag() = %q, want empty", got)
	}
}

func TestFindVideoFileByEpisodeTag_ExtensionPriority(t *testing.T) {
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

	writeTestFile(t, "show.S03E12.mp4")
	writeTestFile(t, "show.S03E12.mkv")

	got, err := FindVideoFileByEpisodeTag("S03E12")
	if err != nil {
		t.Fatalf("FindVideoFileByEpisodeTag() error = %v", err)
	}
	if got != "show.S03E12.mkv" {
		t.Fatalf("FindVideoFileByEpisodeTag() = %q, want %q", got, "show.S03E12.mkv")
	}
}

func writeTestFile(t *testing.T, name string) {
	t.Helper()

	if err := os.WriteFile(name, []byte("content"), 0o644); err != nil {
		t.Fatalf("write file %s failed: %v", name, err)
	}
}
