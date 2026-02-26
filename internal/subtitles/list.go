package subtitles

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrNoSubtitleFiles = errors.New("no .srt or .ass subtitle files found in current directory")

func ListCurrentDirSubtitleFiles() ([]string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".srt" || ext == ".ass" {
			files = append(files, entry.Name())
		}
	}

	if len(files) == 0 {
		return nil, ErrNoSubtitleFiles
	}

	return files, nil
}
