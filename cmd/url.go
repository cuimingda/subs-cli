package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var fontURLs = map[string]string{
	"yahei": "https://raw.githubusercontent.com/chengda/popular-fonts/master/%E5%BE%AE%E8%BD%AF%E9%9B%85%E9%BB%91.ttf",
}

func NewFontCmd() *cobra.Command {
	fontCmd := &cobra.Command{
		Use:   "font",
		Short: "Font resource operations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	fontListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available fonts and download URLs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, name := range sortedFontNames() {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", name, fontURLs[name]); err != nil {
					return err
				}
			}
			return nil
		},
	}

	fontURLCmd := &cobra.Command{
		Use:   "url <name>",
		Short: "Print download URL for a font",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fontName := strings.ToLower(args[0])
			rawURL, ok := fontURLs[fontName]
			if !ok {
				return fmt.Errorf("unsupported font: %s", args[0])
			}
			_, err := io.WriteString(cmd.OutOrStdout(), rawURL+"\n")
			return err
		},
	}

	fontDownloadCmd := &cobra.Command{
		Use:   "download <name>",
		Short: "Download a font to the current directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fontName := strings.ToLower(args[0])
			rawURL, ok := fontURLs[fontName]
			if !ok {
				return fmt.Errorf("unsupported font: %s", args[0])
			}

			fileName, err := deriveFilenameFromURL(rawURL)
			if err != nil {
				return err
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Downloading %s from URL...\n", fileName); err != nil {
				return err
			}
			if err := downloadFile(rawURL, fileName); err != nil {
				return err
			}
			_, err = io.WriteString(cmd.OutOrStdout(), "Download complete: "+fileName+"\n")
			return err
		},
	}

	fontCmd.AddCommand(fontListCmd)
	fontCmd.AddCommand(fontURLCmd)
	fontCmd.AddCommand(fontDownloadCmd)

	return fontCmd
}

func sortedFontNames() []string {
	names := make([]string, 0, len(fontURLs))
	for name := range fontURLs {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

func deriveFilenameFromURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	filename := path.Base(parsed.Path)
	if filename == "." || filename == "/" || filename == "" {
		return "", fmt.Errorf("invalid URL path")
	}

	unescaped, err := url.PathUnescape(filename)
	if err != nil {
		return filename, nil
	}

	return unescaped, nil
}

func downloadFile(rawURL, fileName string) error {
	response, err := http.Get(rawURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("download request failed: %s", response.Status)
	}

	target, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}

	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return err
	}

	return nil
}
