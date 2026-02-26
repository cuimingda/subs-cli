package subtitles

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var episodeTagPattern = regexp.MustCompile(`S[0-9]{2}E[0-9]{2}`)

func ExtractEpisodeTag(fileName string) (string, bool) {
	match := episodeTagPattern.FindString(fileName)
	if match == "" {
		return "", false
	}

	return match, true
}

func FindVideoFileByEpisodeTag(episodeTag string) (string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}

	for _, ext := range []string{".mkv", ".mp4"} {
		for _, entry := range entries {
			if !entry.Type().IsRegular() {
				continue
			}

			if strings.EqualFold(filepath.Ext(entry.Name()), ext) && strings.Contains(entry.Name(), episodeTag) {
				return entry.Name(), nil
			}
		}
	}

	return "", nil
}
