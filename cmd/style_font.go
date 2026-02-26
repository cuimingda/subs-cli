package cmd

import (
	"fmt"
	"strings"

	"github.com/cuimingda/subs-cli/internal/subtitles"
	"github.com/spf13/cobra"
)

func NewStyleCmd() *cobra.Command {
	styleCmd := &cobra.Command{
		Use:   "style",
		Short: "Style subtitle operations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	styleFontCmd := &cobra.Command{
		Use:   "font",
		Short: "ASS style font operations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	styleFontListCmd := &cobra.Command{
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

	styleFontResetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset [V4+ Styles] font names to Microsoft YaHei in ASS files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := subtitles.ResetCurrentDirAssStyleFontsToMicrosoftYaHei()
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(
				cmd.OutOrStdout(),
				"Reset %d font names in %d file(s).\n",
				result.UpdatedFonts,
				result.UpdatedFiles,
			)
			return err
		},
	}

	styleCmd.AddCommand(styleFontCmd)
	styleFontCmd.AddCommand(styleFontListCmd)
	styleFontCmd.AddCommand(styleFontResetCmd)

	return styleCmd
}
