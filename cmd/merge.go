package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

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

			if filepath.Ext(subtitleFile) != ".srt" && filepath.Ext(subtitleFile) != ".ass" &&
				filepath.Ext(subtitleFile) != ".SRT" && filepath.Ext(subtitleFile) != ".ASS" {
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

			if languageTag != "" && !validLanguageTag(languageTag) {
				return fmt.Errorf("invalid language tag: %s", languageTag)
			}

			if _, err := exec.LookPath("ffmpeg"); err != nil {
				return fmt.Errorf("ffmpeg is not installed or not in PATH, please install ffmpeg")
			}

			streams, err := getMKVStreams(targetFile)
			if err != nil {
				return err
			}

			if _, err := cmd.OutOrStdout().Write([]byte(fmt.Sprintf("Found %d existing streams in target.\n", len(streams)))); err != nil {
				return err
			}

			targetStreamCount := len(streams)
			outputFile := mkvMergeOutputPath(targetFile)

			if _, err := os.Stat(outputFile); err == nil {
				if err := os.Remove(outputFile); err != nil {
					return err
				}
			}

			var ffmpegArgs []string
			ffmpegArgs = append(ffmpegArgs,
				"-hide_banner",
				"-y",
				"-i", targetFile,
				"-i", subtitleFile,
				"-c", "copy",
				"-map", "0",
				"-map", "1",
			)
			if languageTag != "" || subtitleTitle != "" {
				targetMetadata := fmt.Sprintf("%d", targetStreamCount)
				if languageTag != "" {
					ffmpegArgs = append(ffmpegArgs, "-metadata:s:s:"+targetMetadata, "language="+languageTag)
				}
				if subtitleTitle != "" {
					ffmpegArgs = append(ffmpegArgs, "-metadata:s:s:"+targetMetadata, "title="+subtitleTitle)
				}
			}
			ffmpegArgs = append(ffmpegArgs, outputFile)

			mergeCmd := exec.Command("ffmpeg", ffmpegArgs...)
			if err := mergeCmd.Run(); err != nil {
				return fmt.Errorf("failed to merge subtitle: %w", err)
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

func mkvMergeOutputPath(targetFile string) string {
	return targetFile + ".tmp_subs"
}

func validLanguageTag(language string) bool {
	re := regexp.MustCompile(`^[a-z]{3}$`)
	return re.MatchString(language)
}
