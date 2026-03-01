package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuimingda/subs-cli/internal/mkv"
	"github.com/spf13/cobra"
)

func NewInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <mkv-file>",
		Short: "List stream ID, type, and title for a mkv file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName := args[0]

			if filepath.Ext(fileName) != ".mkv" && filepath.Ext(fileName) != ".MKV" {
				return fmt.Errorf("file must be an mkv file: %s", fileName)
			}

			if _, err := os.Stat(fileName); err != nil {
				return err
			}

			if err := mkv.RequireFFmpegInstalled(); err != nil {
				return err
			}

			streams, err := getMKVStreams(fileName)
			if err != nil {
				return err
			}

			if len(streams) == 0 {
				return fmt.Errorf("no streams found in %s", fileName)
			}

			for _, stream := range streams {
				title := stream.Title
				if title == "" {
					title = "(EMPTY)"
				}

				if stream.Type == "Subtitle" {
					language := stream.Language
					if language == "" {
						language = "(EMPTY)"
					}
					subtitleFormat := stream.SubtitleFormat
					if subtitleFormat == "" {
						subtitleFormat = "(EMPTY)"
					}

					if _, err := fmt.Fprintf(
						cmd.OutOrStdout(),
						"%s %s %s %s %s\n",
						stream.ID,
						stream.Type,
						language,
						subtitleFormat,
						title,
					); err != nil {
						return err
					}
					continue
				}

				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s\n", stream.ID, stream.Type, title); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
