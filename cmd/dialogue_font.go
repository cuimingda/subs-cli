package cmd

import (
	"fmt"
	"strings"

	"github.com/cuimingda/subs-cli/internal/subtitles"
	"github.com/spf13/cobra"
)

var dialogueCmd = &cobra.Command{
	Use:   "dialogue",
	Short: "Dialogue subtitle operations",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var dialogueFontCmd = &cobra.Command{
	Use:   "font",
	Short: "ASS font operations",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var dialogueFontListCmd = &cobra.Command{
	Use:   "list",
	Short: "List fonts used by \\fn tags in ASS files",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := subtitles.ListDialogueFontsByAssFiles()
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
	rootCmd.AddCommand(dialogueCmd)
	dialogueCmd.AddCommand(dialogueFontCmd)
	dialogueFontCmd.AddCommand(dialogueFontListCmd)
}
