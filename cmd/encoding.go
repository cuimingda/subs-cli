package cmd

import "github.com/spf13/cobra"

func NewEncodingCmd() *cobra.Command {
	encodingCmd := &cobra.Command{
		Use:   "encoding",
		Short: "Subtitle encoding commands",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	encodingCmd.AddCommand(NewEncodingListCmd())
	encodingCmd.AddCommand(NewEncodingResetCmd())

	return encodingCmd
}
