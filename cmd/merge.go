package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuimingda/subs-cli/internal/mkv"
	"github.com/cuimingda/subs-cli/internal/subtitles"

	"github.com/spf13/cobra"
)

func NewMergeCmd() *cobra.Command {
	var targetFile string
	var languageTag string
	var subtitleTitle string

	cmd := &cobra.Command{
		Use:   "merge <subtitle_filename>",
		Short: "Merge a subtitle file into mkv as a new stream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			subtitleFile := args[0]
			subtitleExt := strings.ToLower(filepath.Ext(subtitleFile))

			if subtitleExt != ".srt" && subtitleExt != ".ass" && subtitleExt != ".ssa" {
				return fmt.Errorf("unsupported subtitle format: %s", subtitleFile)
			}

			if _, err := os.Stat(subtitleFile); err != nil {
				return err
			}

			if filepath.Ext(targetFile) != ".mkv" && filepath.Ext(targetFile) != ".MKV" {
				return fmt.Errorf("target must be an mkv file: %s", targetFile)
			}

			if _, err := os.Stat(targetFile); err != nil {
				return err
			}

			if subtitles.IsLikelyChineseEnglishBilingual(subtitleFile) {
				if subtitleTitle == "" {
					subtitleTitle = "Chinese-English"
				}
				if languageTag == "" {
					languageTag = "eng"
				}
			}

			if languageTag != "" && !validLanguageTag(languageTag) {
				return fmt.Errorf("invalid language tag: %s", languageTag)
			}

			if err := mkv.RequireFFmpegInstalled(); err != nil {
				return err
			}

			streams, err := getMKVStreams(targetFile)
			if err != nil {
				return err
			}

			if _, err := cmd.OutOrStdout().Write([]byte(fmt.Sprintf("Found %d existing streams in target.\n", len(streams)))); err != nil {
				return err
			}

			targetSubtitleCount := 0
			for _, stream := range streams {
				if stream.Type == "Subtitle" {
					targetSubtitleCount++
				}
			}
			outputFile := mkvMergeOutputPath(targetFile)

			if err := mkv.RemoveTempOutputIfExists(outputFile); err != nil {
				return err
			}

			mergeArgs := mkv.BuildMergeFFmpegArgs(targetFile, subtitleFile, targetSubtitleCount, languageTag, subtitleTitle)
			mergeArgs = append(mergeArgs, outputFile)

			mergeOutput, err := mkv.RunFFmpeg(mergeArgs...)
			if err != nil {
				return fmt.Errorf("failed to merge subtitle: %w: %s", err, bytes.TrimSpace(mergeOutput))
			}

			if _, err := cmd.OutOrStdout().Write([]byte(fmt.Sprintf("Exported subtitle stream to %s\n", outputFile))); err != nil {
				return err
			}

			return os.Rename(outputFile, targetFile)
		},
	}

	cmd.Flags().StringVar(&targetFile, "target", "", "Target mkv file")
	_ = cmd.MarkFlagRequired("target")
	cmd.Flags().StringVar(&languageTag, "language", "", "Subtitle language tag (lowercase, 3 letters)")
	cmd.Flags().StringVar(&subtitleTitle, "title", "", "Subtitle title")
	return cmd
}
