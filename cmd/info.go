package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type mkvStreamInfo struct {
	ID             string
	Type           string
	Language       string
	SubtitleFormat string
	Title          string
}

func NewInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <mkv-file>",
		Short: "List stream ID, type, and title for a mkv file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName := args[0]

			if filepath.Ext(fileName) != ".mkv" && filepath.Ext(fileName) != ".MKV" {
				return fmt.Errorf("file must be an mkv file: %s", fileName)
			}

			if _, err := os.Stat(fileName); err != nil {
				return err
			}

			if _, err := exec.LookPath("ffmpeg"); err != nil {
				return fmt.Errorf("ffmpeg is not installed or not in PATH, please install ffmpeg")
			}

			streams, err := getMKVStreams(fileName)
			if err != nil {
				return err
			}

			if len(streams) == 0 {
				return fmt.Errorf("no streams found in %s", fileName)
			}

			for _, stream := range streams {
				title := stream.Title
				if title == "" {
					title = "(EMPTY)"
				}

				if stream.Type == "Subtitle" {
					language := stream.Language
					if language == "" {
						language = "(EMPTY)"
					}
					subtitleFormat := stream.SubtitleFormat
					if subtitleFormat == "" {
						subtitleFormat = "(EMPTY)"
					}

					if _, err := fmt.Fprintf(
						cmd.OutOrStdout(),
						"%s %s %s %s %s\n",
						stream.ID,
						stream.Type,
						language,
						subtitleFormat,
						title,
					); err != nil {
						return err
					}
					continue
				}

				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s\n", stream.ID, stream.Type, title); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func getMKVStreams(fileName string) ([]mkvStreamInfo, error) {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-i", fileName)
	output, err := cmd.CombinedOutput()
	streams, parseErr := parseMKVStreams(string(output))
	if parseErr != nil {
		return nil, parseErr
	}
	if len(streams) == 0 {
		return nil, err
	}
	return streams, nil
}

func parseMKVStreams(output string) ([]mkvStreamInfo, error) {
	streamLineRE := regexp.MustCompile(`^\s*Stream #(.+?):\s*([A-Za-z]+):\s*(.+)$`)
	titleRE := regexp.MustCompile(`^\s*title\s*:\s*(.+)$`)
	var streams []mkvStreamInfo

	lines := strings.Split(output, "\n")
	var lastStream *mkvStreamInfo
	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r")
		if strings.HasPrefix(strings.TrimSpace(line), "Stream mapping:") {
			break
		}

		if match := streamLineRE.FindStringSubmatch(line); match != nil {
			rawID := strings.TrimSpace(match[1])
			streamID, language := parseStreamIDAndLanguage(rawID)
			streamType := strings.TrimSpace(match[2])
			streamDesc := strings.TrimSpace(match[3])

			stream := mkvStreamInfo{
				ID:       streamID,
				Type:     streamType,
				Language: language,
			}
			if streamType == "Subtitle" {
				stream.SubtitleFormat = parseSubtitleFormat(streamDesc)
			}

			streams = append(streams, stream)
			lastStream = &streams[len(streams)-1]
			continue
		}

		if lastStream == nil {
			continue
		}

		if titleMatch := titleRE.FindStringSubmatch(line); titleMatch != nil {
			lastStream.Title = strings.TrimSpace(titleMatch[1])
		}
	}
	if len(streams) == 0 {
		return nil, fmt.Errorf("no stream lines found in ffmpeg output")
	}

	return streams, nil
}

func parseStreamIDAndLanguage(rawID string) (streamID, language string) {
	open := strings.Index(rawID, "(")
	if open < 0 {
		return strings.TrimSpace(rawID), ""
	}

	streamID = strings.TrimSpace(rawID[:open])
	rest := rawID[open+1:]
	close := strings.Index(rest, ")")
	if close < 0 {
		return streamID, ""
	}

	return streamID, strings.TrimSpace(rest[:close])
}

func parseSubtitleFormat(description string) string {
	open := strings.Index(description, "(")
	if open >= 0 {
		rest := description[open+1:]
		close := strings.Index(rest, ")")
		if close >= 0 {
			return strings.TrimSpace(rest[:close])
		}
	}

	if comma := strings.Index(description, ","); comma >= 0 {
		description = description[:comma]
	}

	return strings.TrimSpace(description)
}
