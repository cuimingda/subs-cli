package cmd

import (
	"fmt"

	"github.com/cuimingda/subs-cli/internal/subtitles"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List subtitle files in current directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := subtitles.ListCurrentDirSubtitleFiles()
			if err != nil {
				return err
			}

			for _, file := range files {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), file); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
