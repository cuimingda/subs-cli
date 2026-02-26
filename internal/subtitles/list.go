package subtitles

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrNoSubtitleFiles = errors.New("当前目录下未找到 .srt 或 .ass 文件")

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
