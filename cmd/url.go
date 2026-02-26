package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

var yaheiFontURL = "https://raw.githubusercontent.com/chengda/popular-fonts/master/%E5%BE%AE%E8%BD%AF%E9%9B%85%E9%BB%91.ttf"

func NewURLCmd() *cobra.Command {
	urlCmd := &cobra.Command{
		Use:   "url",
		Short: "Font resource links",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	urlYaheiCmd := &cobra.Command{
		Use:   "yahei",
		Short: "Print Microsoft YaHei font resource URL",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			download, err := cmd.Flags().GetBool("download")
			if err != nil {
				return err
			}
			if !download {
				_, err := io.WriteString(cmd.OutOrStdout(), yaheiFontURL+"\n")
				return err
			}

			fileName, err := deriveFilenameFromURL(yaheiFontURL)
			if err != nil {
				return err
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Downloading %s from URL...\n", fileName); err != nil {
				return err
			}
			if err := downloadFile(yaheiFontURL, fileName); err != nil {
				return err
			}
			_, err = io.WriteString(cmd.OutOrStdout(), "Download complete: "+fileName+"\n")
			return err
		},
	}
	urlYaheiCmd.Flags().BoolP("download", "d", false, "Download font file to current directory")

	urlCmd.AddCommand(urlYaheiCmd)

	return urlCmd
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
