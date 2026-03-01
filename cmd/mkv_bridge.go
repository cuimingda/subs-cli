package cmd

import "github.com/cuimingda/subs-cli/internal/mkv"

type mkvStreamInfo = mkv.StreamInfo

func getMKVStreams(fileName string) ([]mkvStreamInfo, error) {
	return mkv.ListStreams(fileName)
}

func parseMKVStreams(output string) ([]mkvStreamInfo, error) {
	return mkv.ParseMKVStreams(output)
}

func parseStreamIDAndLanguage(rawID string) (string, string) {
	return mkv.ParseStreamIDAndLanguage(rawID)
}

func parseSubtitleFormat(description string) string {
	return mkv.ParseSubtitleFormat(description)
}

func streamIDMatch(rawID, target string) bool {
	return mkv.StreamIDMatch(rawID, target)
}

func streamIDTail(streamID string) string {
	return mkv.StreamIDTail(streamID)
}

func selectSubtitleStreams(subtitleStreams, allStreams []mkvStreamInfo, streamID string) ([]mkvStreamInfo, error) {
	return mkv.SelectSubtitleStreams(subtitleStreams, allStreams, streamID)
}

func mkvSubtitleOutputDir(fileName, baseOutputDir string) string {
	return mkv.SubtitleOutputDir(fileName, baseOutputDir)
}

func mkvSubtitleOutputPath(fileName, outputBaseDir string, stream mkvStreamInfo) (string, error) {
	return mkv.SubtitleOutputPath(fileName, outputBaseDir, stream)
}

func sanitizeStreamID(streamID string) string {
	return mkv.SanitizeStreamID(streamID)
}

func sanitizeStreamTitle(title string) string {
	return mkv.SanitizeStreamTitle(title)
}

func sanitizeFileNamePart(value string) string {
	return mkv.SanitizeFileNamePart(value)
}

func displayOrEmpty(value string) string {
	if value == "" {
		return "(EMPTY)"
	}
	return value
}

func mkvMergeOutputPath(targetFile string) string {
	return mkv.MergeOutputPath(targetFile)
}

func validLanguageTag(language string) bool {
	return mkv.ValidLanguageTag(language)
}

func findStreamForSubtitleRemoval(allStreams []mkvStreamInfo, targetID string) (mkvStreamInfo, error) {
	return mkv.FindStreamForSubtitleRemoval(allStreams, targetID)
}
