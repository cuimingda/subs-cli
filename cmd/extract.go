package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuimingda/subs-cli/internal/mkv"

	"github.com/spf13/cobra"
)

func NewExtractCmd() *cobra.Command {
	var streamID string
	var outputDir string

	cmd := &cobra.Command{
		Use:   "extract <mkv-file>",
		Short: "Extract all subtitle streams from an mkv file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			fileName := args[0]

			if filepath.Ext(fileName) != ".mkv" && filepath.Ext(fileName) != ".MKV" {
				return fmt.Errorf("file must be an mkv file: %s", fileName)
			}

			if _, err := os.Stat(fileName); err != nil {
				return err
			}

			if err := mkv.RequireFFmpegInstalled(); err != nil {
				return err
			}

			mkvOutputDir := filepath.Dir(fileName)
			if outputDir != "" {
				info, err := os.Stat(outputDir)
				if err != nil {
					return err
				}
				if !info.IsDir() {
					return fmt.Errorf("output is not a directory: %s", outputDir)
				}
				mkvOutputDir = outputDir
			}

			streams, err := getMKVStreams(fileName)
			if err != nil {
				return err
			}

			subtitleStreams := make([]mkvStreamInfo, 0, len(streams))
			for _, stream := range streams {
				if stream.Type == "Subtitle" {
					subtitleStreams = append(subtitleStreams, stream)
				}
			}

			selectedStreams, err := selectSubtitleStreams(subtitleStreams, streams, streamID)
			if err != nil {
				return err
			}

			outDir := mkvSubtitleOutputDir(fileName, mkvOutputDir)
			if _, err := os.Stat(outDir); err == nil {
				return fmt.Errorf("subtitle output directory already exists: %s", outDir)
			}
			if err := os.MkdirAll(outDir, 0o755); err != nil {
				return err
			}

			if _, err := fmt.Fprintf(cobraCmd.OutOrStdout(), "Found %d subtitle stream(s) to extract.\n", len(selectedStreams)); err != nil {
				return err
			}

			for _, stream := range selectedStreams {
				outputPath, err := mkvSubtitleOutputPath(fileName, mkvOutputDir, stream)
				if err != nil {
					return err
				}

				if _, err := fmt.Fprintf(
					cobraCmd.OutOrStdout(),
					"Exporting stream %s (lang=%s, format=%s) -> %s\n",
					stream.ID,
					displayOrEmpty(stream.Language),
					displayOrEmpty(stream.SubtitleFormat),
					outputPath,
				); err != nil {
					return err
				}

				extractOutput, err := mkv.RunFFmpeg(
					mkv.BuildExtractFFmpegArgs(fileName, stream, outputPath)...,
				)
				if err != nil {
					return fmt.Errorf("failed to export stream %s: %w: %s", stream.ID, err, strings.TrimSpace(string(extractOutput)))
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&streamID, "id", "", "Only export one subtitle stream by stream id (for example: 4)")
	cmd.Flags().StringVar(&outputDir, "output", "", "Output directory for extracted subtitle files")

	return cmd
}
