package subtitles

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	defaultMaxSubtitleFileBytes  = 50 * 1024 * 1024
	maxSubtitleFileBytesEnvVar = "SUBS_MAX_SUBTITLE_FILE_BYTES"
)

var maxSubtitleFileBytes int64 = defaultMaxSubtitleFileBytes

func init() {
	rawLimit := strings.TrimSpace(os.Getenv(maxSubtitleFileBytesEnvVar))
	if rawLimit == "" {
		return
	}

	limit, err := strconv.ParseInt(rawLimit, 10, 64)
	if err != nil || limit <= 0 {
		return
	}

	maxSubtitleFileBytes = limit
}

func validateSubtitleFileSize(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("invalid input path %s: expected a file", filePath)
	}

	if info.Size() > maxSubtitleFileBytes {
		return fmt.Errorf(
			"file %s is too large (%d bytes): maximum allowed is %d bytes. Set %s to override",
			filePath,
			info.Size(),
			maxSubtitleFileBytes,
			maxSubtitleFileBytesEnvVar,
		)
	}

	return nil
}
