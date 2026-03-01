package mkv

import "testing"

func TestParseStreamIDAndLanguage(t *testing.T) {
	streamID, language := ParseStreamIDAndLanguage("0:2(eng)")
	if streamID != "0:2" || language != "eng" {
		t.Fatalf("unexpected parse result: id=%q language=%q", streamID, language)
	}

	streamID, language = ParseStreamIDAndLanguage("0:3")
	if streamID != "0:3" || language != "" {
		t.Fatalf("unexpected parse result without language: id=%q language=%q", streamID, language)
	}
}

func TestParseSubtitleFormat(t *testing.T) {
	if got := ParseSubtitleFormat("SubRip subtitle, srt"); got != "SubRip subtitle" {
		t.Fatalf("ParseSubtitleFormat() = %q, want %q", got, "SubRip subtitle")
	}

	if got := ParseSubtitleFormat("SubRip (srt)"); got != "srt" {
		t.Fatalf("ParseSubtitleFormat() with parentheses = %q, want %q", got, "srt")
	}
}

func TestStreamIDMatch(t *testing.T) {
	if !StreamIDMatch("0:2", "2") {
		t.Fatal("expected 0:2 match 2")
	}
	if !StreamIDMatch("0:10", "10") {
		t.Fatal("expected 0:10 match 10")
	}
	if !StreamIDMatch("0:2", "0:2") {
		t.Fatal("expected exact match")
	}
	if StreamIDMatch("0:2", "3") {
		t.Fatal("unexpected match")
	}
}

func TestStreamIDTail(t *testing.T) {
	if got, want := StreamIDTail("0:10"), "10"; got != want {
		t.Fatalf("StreamIDTail = %q, want %q", got, want)
	}
}

func TestSelectSubtitleStreams(t *testing.T) {
	streams := []StreamInfo{
		{ID: "0:0", Type: "Video"},
		{ID: "0:1", Type: "Subtitle", Language: "eng"},
		{ID: "0:2", Type: "Audio"},
	}
	found, err := SelectSubtitleStreams(streams, streams, "1")
	if err != nil {
		t.Fatalf("SelectSubtitleStreams() error = %v", err)
	}
	if len(found) != 1 || found[0].ID != "0:1" {
		t.Fatalf("unexpected result %+v", found)
	}
}

func TestFindStreamForSubtitleRemoval(t *testing.T) {
	streams := []StreamInfo{
		{ID: "0:0", Type: "Video"},
		{ID: "0:1", Type: "Subtitle"},
	}

	stream, err := FindStreamForSubtitleRemoval(streams, "1")
	if err != nil {
		t.Fatalf("FindStreamForSubtitleRemoval() error = %v", err)
	}
	if stream.ID != "0:1" {
		t.Fatalf("unexpected stream %+v", stream)
	}

	if _, err := FindStreamForSubtitleRemoval(streams, "0"); err == nil {
		t.Fatal("expected non-subtitle stream error")
	}
	if _, err := FindStreamForSubtitleRemoval(streams, "x"); err == nil {
		t.Fatal("expected invalid stream id error")
	}
}

func TestSubtitleOutputPath(t *testing.T) {
	path, err := SubtitleOutputPath("movie.mkv", "/tmp", StreamInfo{
		ID:             "0:4",
		Language:       "eng",
		Title:          "Title",
		SubtitleFormat: "ass",
	})
	if err != nil {
		t.Fatalf("SubtitleOutputPath() error = %v", err)
	}
	if path == "" || path[len(path)-4:] != ".ass" {
		t.Fatalf("path = %q", path)
	}
}

func TestSanitizeFileNamePart(t *testing.T) {
	if got := SanitizeFileNamePart(""); got != "empty" {
		t.Fatalf("SanitizeFileNamePart() = %q, want %q", got, "empty")
	}
	if got := SanitizeFileNamePart(" <<>\"`"); got != "____`" {
		t.Fatalf("SanitizeFileNamePart() = %q", got)
	}
}

func TestBuildMergeFFmpegArgs(t *testing.T) {
	args := BuildMergeFFmpegArgs("target.mkv", "sub.srt", 2, "eng", "title")
	if len(args) == 0 {
		t.Fatal("BuildMergeFFmpegArgs() expected non-empty args")
	}
}
