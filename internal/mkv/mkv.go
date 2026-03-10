package mkv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	streamLineRE       = regexp.MustCompile(`^\s*Stream #(.+?):\s*([A-Za-z]+):\s*(.+)$`)
	titleLineRE        = regexp.MustCompile(`^\s*title\s*:\s*(.+)$`)
	streamIDSplitterRE = regexp.MustCompile(`[\\/:*?"<>|]`)
	languageTagRE      = regexp.MustCompile(`^[a-z]{3}$`)
	streamDefaultRE    = regexp.MustCompile(`(?i)\bdefault\b`)
	streamForcedRE     = regexp.MustCompile(`(?i)\bforced\b`)
)

type StreamInfo struct {
	ID             string
	Type           string
	Language       string
	SubtitleFormat string
	Title          string
	IsDefault      bool
	IsForced       bool
}

type FFmpegRunner interface {
	IsInstalled() error
	Run(args ...string) ([]byte, error)
}

type commandFFmpegRunner struct{}

func (commandFFmpegRunner) IsInstalled() error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg is not installed or not in PATH, please install ffmpeg")
	}
	return nil
}

func (commandFFmpegRunner) Run(args ...string) ([]byte, error) {
	cmd := exec.Command("ffmpeg", args...)
	return cmd.CombinedOutput()
}

var ffmpegRunner FFmpegRunner = commandFFmpegRunner{}

func SetFFmpegRunner(runner FFmpegRunner) {
	if runner == nil {
		ffmpegRunner = commandFFmpegRunner{}
		return
	}
	ffmpegRunner = runner
}

func RequireFFmpegInstalled() error {
	if err := ffmpegRunner.IsInstalled(); err != nil {
		return err
	}
	return nil
}

func ParseMKVStreams(output string) ([]StreamInfo, error) {
	var streams []StreamInfo

	lines := strings.Split(output, "\n")
	var lastStream *StreamInfo
	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r")
		if strings.HasPrefix(strings.TrimSpace(line), "Stream mapping:") {
			break
		}

		if match := streamLineRE.FindStringSubmatch(line); match != nil {
			rawID := strings.TrimSpace(match[1])
			streamID, language := ParseStreamIDAndLanguage(rawID)
			streamType := strings.TrimSpace(match[2])
			streamDesc := strings.TrimSpace(match[3])

			stream := StreamInfo{
				ID:        streamID,
				Type:      streamType,
				Language:  language,
				IsDefault: isSubtitleDefault(streamDesc),
				IsForced:  isSubtitleForced(streamDesc),
			}
			if streamType == "Subtitle" {
				stream.SubtitleFormat = ParseSubtitleFormat(streamDesc)
			}

			streams = append(streams, stream)
			lastStream = &streams[len(streams)-1]
			continue
		}

		if lastStream == nil {
			continue
		}

		if titleMatch := titleLineRE.FindStringSubmatch(line); titleMatch != nil {
			lastStream.Title = strings.TrimSpace(titleMatch[1])
		}
	}

	if len(streams) == 0 {
		return nil, fmt.Errorf("no stream lines found in ffmpeg output")
	}

	return streams, nil
}

func ListStreams(fileName string) ([]StreamInfo, error) {
	output, err := RunFFmpeg("-hide_banner", "-i", fileName)
	streams, parseErr := ParseMKVStreams(string(output))
	if parseErr != nil {
		return nil, parseErr
	}
	if len(streams) == 0 {
		return nil, err
	}
	return streams, nil
}

func ParseStreamIDAndLanguage(rawID string) (streamID, language string) {
	open := strings.Index(rawID, "(")
	if open < 0 {
		return strings.TrimSpace(rawID), ""
	}

	streamID = strings.TrimSpace(rawID[:open])
	rest := rawID[open+1:]
	close := strings.Index(rest, ")")
	if close < 0 {
		return streamID, ""
	}

	return streamID, strings.TrimSpace(rest[:close])
}

func isSubtitleDefault(description string) bool {
	return streamDefaultRE.MatchString(strings.ToLower(description))
}

func isSubtitleForced(description string) bool {
	return streamForcedRE.MatchString(strings.ToLower(description))
}

func ParseSubtitleFormat(description string) string {
	description = strings.TrimSpace(description)
	if comma := strings.Index(description, ","); comma >= 0 {
		description = strings.TrimSpace(description[:comma])
	}

	open := strings.Index(description, "(")
	if open >= 0 {
		rest := description[open+1:]
		close := strings.Index(rest, ")")
		if close >= 0 {
			formatPrefix := strings.TrimSpace(description[:open])
			if strings.EqualFold(formatPrefix, "ass") {
				return strings.TrimSpace(description[:open+close+2])
			}

			formatValue := strings.TrimSpace(rest[:close])
			if isSubtitleMetadataKeyword(formatValue) {
				return formatPrefix
			}
			return formatValue
		}
	}

	return strings.TrimSpace(description)
}

func isSubtitleMetadataKeyword(keyword string) bool {
	switch strings.ToLower(strings.TrimSpace(keyword)) {
	case "default", "forced", "hearing_impaired", "visual_impaired":
		return true
	default:
		return false
	}
}

func SubtitleFileExtension(subtitleFormat string) string {
	normalizedFormat := strings.ToLower(strings.TrimSpace(subtitleFormat))
	if normalizedFormat == "" {
		return "srt"
	}

	if strings.HasPrefix(normalizedFormat, "ass") {
		return "ass"
	}

	return strings.TrimPrefix(normalizedFormat, ".")
}

func StreamIDMatch(rawID, target string) bool {
	rawID = strings.TrimSpace(rawID)
	target = strings.TrimSpace(target)
	if rawID == target {
		return true
	}
	return StreamIDTail(rawID) == target
}

func StreamIDTail(streamID string) string {
	lastColon := strings.LastIndex(streamID, ":")
	if lastColon < 0 {
		return streamID
	}
	return strings.TrimSpace(streamID[lastColon+1:])
}

func SelectSubtitleStreams(subtitleStreams, allStreams []StreamInfo, target string) ([]StreamInfo, error) {
	if target == "" {
		return subtitleStreams, nil
	}

	targetMatched := false
	for _, stream := range allStreams {
		if StreamIDMatch(stream.ID, target) {
			targetMatched = true
			if stream.Type == "Subtitle" {
				return []StreamInfo{stream}, nil
			}
			return nil, fmt.Errorf("stream id %s is not a subtitle stream", target)
		}
	}

	if !targetMatched {
		return nil, fmt.Errorf("stream id %s not found", target)
	}

	return nil, nil
}

func FindStreamForSubtitleRemoval(allStreams []StreamInfo, targetID string) (StreamInfo, error) {
	idRE := regexp.MustCompile(`^[0-9]+$`)
	if !idRE.MatchString(targetID) {
		return StreamInfo{}, fmt.Errorf("invalid stream id: %s", targetID)
	}

	for _, stream := range allStreams {
		if StreamIDMatch(stream.ID, targetID) {
			if stream.Type != "Subtitle" {
				return StreamInfo{}, fmt.Errorf("stream id %s is not a subtitle stream", targetID)
			}
			return stream, nil
		}
	}

	return StreamInfo{}, fmt.Errorf("stream id %s not found", targetID)
}

func SubtitleOutputDir(fileName, baseOutputDir string) string {
	base := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	return filepath.Join(baseOutputDir, base+"_subs")
}

func SubtitleOutputPath(fileName, outputBaseDir string, stream StreamInfo) (string, error) {
	base := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	parts := []string{base, SanitizeStreamID(stream.ID)}
	if stream.Language != "" {
		parts = append(parts, stream.Language)
	}
	if stream.Title != "" {
		parts = append(parts, SanitizeStreamTitle(stream.Title))
	}

	filename := strings.Join(parts, "_")
	ext := SubtitleFileExtension(stream.SubtitleFormat)

	return filepath.Join(SubtitleOutputDir(fileName, outputBaseDir), fmt.Sprintf("%s.%s", filename, ext)), nil
}

func SanitizeStreamID(streamID string) string {
	safeID := strings.ReplaceAll(streamID, ":", "_")
	return SanitizeFileNamePart(safeID)
}

func SanitizeStreamTitle(title string) string {
	return SanitizeFileNamePart(title)
}

func SanitizeFileNamePart(value string) string {
	value = strings.TrimSpace(streamIDSplitterRE.ReplaceAllString(value, "_"))
	if value == "" {
		return "empty"
	}
	return value
}

func MergeOutputPath(targetFile string) string {
	return targetFile + ".tmp_subs.mkv"
}

func RemoveTempOutputIfExists(outputPath string) error {
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}
	return nil
}

func ValidLanguageTag(language string) bool {
	return languageTagRE.MatchString(language)
}

func BuildMergeFFmpegArgs(targetFile, subtitleFile string, targetSubtitleCount int, languageTag, subtitleTitle string) []string {
	ffmpegArgs := []string{
		"-hide_banner",
		"-y",
		"-i",
		targetFile,
		"-i",
		subtitleFile,
		"-c",
		"copy",
		"-map",
		"0",
		"-map",
		"1",
	}
	for subtitleIndex := 0; subtitleIndex < targetSubtitleCount; subtitleIndex++ {
		ffmpegArgs = append(ffmpegArgs, "-disposition:s:"+strconv.Itoa(subtitleIndex), "-default")
	}
	newSubtitleIndex := strconv.Itoa(targetSubtitleCount)
	ffmpegArgs = append(ffmpegArgs, "-disposition:s:"+newSubtitleIndex, "default")
	if languageTag != "" || subtitleTitle != "" {
		if languageTag != "" {
			ffmpegArgs = append(ffmpegArgs, "-metadata:s:s:"+newSubtitleIndex, "language="+languageTag)
		}
		if subtitleTitle != "" {
			ffmpegArgs = append(ffmpegArgs, "-metadata:s:s:"+newSubtitleIndex, "title="+subtitleTitle)
		}
	}
	return ffmpegArgs
}

func BuildExtractFFmpegArgs(sourceFile string, stream StreamInfo, outputPath string) []string {
	return []string{
		"-hide_banner",
		"-i",
		sourceFile,
		"-map",
		stream.ID,
		"-c",
		"copy",
		outputPath,
	}
}

func BuildRemoveFFmpegArgs(sourceFile string, stream StreamInfo) []string {
	return []string{
		"-hide_banner",
		"-y",
		"-i",
		sourceFile,
		"-map",
		"0",
		"-map",
		"-" + stream.ID,
		"-c",
		"copy",
	}
}

func StreamDefaultSubtitleIndex(allStreams []StreamInfo, targetID string) (int, error) {
	subtitleIndex := 0
	for _, stream := range allStreams {
		if stream.Type != "Subtitle" {
			continue
		}
		if StreamIDMatch(stream.ID, targetID) {
			return subtitleIndex, nil
		}
		subtitleIndex++
	}
	return -1, fmt.Errorf("subtitle stream %s not found", targetID)
}

func BuildDefaultToggleFFmpegArgs(sourceFile string, allStreams []StreamInfo, target StreamInfo) []string {
	targetIndex, err := StreamDefaultSubtitleIndex(allStreams, target.ID)
	if err != nil {
		return nil
	}

	ffmpegArgs := []string{
		"-hide_banner",
		"-y",
		"-i",
		sourceFile,
		"-map",
		"0",
		"-c",
		"copy",
	}

	if target.IsDefault {
		ffmpegArgs = append(ffmpegArgs, "-disposition:s:"+strconv.Itoa(targetIndex), "0")
		return ffmpegArgs
	}

	subtitleIndex := 0
	for _, stream := range allStreams {
		if stream.Type != "Subtitle" {
			continue
		}
		value := "0"
		if stream.ID == target.ID {
			value = "default"
		}
		ffmpegArgs = append(ffmpegArgs, "-disposition:s:"+strconv.Itoa(subtitleIndex), value)
		subtitleIndex++
	}
	return ffmpegArgs
}

func BuildForceToggleFFmpegArgs(sourceFile string, allStreams []StreamInfo, target StreamInfo) []string {
	targetIndex, err := StreamDefaultSubtitleIndex(allStreams, target.ID)
	if err != nil {
		return nil
	}

	ffmpegArgs := []string{
		"-hide_banner",
		"-y",
		"-i",
		sourceFile,
		"-map",
		"0",
		"-c",
		"copy",
	}

	if target.IsForced {
		ffmpegArgs = append(ffmpegArgs, "-disposition:s:"+strconv.Itoa(targetIndex), "0")
		return ffmpegArgs
	}

	subtitleIndex := 0
	for _, stream := range allStreams {
		if stream.Type != "Subtitle" {
			continue
		}
		value := "0"
		if stream.ID == target.ID {
			value = "forced"
		}
		ffmpegArgs = append(ffmpegArgs, "-disposition:s:"+strconv.Itoa(subtitleIndex), value)
		subtitleIndex++
	}
	return ffmpegArgs
}

func RunFFmpeg(args ...string) ([]byte, error) {
	return ffmpegRunner.Run(args...)
}
