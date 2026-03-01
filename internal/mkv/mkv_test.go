package mkv

import (
	"errors"
	"strings"
	"testing"
)

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

	if got := ParseSubtitleFormat("ass (ssa)"); got != "ass (ssa)" {
		t.Fatalf("ParseSubtitleFormat() for ass/ssa = %q, want %q", got, "ass (ssa)")
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

	path, err = SubtitleOutputPath("movie.mkv", "/tmp", StreamInfo{
		ID:             "0:5",
		Language:       "eng",
		SubtitleFormat: "ass (ssa)",
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

type fakeFFmpegRunner struct {
	installed bool
	args      []string
	output    string
	runErr    error
}

func (f *fakeFFmpegRunner) IsInstalled() error {
	if f.installed {
		return nil
	}
	return errors.New("ffmpeg not found")
}

func (f *fakeFFmpegRunner) Run(args ...string) ([]byte, error) {
	f.args = append(f.args, strings.Join(args, "|"))
	return []byte(f.output), f.runErr
}

func TestFFmpegRunnerIntegration(t *testing.T) {
	old := ffmpegRunner
	t.Cleanup(func() {
		ffmpegRunner = old
	})

	runner := &fakeFFmpegRunner{
		installed: true,
		output:    "Stream #0:0: Video: h264\nStream #0:1: Subtitle: SubRip\n",
	}
	SetFFmpegRunner(runner)
	streams, err := ListStreams("sample.mkv")
	if err != nil {
		t.Fatalf("ListStreams() error = %v", err)
	}
	if len(streams) != 2 || streams[0].ID != "0:0" || streams[1].Type != "Subtitle" {
		t.Fatalf("unexpected streams %+v", streams)
	}

	if len(runner.args) != 1 {
		t.Fatalf("runner args call count = %d, want 1", len(runner.args))
	}
}

func TestFFmpegRunnerErrorPath(t *testing.T) {
	old := ffmpegRunner
	t.Cleanup(func() { ffmpegRunner = old })

	runner := &fakeFFmpegRunner{installed: false}
	SetFFmpegRunner(runner)
	if err := RequireFFmpegInstalled(); err == nil {
		t.Fatal("expected ffmpeg installed error")
	}
}
