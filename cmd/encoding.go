package cmd

import "github.com/spf13/cobra"

var encodingCmd = &cobra.Command{
	Use:   "encoding",
	Short: "Subtitle encoding commands",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(encodingCmd)
}
