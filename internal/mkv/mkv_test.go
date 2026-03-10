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

	if got := ParseSubtitleFormat("SubRip (default)"); got != "SubRip" {
		t.Fatalf("ParseSubtitleFormat() with default flag = %q, want %q", got, "SubRip")
	}

	if got := ParseSubtitleFormat("SubRip (default) (forced)"); got != "SubRip" {
		t.Fatalf("ParseSubtitleFormat() with multiple flags = %q, want %q", got, "SubRip")
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

func TestParseMKVStreams_DefaultAndForcedFlags(t *testing.T) {
	output := strings.Join([]string{
		"Stream #0:0: Video: h264",
		"Stream #0:1(eng): Subtitle: SubRip (default) (forced)",
		"Stream #0:2: Audio: aac",
	}, "\n")

	streams, err := ParseMKVStreams(output)
	if err != nil {
		t.Fatalf("ParseMKVStreams() error = %v", err)
	}
	if len(streams) != 3 {
		t.Fatalf("stream count = %d, want 3", len(streams))
	}
	if !streams[1].IsDefault {
		t.Fatalf("expected subtitle stream to be default")
	}
	if !streams[1].IsForced {
		t.Fatalf("expected subtitle stream to be forced")
	}
	if streams[1].SubtitleFormat != "SubRip" {
		t.Fatalf("subtitle format = %q, want SubRip", streams[1].SubtitleFormat)
	}
}

func TestStreamDefaultSubtitleIndex(t *testing.T) {
	streams := []StreamInfo{
		{ID: "0:0", Type: "Video"},
		{ID: "0:1", Type: "Subtitle"},
		{ID: "0:2", Type: "Audio"},
		{ID: "0:3", Type: "Subtitle"},
	}
	index, err := StreamDefaultSubtitleIndex(streams, "3")
	if err != nil {
		t.Fatalf("StreamDefaultSubtitleIndex() error = %v", err)
	}
	if index != 1 {
		t.Fatalf("StreamDefaultSubtitleIndex() = %d, want 1", index)
	}
}

func TestBuildDefaultToggleFFmpegArgs(t *testing.T) {
	streams := []StreamInfo{
		{ID: "0:0", Type: "Video"},
		{ID: "0:1", Type: "Subtitle"},
		{ID: "0:2", Type: "Subtitle", IsDefault: true},
		{ID: "0:3", Type: "Audio"},
		{ID: "0:4", Type: "Subtitle"},
	}

	offTarget := StreamInfo{ID: "0:2", Type: "Subtitle", IsDefault: true}
	offArgs := BuildDefaultToggleFFmpegArgs("target.mkv", streams, offTarget)
	if len(offArgs) == 0 || !containsArgPair(offArgs, "-disposition:s:1", "0") {
		t.Fatalf("unexpected off args: %#v", offArgs)
	}

	onTarget := StreamInfo{ID: "0:1", Type: "Subtitle", IsDefault: false}
	onArgs := BuildDefaultToggleFFmpegArgs("target.mkv", streams, onTarget)
	if len(onArgs) == 0 {
		t.Fatalf("BuildDefaultToggleFFmpegArgs() for enable returned empty")
	}
	if !containsArgPair(onArgs, "-disposition:s:0", "default") {
		t.Fatalf("expected enable target to be set as default: %#v", onArgs)
	}
	if !containsArgPair(onArgs, "-disposition:s:1", "0") {
		t.Fatalf("expected other subtitle to be unset: %#v", onArgs)
	}
}

func TestBuildForceToggleFFmpegArgs(t *testing.T) {
	streams := []StreamInfo{
		{ID: "0:0", Type: "Video"},
		{ID: "0:1", Type: "Subtitle"},
		{ID: "0:2", Type: "Subtitle", IsForced: true},
		{ID: "0:3", Type: "Audio"},
		{ID: "0:4", Type: "Subtitle"},
	}

	offTarget := StreamInfo{ID: "0:2", Type: "Subtitle", IsForced: true}
	offArgs := BuildForceToggleFFmpegArgs("target.mkv", streams, offTarget)
	if len(offArgs) == 0 || !containsArgPair(offArgs, "-disposition:s:1", "0") {
		t.Fatalf("unexpected off args: %#v", offArgs)
	}

	onTarget := StreamInfo{ID: "0:1", Type: "Subtitle", IsForced: false}
	onArgs := BuildForceToggleFFmpegArgs("target.mkv", streams, onTarget)
	if len(onArgs) == 0 {
		t.Fatalf("BuildForceToggleFFmpegArgs() for enable returned empty")
	}
	if !containsArgPair(onArgs, "-disposition:s:0", "forced") {
		t.Fatalf("expected enable target to be set as forced: %#v", onArgs)
	}
	if !containsArgPair(onArgs, "-disposition:s:1", "0") {
		t.Fatalf("expected existing forced stream to be unset: %#v", onArgs)
	}
}

func containsArgPair(args []string, key string, value string) bool {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == key && args[i+1] == value {
			return true
		}
	}
	return false
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
	if !containsArgPair(args, "-disposition:s:0", "-default") {
		t.Fatalf("expected previous subtitle 0 default to be cleared: %#v", args)
	}
	if !containsArgPair(args, "-disposition:s:1", "-default") {
		t.Fatalf("expected previous subtitle 1 default to be cleared: %#v", args)
	}
	if !containsArgPair(args, "-disposition:s:2", "default") {
		t.Fatalf("expected new subtitle to be default: %#v", args)
	}
	if !containsArgPair(args, "-metadata:s:s:2", "language=eng") {
		t.Fatalf("expected new subtitle language metadata: %#v", args)
	}
	if !containsArgPair(args, "-metadata:s:s:2", "title=title") {
		t.Fatalf("expected new subtitle title metadata: %#v", args)
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
