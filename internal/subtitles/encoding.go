package subtitles

import (
	"errors"
	"os"
	"strings"

	"github.com/saintfish/chardet"
)

const UnknownEncoding = "UNKNOWN"

type FileEncoding struct {
	FileName string
	Encoding string
}

var detectTextEncoding = defaultDetectTextEncoding

func ListCurrentDirSubtitleFileEncodings() ([]FileEncoding, error) {
	files, err := ListCurrentDirSubtitleFiles()
	if err != nil {
		return nil, err
	}

	encodings := make([]FileEncoding, 0, len(files))
	for _, file := range files {
		encoding := detectFileEncoding(file)
		encodings = append(encodings, FileEncoding{
			FileName: file,
			Encoding: encoding,
		})
	}

	return encodings, nil
}

func detectFileEncoding(file string) string {
	content, err := os.ReadFile(file)
	if err != nil || len(content) == 0 {
		return UnknownEncoding
	}

	encoding, err := detectTextEncoding(content)
	if err != nil || encoding == "" {
		return UnknownEncoding
	}

	return strings.ToUpper(encoding)
}

func defaultDetectTextEncoding(content []byte) (string, error) {
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(content)
	if err != nil {
		return "", err
	}

	if result == nil || result.Charset == "" {
		return "", errors.New("encoding detection failed")
	}

	return result.Charset, nil
}
