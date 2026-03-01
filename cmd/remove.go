package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
)

func NewRemoveCmd() *cobra.Command {
	var streamID string

	cmd := &cobra.Command{
		Use:   "remove <mkv_filename>",
		Short: "Remove a subtitle stream from an mkv file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			targetFile := args[0]

			if filepath.Ext(targetFile) != ".mkv" && filepath.Ext(targetFile) != ".MKV" {
				return fmt.Errorf("file must be an mkv file: %s", targetFile)
			}

			if _, err := os.Stat(targetFile); err != nil {
				return err
			}

			if _, err := strconv.Atoi(streamID); err != nil {
				return fmt.Errorf("invalid stream id: %s", streamID)
			}

			if _, err := exec.LookPath("ffmpeg"); err != nil {
				return fmt.Errorf("ffmpeg is not installed or not in PATH, please install ffmpeg")
			}

			streams, err := getMKVStreams(targetFile)
			if err != nil {
				return err
			}

			targetStream, err := findStreamForSubtitleRemoval(streams, streamID)
			if err != nil {
				return err
			}

			confirmed, err := confirmAction(
				cobraCmd.InOrStdin(),
				cobraCmd.ErrOrStderr(),
				fmt.Sprintf(
					"This will remove stream id=%s, type=%s, language=%s, format=%s, title=%s",
					targetStream.ID,
					targetStream.Type,
					displayOrEmpty(targetStream.Language),
					displayOrEmpty(targetStream.SubtitleFormat),
					displayOrEmpty(targetStream.Title),
				),
			)
			if err != nil {
				return err
			}
			if !confirmed {
				return nil
			}

			outputFile := mkvMergeOutputPath(targetFile)
			if _, err := os.Stat(outputFile); err == nil {
				if err := os.Remove(outputFile); err != nil {
					return err
				}
			}

			ffmpegArgs := []string{
				"-hide_banner",
				"-y",
				"-i",
				targetFile,
				"-map",
				"0",
				"-map",
				"-" + targetStream.ID,
				"-c",
				"copy",
				outputFile,
			}
			mergeCmd := exec.Command("ffmpeg", ffmpegArgs...)
			mergeOutput, err := mergeCmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to remove stream %s: %w: %s", streamID, err, bytes.TrimSpace(mergeOutput))
			}

			if _, err := fmt.Fprintf(cobraCmd.OutOrStdout(), "Removed stream %s\n", targetStream.ID); err != nil {
				return err
			}

			return os.Rename(outputFile, targetFile)
		},
	}

	cmd.Flags().StringVar(&streamID, "id", "", "Target subtitle stream id (pure number)")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func findStreamForSubtitleRemoval(allStreams []mkvStreamInfo, targetID string) (mkvStreamInfo, error) {
	idRE := regexp.MustCompile(`^[0-9]+$`)
	if !idRE.MatchString(targetID) {
		return mkvStreamInfo{}, fmt.Errorf("invalid stream id: %s", targetID)
	}

	for _, stream := range allStreams {
		if streamIDMatch(stream.ID, targetID) {
			if stream.Type != "Subtitle" {
				return mkvStreamInfo{}, fmt.Errorf("stream id %s is not a subtitle stream", targetID)
			}
			return stream, nil
		}
	}

	return mkvStreamInfo{}, fmt.Errorf("stream id %s not found", targetID)
}
