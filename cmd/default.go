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

func NewDefaultCmd() *cobra.Command {
	var streamID string

	cmd := &cobra.Command{
		Use:   "default <mkv_filename>",
		Short: "Toggle default disposition for a subtitle stream in an mkv file",
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

			outputFile := mkvMergeOutputPath(targetFile)
			if err := mkv.RemoveTempOutputIfExists(outputFile); err != nil {
				return err
			}

			ffmpegArgs := mkvDefaultToggleFFmpegArgs(targetFile, streams, targetStream)
			if len(ffmpegArgs) == 0 {
				return fmt.Errorf("failed to build ffmpeg args for stream %s", streamID)
			}
			ffmpegArgs = append(ffmpegArgs, outputFile)

			defaultOutput, err := mkv.RunFFmpeg(ffmpegArgs...)
			if err != nil {
				return fmt.Errorf("failed to set default for stream %s: %w: %s", streamID, err, bytes.TrimSpace(defaultOutput))
			}

			if _, err := fmt.Fprintf(cobraCmd.OutOrStdout(), "Toggled default for stream %s\n", targetStream.ID); err != nil {
				return err
			}

			return os.Rename(outputFile, targetFile)
		},
	}

	cmd.Flags().StringVar(&streamID, "id", "", "Target subtitle stream id (pure number)")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
