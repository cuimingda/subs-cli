package cmd

import (
	"fmt"
	"strings"

	"github.com/cuimingda/subs-cli/internal/subtitles"
	"github.com/spf13/cobra"
)

var styleCmd = &cobra.Command{
	Use:   "style",
	Short: "Style subtitle operations",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var styleFontCmd = &cobra.Command{
	Use:   "font",
	Short: "ASS style font operations",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var styleFontListCmd = &cobra.Command{
	Use:   "list",
	Short: "List font names from [V4+ Styles] in ASS files",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := subtitles.ListStyleFontsByAssFiles()
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fontText := strings.Join(entry.Fonts, ",")
			if fontText == "" {
				fontText = "None"
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", entry.FileName, fontText); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(styleCmd)
	styleCmd.AddCommand(styleFontCmd)
	styleFontCmd.AddCommand(styleFontListCmd)
}
