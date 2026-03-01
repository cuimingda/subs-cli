package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func NewExtractCmd() *cobra.Command {
	var streamID string
	var outputDir string

	cmd := &cobra.Command{
		Use:   "extract <mkv-file>",
		Short: "Extract all subtitle streams from an mkv file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			fileName := args[0]

			if filepath.Ext(fileName) != ".mkv" && filepath.Ext(fileName) != ".MKV" {
				return fmt.Errorf("file must be an mkv file: %s", fileName)
			}

			if _, err := os.Stat(fileName); err != nil {
				return err
			}

			if _, err := exec.LookPath("ffmpeg"); err != nil {
				return fmt.Errorf("ffmpeg is not installed or not in PATH, please install ffmpeg")
			}

			mkvOutputDir := filepath.Dir(fileName)
			if outputDir != "" {
				info, err := os.Stat(outputDir)
				if err != nil {
					return err
				}
				if !info.IsDir() {
					return fmt.Errorf("output is not a directory: %s", outputDir)
				}
				mkvOutputDir = outputDir
			}

			streams, err := getMKVStreams(fileName)
			if err != nil {
				return err
			}

			subtitleStreams := make([]mkvStreamInfo, 0, len(streams))
			for _, stream := range streams {
				if stream.Type == "Subtitle" {
					subtitleStreams = append(subtitleStreams, stream)
				}
			}

			selectedStreams, err := selectSubtitleStreams(subtitleStreams, streams, streamID)
			if err != nil {
				return err
			}

			outDir := mkvSubtitleOutputDir(fileName, mkvOutputDir)
			if _, err := os.Stat(outDir); err == nil {
				return fmt.Errorf("subtitle output directory already exists: %s", outDir)
			}
			if err := os.MkdirAll(outDir, 0o755); err != nil {
				return err
			}

			if _, err := fmt.Fprintf(cobraCmd.OutOrStdout(), "Found %d subtitle stream(s) to extract.\n", len(selectedStreams)); err != nil {
				return err
			}

			for _, stream := range selectedStreams {
				outputPath, err := mkvSubtitleOutputPath(fileName, mkvOutputDir, stream)
				if err != nil {
					return err
				}

				if _, err := fmt.Fprintf(
					cobraCmd.OutOrStdout(),
					"Exporting stream %s (lang=%s, format=%s) -> %s\n",
					stream.ID,
					displayOrEmpty(stream.Language),
					displayOrEmpty(stream.SubtitleFormat),
					outputPath,
				); err != nil {
					return err
				}

				extractCmd := exec.Command(
					"ffmpeg",
					"-hide_banner",
					"-i",
					fileName,
					"-map",
					stream.ID,
					"-c",
					"copy",
					outputPath,
				)
				extractCmd.Stderr = &bytes.Buffer{}
				if err := extractCmd.Run(); err != nil {
					return fmt.Errorf("failed to export stream %s: %w", stream.ID, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&streamID, "id", "", "Only export one subtitle stream by stream id (for example: 4)")
	cmd.Flags().StringVar(&outputDir, "output", "", "Output directory for extracted subtitle files")

	return cmd
}

func selectSubtitleStreams(subtitleStreams, allStreams []mkvStreamInfo, streamID string) ([]mkvStreamInfo, error) {
	if streamID == "" {
		return subtitleStreams, nil
	}

	targetMatched := false
	for _, stream := range allStreams {
		if streamIDMatch(stream.ID, streamID) {
			targetMatched = true
			if stream.Type == "Subtitle" {
				return []mkvStreamInfo{stream}, nil
			}
			return nil, fmt.Errorf("stream id %s is not a subtitle stream", streamID)
		}
	}

	if !targetMatched {
		return nil, fmt.Errorf("stream id %s not found", streamID)
	}

	return nil, nil
}

func streamIDMatch(rawID, target string) bool {
	rawID = strings.TrimSpace(rawID)
	target = strings.TrimSpace(target)
	if rawID == target {
		return true
	}
	return streamIDTail(rawID) == target
}

func streamIDTail(streamID string) string {
	lastColon := strings.LastIndex(streamID, ":")
	if lastColon < 0 {
		return streamID
	}
	return strings.TrimSpace(streamID[lastColon+1:])
}

func mkvSubtitleOutputDir(fileName, baseOutputDir string) string {
	base := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	return filepath.Join(baseOutputDir, base+"_subs")
}

func mkvSubtitleOutputPath(fileName, outputBaseDir string, stream mkvStreamInfo) (string, error) {
	base := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	parts := []string{base, sanitizeStreamID(stream.ID)}
	if stream.Language != "" {
		parts = append(parts, stream.Language)
	}
	if stream.Title != "" {
		parts = append(parts, sanitizeStreamTitle(stream.Title))
	}

	filename := strings.Join(parts, "_")
	ext := strings.ToLower(stream.SubtitleFormat)
	if ext == "" {
		ext = "srt"
	}

	return filepath.Join(mkvSubtitleOutputDir(fileName, outputBaseDir), fmt.Sprintf("%s.%s", filename, ext)), nil
}

func sanitizeStreamID(streamID string) string {
	safeID := strings.ReplaceAll(streamID, ":", "_")
	return sanitizeFileNamePart(safeID)
}

func sanitizeStreamTitle(title string) string {
	return sanitizeFileNamePart(title)
}

func sanitizeFileNamePart(value string) string {
	re := regexp.MustCompile(`[\\/:*?"<>|]`)
	value = strings.TrimSpace(re.ReplaceAllString(value, "_"))
	if value == "" {
		return "empty"
	}
	return value
}

func displayOrEmpty(value string) string {
	if value == "" {
		return "(EMPTY)"
	}
	return value
}
