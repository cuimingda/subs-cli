package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cuimingda/subs-cli/internal/mkv"

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

			if err := mkv.RequireFFmpegInstalled(); err != nil {
				return err
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
			if err := mkv.RemoveTempOutputIfExists(outputFile); err != nil {
				return err
			}

			ffmpegArgs := mkv.BuildRemoveFFmpegArgs(targetFile, targetStream)
			ffmpegArgs = append(ffmpegArgs, outputFile)
			mergeOutput, err := mkv.RunFFmpeg(ffmpegArgs...)
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
