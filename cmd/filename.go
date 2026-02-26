package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	subtitles "github.com/cuimingda/subs-cli/internal/subtitles"
	"github.com/spf13/cobra"
)

func NewFilenameCmd() *cobra.Command {
	filenameCmd := &cobra.Command{
		Use:   "filename",
		Short: "Operations on subtitle filenames",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	filenameSearchCmd := &cobra.Command{
		Use:     "seach",
		Aliases: []string{"search"},
		Short:   "Search for current directory videos that match subtitle episode tags",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			subtitleFiles, err := subtitles.ListCurrentDirSubtitleFiles()
			if err != nil {
				return err
			}

			for _, subtitleFile := range subtitleFiles {
				ignore := colorize("ignore", "31")
				notFound := colorize("not found", "31")

				episodeTag, ok := subtitles.ExtractEpisodeTag(subtitleFile)
				if !ok {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s\n", subtitleFile, ignore); err != nil {
						return err
					}
					continue
				}

				videoFile, err := subtitles.FindVideoFileByEpisodeTag(episodeTag)
				if err != nil {
					return err
				}

				if videoFile == "" {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s\n", subtitleFile, notFound); err != nil {
						return err
					}
					continue
				}

				subtitleBase := strings.TrimSuffix(subtitleFile, filepath.Ext(subtitleFile))
				videoBase := strings.TrimSuffix(videoFile, filepath.Ext(videoFile))
				if subtitleBase == videoBase {
					suffix := colorize("(same)", "32")
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s %s\n", subtitleFile, videoFile, suffix); err != nil {
						return err
					}
					continue
				}

				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s (found)\n", subtitleFile, videoFile); err != nil {
					return err
				}
			}

			return nil
		},
	}

	filenameCmd.AddCommand(filenameSearchCmd)

	return filenameCmd
}

func colorize(text, color string) string {
	return "\x1b[" + color + "m" + text + "\x1b[0m"
}
