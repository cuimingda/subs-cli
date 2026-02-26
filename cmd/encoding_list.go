package cmd

import (
	"fmt"

	"github.com/cuimingda/subs-cli/internal/subtitles"
	"github.com/spf13/cobra"
)

var encodingListCmd = &cobra.Command{
	Use:   "list",
	Short: "List subtitle files and encodings in current directory",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := subtitles.ListCurrentDirSubtitleFileEncodings()
		if err != nil {
			return err
		}

		for _, file := range files {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s - %s\n", file.FileName, file.Encoding); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	encodingCmd.AddCommand(encodingListCmd)
}
