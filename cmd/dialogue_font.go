package cmd

import (
	"fmt"
	"strings"

	"github.com/cuimingda/subs-cli/internal/subtitles"
	"github.com/spf13/cobra"
)

func NewDialogueCmd() *cobra.Command {
	dialogueCmd := &cobra.Command{
		Use:   "dialogue",
		Short: "Dialogue subtitle operations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	dialogueFontCmd := &cobra.Command{
		Use:   "font",
		Short: "ASS font operations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	dialogueFontListCmd := &cobra.Command{
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

	dialogueFontPruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove \\fn font tags from ASS files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := subtitles.PruneDialogueFontTagsFromAssFiles()
			if err != nil {
				return err
			}

			if _, err := fmt.Fprintf(
				cmd.OutOrStdout(),
				"Pruned %d font tags in %d files.\n",
				result.RemovedTags,
				result.TotalAssFiles,
			); err != nil {
				return err
			}

			return nil
		},
	}

	dialogueCmd.AddCommand(dialogueFontCmd)
	dialogueFontCmd.AddCommand(dialogueFontListCmd)
	dialogueFontCmd.AddCommand(dialogueFontPruneCmd)

	return dialogueCmd
}
