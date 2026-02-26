package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	subtitles "github.com/cuimingda/subs-cli/internal/subtitles"
	"github.com/spf13/cobra"
)

func NewFileCmd() *cobra.Command {
	fileCmd := &cobra.Command{
		Use:   "file",
		Short: "Operations on subtitle files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	fileSearchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search for current directory videos that match subtitle episode tags",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			subtitleFiles, err := subtitles.ListCurrentDirSubtitleFiles()
			if err != nil {
				return err
			}

			for _, subtitleFile := range subtitleFiles {
				ignore := colorize("ignore", "31")
				notFound := colorize("not found", "31")

				episodeTag, ok := subtitles.ExtractEpisodeTag(subtitleFile)
				if !ok {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s\n", subtitleFile, ignore); err != nil {
						return err
					}
					continue
				}

				videoFile, err := subtitles.FindVideoFileByEpisodeTag(episodeTag)
				if err != nil {
					return err
				}

				if videoFile == "" {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s\n", subtitleFile, notFound); err != nil {
						return err
					}
					continue
				}

				subtitleBase := strings.TrimSuffix(subtitleFile, filepath.Ext(subtitleFile))
				videoBase := strings.TrimSuffix(videoFile, filepath.Ext(videoFile))
				if subtitleBase == videoBase {
					suffix := colorize("(same)", "32")
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s %s\n", subtitleFile, videoFile, suffix); err != nil {
						return err
					}
					continue
				}

				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s (found)\n", subtitleFile, videoFile); err != nil {
					return err
				}
			}

			return nil
		},
	}

	fileRenameCmd := &cobra.Command{
		Use:   "rename",
		Short: "Rename subtitle files according to matching video files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			subtitleFiles, err := subtitles.ListCurrentDirSubtitleFiles()
			if err != nil {
				return err
			}

			for _, subtitleFile := range subtitleFiles {
				ignore := colorize("ignore", "31")
				notFound := colorize("not found", "31")

				episodeTag, ok := subtitles.ExtractEpisodeTag(subtitleFile)
				if !ok {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s\n", subtitleFile, ignore); err != nil {
						return err
					}
					continue
				}

				videoFile, err := subtitles.FindVideoFileByEpisodeTag(episodeTag)
				if err != nil {
					return err
				}

				if videoFile == "" {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s\n", subtitleFile, notFound); err != nil {
						return err
					}
					continue
				}

				subtitleBase := strings.TrimSuffix(subtitleFile, filepath.Ext(subtitleFile))
				videoBase := strings.TrimSuffix(videoFile, filepath.Ext(videoFile))
				if subtitleBase == videoBase {
					suffix := colorize("(same)", "32")
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s %s\n", subtitleFile, videoFile, suffix); err != nil {
						return err
					}
					continue
				}

				newName := videoBase + filepath.Ext(subtitleFile)
				if err := os.Rename(subtitleFile, newName); err != nil {
					return err
				}

				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s => %s (renamed)\n", subtitleFile, newName); err != nil {
					return err
				}
			}

			return nil
		},
	}

	fileRmCmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove all subtitle files in current directory by moving them to system trash",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			subtitleFiles, err := subtitles.ListCurrentDirSubtitleFiles()
			if err != nil {
				return err
			}

			confirmed, err := confirmAction(cmd.InOrStdin(), cmd.ErrOrStderr(), "This will remove all subtitle files in current directory (srt/ass). Continue?")
			if err != nil {
				return err
			}
			if !confirmed {
				return nil
			}

			for _, subtitleFile := range subtitleFiles {
				if err := moveToTrash(subtitleFile); err != nil {
					return err
				}
			}

			return nil
		},
	}

	fileCmd.AddCommand(fileSearchCmd)
	fileCmd.AddCommand(fileRenameCmd)
	fileCmd.AddCommand(fileRmCmd)

	return fileCmd
}

func colorize(text, color string) string {
	return "\x1b[" + color + "m" + text + "\x1b[0m"
}

func confirmAction(in io.Reader, out io.Writer, message string) (bool, error) {
	if _, err := fmt.Fprintf(out, "%s [y/N]: ", message); err != nil {
		return false, err
	}

	inReader := bufio.NewReader(in)
	response, err := inReader.ReadString('\n')
	if err != nil {
		return false, nil
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

func moveToTrash(fileName string) error {
	trashDir, err := getTrashDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(trashDir, 0o755); err != nil {
		return err
	}

	target := resolveUniqueTrashPath(trashDir, fileName)
	return os.Rename(fileName, target)
}

func resolveUniqueTrashPath(trashDir, fileName string) string {
	baseName := filepath.Base(fileName)
	ext := filepath.Ext(baseName)
	nameOnly := strings.TrimSuffix(baseName, ext)
	target := filepath.Join(trashDir, baseName)
	if _, err := os.Stat(target); err == nil {
		counter := 1
		for {
			candidate := filepath.Join(trashDir, fmt.Sprintf("%s (%d)%s", nameOnly, counter, ext))
			if _, err := os.Stat(candidate); os.IsNotExist(err) {
				return candidate
			}
			counter++
		}
	}

	return target
}

func getTrashDir() (string, error) {
	home := os.Getenv("HOME")
	switch runtime.GOOS {
	case "windows":
		return "", fmt.Errorf("unsupported platform for move to trash")
	case "darwin":
		if home == "" {
			return "", fmt.Errorf("HOME not set")
		}
		return filepath.Join(home, ".Trash"), nil
	case "linux":
		if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
			return filepath.Join(dataHome, "Trash", "files"), nil
		}
		if home == "" {
			return "", fmt.Errorf("HOME not set")
		}
		return filepath.Join(home, ".local", "share", "Trash", "files"), nil
	default:
		if home == "" {
			return "", fmt.Errorf("HOME not set")
		}
		return filepath.Join(home, ".Trash"), nil
	}
}
