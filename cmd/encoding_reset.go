package cmd

import (
	"fmt"

	"github.com/cuimingda/subs-cli/internal/subtitles"
	"github.com/spf13/cobra"
)

var encodingResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset subtitle file encoding to UTF-8",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := subtitles.ResetCurrentDirSubtitleFilesToUTF8()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(cmd.OutOrStdout(), "总共 %d 个文件，修改了 %d 个文件\n", result.Total, result.Updated)
		return err
	},
}

func init() {
	encodingCmd.AddCommand(encodingResetCmd)
}
