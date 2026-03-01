package cmd

import (
	"strings"
	"testing"
)

func TestStreamIDMatch(t *testing.T) {
	if !streamIDMatch("0:2", "2") {
		t.Fatalf("expected 0:2 match 2")
	}
	if !streamIDMatch("0:10", "10") {
		t.Fatalf("expected 0:10 match 10")
	}
	if !streamIDMatch("abc:3", "3") {
		t.Fatalf("expected abc:3 match 3")
	}
	if !streamIDMatch("0:2", "0:2") {
		t.Fatalf("expected exact match")
	}
	if streamIDMatch("0:2", "3") {
		t.Fatalf("unexpected match for 0:2 and 3")
	}
}

func TestStreamIDTail(t *testing.T) {
	if got, want := streamIDTail("0:10"), "10"; got != want {
		t.Fatalf("streamIDTail(0:10) = %q, want %q", got, want)
	}
	if got, want := streamIDTail("stream-no-colon"), "stream-no-colon"; got != want {
		t.Fatalf("streamIDTail(no colon) = %q, want %q", got, want)
	}
}

func TestSanitizeStreamTitle(t *testing.T) {
	got := sanitizeStreamTitle(`title/with:bad*chars?`)
	if strings.ContainsAny(got, `\\/:*?"<>|`) {
		t.Fatalf("sanitizeStreamTitle should remove invalid filename chars, got %q", got)
	}
	if got == "" {
		t.Fatalf("sanitizeStreamTitle should not return empty")
	}
}

func TestMKVMergeOutputPath(t *testing.T) {
	got := mkvMergeOutputPath("test.mkv")
	want := "test.mkv.tmp_subs.mkv"
	if got != want {
		t.Fatalf("mkvMergeOutputPath = %q, want %q", got, want)
	}
}

func TestFindStreamForSubtitleRemoval(t *testing.T) {
	streams := []mkvStreamInfo{
		{ID: "0:0", Type: "Video"},
		{ID: "0:1", Type: "Audio"},
		{ID: "0:2", Type: "Subtitle", Language: "eng"},
	}

	if _, err := findStreamForSubtitleRemoval(streams, "x"); err == nil {
		t.Fatalf("expected invalid id error")
	}

	subtitle, err := findStreamForSubtitleRemoval(streams, "2")
	if err != nil {
		t.Fatalf("findStreamForSubtitleRemoval() error = %v", err)
	}
	if subtitle.ID != "0:2" || subtitle.Type != "Subtitle" {
		t.Fatalf("unexpected stream %+v", subtitle)
	}

	if _, err := findStreamForSubtitleRemoval(streams, "1"); err == nil {
		t.Fatalf("expected non-subtitle error")
	}
	if _, err := findStreamForSubtitleRemoval(streams, "999"); err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestParseStreamIDAndLanguage(t *testing.T) {
	streamID, language := parseStreamIDAndLanguage("0:2(eng)")
	if streamID != "0:2" || language != "eng" {
		t.Fatalf("unexpected parse result: id=%q language=%q", streamID, language)
	}

	streamID, language = parseStreamIDAndLanguage("0:3")
	if streamID != "0:3" || language != "" {
		t.Fatalf("unexpected parse result without language: id=%q language=%q", streamID, language)
	}

	streamID, language = parseStreamIDAndLanguage("0:4(eng")
	if streamID != "0:4" || language != "" {
		t.Fatalf("unexpected parse result for incomplete language token: id=%q language=%q", streamID, language)
	}
}

func TestParseSubtitleFormat(t *testing.T) {
	if got := parseSubtitleFormat("SubRip subtitle, srt"); got != "SubRip subtitle" {
		t.Fatalf("parseSubtitleFormat() = %q, want %q", got, "SubRip subtitle")
	}

	if got := parseSubtitleFormat("ass"); got != "ass" {
		t.Fatalf("parseSubtitleFormat() without parentheses = %q, want %q", got, "ass")
	}

	if got := parseSubtitleFormat("SubRip (srt)"); got != "srt" {
		t.Fatalf("parseSubtitleFormat() with parentheses = %q, want %q", got, "srt")
	}
}

func TestMkvSubtitleOutputPathDefaultExtension(t *testing.T) {
	path, err := mkvSubtitleOutputPath("movie.mkv", "/tmp", mkvStreamInfo{
		ID:       "0:4",
		Language: "eng",
		Title:    "Title",
	})
	if err != nil {
		t.Fatalf("mkvSubtitleOutputPath() error = %v", err)
	}

	if strings.HasSuffix(path, "_4_eng_Title.srt") == false {
		t.Fatalf("path = %q, want suffix _4_eng_Title.srt", path)
	}

	pathNoInfo, err := mkvSubtitleOutputPath("movie.mkv", "/tmp", mkvStreamInfo{ID: "0:4"})
	if err != nil {
		t.Fatalf("mkvSubtitleOutputPath() error = %v", err)
	}
	if !strings.HasSuffix(pathNoInfo, "_4.srt") {
		t.Fatalf("path = %q, want suffix _4.srt when no format info", pathNoInfo)
	}
}

func TestSelectSubtitleStreams(t *testing.T) {
	streams := []mkvStreamInfo{
		{ID: "0:0", Type: "Video"},
		{ID: "0:1", Type: "Subtitle", Language: "eng"},
		{ID: "0:2", Type: "Audio"},
	}

	allStreams := streams

	selected, err := selectSubtitleStreams(streams, allStreams, "1")
	if err != nil {
		t.Fatalf("selectSubtitleStreams() error = %v", err)
	}
	if len(selected) != 1 || selected[0].ID != "0:1" || selected[0].Type != "Subtitle" {
		t.Fatalf("unexpected selected streams %+v", selected)
	}

	if _, err := selectSubtitleStreams(streams, allStreams, "2"); err == nil {
		t.Fatalf("expected non-subtitle stream error")
	}

	if _, err := selectSubtitleStreams(streams, allStreams, "9"); err == nil {
		t.Fatalf("expected stream not found error")
	}
}

func TestSanitizeFileNamePart(t *testing.T) {
	if got := sanitizeFileNamePart(""); got != "empty" {
		t.Fatalf("sanitizeFileNamePart(\"\") = %q, want %q", got, "empty")
	}

	if got := sanitizeFileNamePart(" <<>\"`"); got != "____`" {
		t.Fatalf("sanitizeFileNamePart(\" <><>\\\"`\") = %q, want %q", got, "____`")
	}
}
